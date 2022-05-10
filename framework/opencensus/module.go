package opencensus

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sync"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/zipkin"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
	openzipkin "github.com/openzipkin/zipkin-go"
	reporterHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/runmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var (
	registerOnce = new(sync.Once)
	// KeyArea is the key to represent the current flamingo area
	KeyArea, _ = tag.NewKey("area")
)

// View helps to register opencensus views with the default "area" tag
func View(name string, m stats.Measure, aggr *view.Aggregation, tagKeys ...tag.Key) error {
	return view.Register(&view.View{
		Name:        name,
		Measure:     m,
		Aggregation: aggr,
		TagKeys:     append([]tag.Key{KeyArea}, tagKeys...),
	})
}

type correlationIDInjector struct {
	next http.RoundTripper
}

// RoundTrip a request
func (rt *correlationIDInjector) RoundTrip(req *http.Request) (*http.Response, error) {
	if span := trace.FromContext(req.Context()); span != nil {
		req.Header.Add("X-Correlation-ID", span.SpanContext().TraceID.String())
	}

	return rt.next.RoundTrip(req)
}

// Module registers the opencensus module which in turn enables jaeger & co
type Module struct {
	endpoint       string
	serviceName    string
	serviceAddr    string
	jaegerEnable   bool
	zipkinEnable   bool
	zipkinEndpoint string
}

// Inject dependencies
func (m *Module) Inject(
	cfg *struct {
		Endpoint       string `inject:"config:flamingo.opencensus.jaeger.endpoint"`
		ServiceName    string `inject:"config:flamingo.opencensus.serviceName"`
		ServiceAddr    string `inject:"config:flamingo.opencensus.serviceAddr"`
		JaegerEnable   bool   `inject:"config:flamingo.opencensus.jaeger.enable"`
		ZipkinEnable   bool   `inject:"config:flamingo.opencensus.zipkin.enable"`
		ZipkinEndpoint string `inject:"config:flamingo.opencensus.zipkin.endpoint"`
	},
) *Module {
	if cfg != nil {
		m.endpoint = cfg.Endpoint
		m.serviceName = cfg.ServiceName
		m.serviceAddr = cfg.ServiceAddr
		m.jaegerEnable = cfg.JaegerEnable
		m.zipkinEnable = cfg.ZipkinEnable
		m.zipkinEndpoint = cfg.ZipkinEndpoint
	}
	return m
}

// find first not-loopback ipv4 address
func localAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipnet.IP.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() == nil {
			continue
		}
		return ipnet.IP.To4().String()
	}

	// return a random loopback addr, to ensure this is unqiue at least
	return fmt.Sprintf("127.%d.%d.%d:3210", rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

// Configure the opencensus Module
func (m *Module) Configure(injector *dingo.Injector) {
	registerOnce.Do(func() {
		// For demoing purposes, always sample.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
		http.DefaultTransport = &correlationIDInjector{next: &ochttp.Transport{Base: http.DefaultTransport}}

		if m.jaegerEnable {
			// Register the Jaeger exporter to be able to retrieve
			// the collected spans.
			exporter, err := jaeger.NewExporter(jaeger.Options{
				CollectorEndpoint: m.endpoint,
				Process: jaeger.Process{
					ServiceName: m.serviceName,
					Tags: []jaeger.Tag{
						jaeger.StringTag("ip", localAddr()),
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			trace.RegisterExporter(exporter)
		}

		if m.zipkinEnable {
			localEndpoint, err := openzipkin.NewEndpoint(m.serviceName, localAddr())
			if err != nil {
				log.Fatal(err)
			}

			// The Zipkin reporter takes collected spans from the app and reports them to the backend
			// http://localhost:9411/api/v2/spans is the default for the Zipkin Span v2
			reporter := reporterHttp.NewReporter(m.zipkinEndpoint)
			// defer reporter.Close()

			// The OpenCensus exporter wraps the Zipkin reporter
			exporter := zipkin.NewExporter(reporter, localEndpoint)
			trace.RegisterExporter(exporter)
		}
	})

	if err := runmetrics.Enable(runmetrics.RunMetricOptions{
		EnableCPU:    true,
		EnableMemory: true,
	}); err != nil {
		log.Fatal(err)
	}

	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)
	injector.BindMap((*domain.Handler)(nil), "/metrics").ToInstance(exporter)
}

// CueConfig defines the opencensus config scheme
func (m *Module) CueConfig() string {
	return `
flamingo: opencensus: {
	jaeger: {
		enable: bool | *false
		endpoint: string | *"http://localhost:14268/api/traces"
	}
	zipkin: {
		enable: bool | *false
		endpoint: string | *"http://localhost:9411/api/v2/spans"
	}
	serviceName: string | *"flamingo"
	serviceAddr: string | *":13210"
	tracing: sampler: {
		whitelist: [...string]
		blacklist: [...string]
		allowParentTrace: bool | *true
	}
	publicEndpoint: bool | *true
}
`
}

// FlamingoLegacyConfigAlias maps legacy config to new
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"opencensus.jaeger.enable":                    "flamingo.opencensus.jaeger.enable",
		"opencensus.jaeger.endpoint":                  "flamingo.opencensus.jaeger.endpoint",
		"opencensus.zipkin.enable":                    "flamingo.opencensus.zipkin.enable",
		"opencensus.zipkin.endpoint":                  "flamingo.opencensus.zipkin.endpoint",
		"opencensus.serviceName":                      "flamingo.opencensus.serviceName",
		"opencensus.serviceAddr":                      "flamingo.opencensus.serviceAddr",
		"opencensus.tracing.sampler.whitelist":        "flamingo.opencensus.tracing.sampler.whitelist",
		"opencensus.tracing.sampler.blacklist":        "flamingo.opencensus.tracing.sampler.blacklist",
		"opencensus.tracing.sampler.allowParentTrace": "flamingo.opencensus.tracing.sampler.allowParentTrace",
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(systemendpoint.Module),
	}
}
