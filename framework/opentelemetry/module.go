package opentelemetry

import (
	"log"
	"net/http"
	"sync"

	"go.opentelemetry.io/otel/baggage"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
	octrace "go.opencensus.io/trace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/bridge/opencensus"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var (
	createTracerOnce sync.Once
	createMeterOnce  sync.Once
	KeyArea, _       = baggage.NewKeyProperty("area")
)

type Module struct {
	endpoint       string
	serviceName    string
	jaegerEnable   bool
	zipkinEnable   bool
	zipkinEndpoint string
}

func (m *Module) Inject(
	cfg *struct {
		Endpoint       string `inject:"config:flamingo.opentelemetry.jaeger.endpoint"`
		ServiceName    string `inject:"config:flamingo.opentelemetry.serviceName"`
		JaegerEnable   bool   `inject:"config:flamingo.opentelemetry.jaeger.enable"`
		ZipkinEnable   bool   `inject:"config:flamingo.opentelemetry.zipkin.enable"`
		ZipkinEndpoint string `inject:"config:flamingo.opentelemetry.zipkin.endpoint"`
	},
) *Module {
	if cfg != nil {
		m.endpoint = cfg.Endpoint
		m.serviceName = cfg.ServiceName
		m.jaegerEnable = cfg.JaegerEnable
		m.zipkinEnable = cfg.ZipkinEnable
		m.zipkinEndpoint = cfg.ZipkinEndpoint
	}
	return m
}

const (
	name      = "instrumentation/flamingo"
	schemaURL = "https://flamingo.me/schemas/1.0.0"
)

func (m *Module) Configure(injector *dingo.Injector) {
	http.DefaultTransport = &correlationIDInjector{next: otelhttp.NewTransport(http.DefaultTransport)}
	// traces
	tracerProviderOptions := make([]tracesdk.TracerProviderOption, 0, 3)
	// Create the Jaeger exporter
	if m.jaegerEnable {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(m.endpoint)))
		if err != nil {
			log.Fatalf("Failed to initialze Jeager exporter: %v", err)
		}
		tracerProviderOptions = append(tracerProviderOptions, tracesdk.WithBatcher(exp))
	}
	// Create the Zipkin exporter
	if m.zipkinEnable {
		exp, err := zipkin.New(
			m.zipkinEndpoint,
		)
		if err != nil {
			log.Fatalf("Failed to initialize Zipkin exporter: %v", err)
		}
		tracerProviderOptions = append(tracerProviderOptions, tracesdk.WithBatcher(exp))
	}
	tracerProviderOptions = append(tracerProviderOptions,
		tracesdk.WithResource(resource.NewWithAttributes(
			schemaURL,
			attribute.String("service.name", m.serviceName),
		)),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.NeverSample(),
			tracesdk.WithLocalParentSampled(tracesdk.AlwaysSample()), tracesdk.WithLocalParentNotSampled(tracesdk.NeverSample()),
			tracesdk.WithRemoteParentSampled(tracesdk.AlwaysSample()), tracesdk.WithRemoteParentNotSampled(tracesdk.NeverSample()),
		)),
	)
	tp := tracesdk.NewTracerProvider(
		tracerProviderOptions...,
	)
	otel.SetTracerProvider(tp)
	// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/context/api-propagators.md#propagators-distribution
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// metrics
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resource.NewWithAttributes(
			schemaURL,
			attribute.String("service.name", m.serviceName),
		)),
	)
	exp, err := prometheus.New(config, c)
	if err != nil {
		log.Fatalf("Failed to initialize Prometheus exporter: %v", err)
	}
	global.SetMeterProvider(exp.MeterProvider())
	if err := runtimemetrics.Start(); err != nil {
		log.Fatal(err)
	}
	injector.BindMap((*domain.Handler)(nil), "/metrics").ToInstance(exp)
}

type correlationIDInjector struct {
	next http.RoundTripper
}

func (rt *correlationIDInjector) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(req.Context())
	if span.SpanContext().IsSampled() {
		req.Header.Add("X-Correlation-ID", span.SpanContext().TraceID().String())
	}
	return rt.next.RoundTrip(req)
}

type Instrumentation struct {
	Tracer trace.Tracer
	Meter  metric.Meter
}

var (
	tracer trace.Tracer
	meter  metric.Meter
)

func GetTracer() trace.Tracer {
	createTracerOnce.Do(func() {
		tp := otel.GetTracerProvider()
		tr := tp.Tracer(name, trace.WithInstrumentationVersion(SemVersion()))
		octrace.DefaultTracer = opencensus.NewTracer(tr)
		tracer = tr
	})
	return tracer
}

func GetMeter() metric.Meter {
	createMeterOnce.Do(func() {
		mp := global.MeterProvider()
		meter = mp.Meter(name, metric.WithInstrumentationVersion(SemVersion()))
	})
	return meter
}

func (m *Module) CueConfig() string {
	return `
flamingo: opentelemetry: {
	jaeger: {
		enable: bool | *false
		endpoint: string | *"http://localhost:14268/api/traces"
	}
	zipkin: {
		enable: bool | *false
		endpoint: string | *"http://localhost:9411/api/v2/spans"
	}
	serviceName: string | *"flamingo"
	tracing: sampler: {
		whitelist: [...string]
		blacklist: [...string]
		allowParentTrace: bool | *true
	}
	publicEndpoint: bool | *true
}
`
}
