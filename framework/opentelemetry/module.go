package opentelemetry

import (
	"log"
	"net/http"
	"sync"

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
	"go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	registerOnce sync.Once
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
	version   = "1.0.0"
	schemaURL = "https://flamingo.me/schemas/1.0.0"
)

func (m *Module) Configure(injector *dingo.Injector) {
	registerOnce.Do(func() {
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
				// TODO: should we use a schemaURL? https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/schemas/overview.md#how-schemas-work
				schemaURL,
				attribute.String("service.name", m.serviceName),
			)),
			tracesdk.WithSampler(tracesdk.NeverSample()),
		)
		tp := tracesdk.NewTracerProvider(
			tracerProviderOptions...,
		)
		otel.SetTracerProvider(tp)
		// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/context/api-propagators.md#propagators-distribution
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	})

	// metrics
	// TODO: opentelemetry-go doesn't provide cpu metrics, is there a better solution? https://pkg.go.dev/go.opentelemetry.io/contrib/plugins/runtime
	// TODO: which Config can be tweaked?
	exp, err := prometheus.New(prometheus.Config{}, &basic.Controller{})
	if err != nil {
		log.Fatalf("Failed to initialize Prometheus exporter: %v", err)
	}
	global.SetMeterProvider(exp.MeterProvider())
	injector.BindMap((*domain.Handler)(nil), "/metrics").ToInstance(exp)
}

type correlationIDInjector struct {
	next http.RoundTripper
}

func (rt *correlationIDInjector) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.SpanFromContext(req.Context())
	if span.IsRecording() {
		req.Header.Add("X-Correlation-ID", span.SpanContext().TraceID().String())
	}
	return rt.next.RoundTrip(req)
}

type Instrumentation struct {
	tracer trace.Tracer
}

func NewInstrumentation(tp trace.TracerProvider) *Instrumentation {
	if tp == nil {
		tp = otel.GetTracerProvider()
	}
	tr := tp.Tracer(name, trace.WithInstrumentationVersion(version))
	octrace.DefaultTracer = opencensus.NewTracer(tr)
	return &Instrumentation{
		tracer: tr,
	}
}
