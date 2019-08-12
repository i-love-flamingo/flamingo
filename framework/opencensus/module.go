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
	openzipkin "github.com/openzipkin/zipkin-go"
	reporterHttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
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
	Endpoint       string `inject:"config:opencensus.jaeger.endpoint"`
	ServiceName    string `inject:"config:opencensus.serviceName"`
	ServiceAddr    string `inject:"config:opencensus.serviceAddr"`
	JaegerEnable   bool   `inject:"config:opencensus.jaeger.enable"`
	ZipkinEnable   bool   `inject:"config:opencensus.zipkin.enable"`
	ZipkinEndpoint string `inject:"config:opencensus.zipkin.endpoint"`
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

		if m.JaegerEnable {
			// Register the Jaeger exporter to be able to retrieve
			// the collected spans.
			exporter, err := jaeger.NewExporter(jaeger.Options{
				CollectorEndpoint: m.Endpoint,
				Process: jaeger.Process{
					ServiceName: m.ServiceName,
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

		if m.ZipkinEnable {
			localEndpoint, err := openzipkin.NewEndpoint(m.ServiceName, localAddr())
			if err != nil {
				log.Fatal(err)
			}

			// The Zipkin reporter takes collected spans from the app and reports them to the backend
			// http://localhost:9411/api/v2/spans is the default for the Zipkin Span v2
			reporter := reporterHttp.NewReporter(m.ZipkinEndpoint)
			// defer reporter.Close()

			// The OpenCensus exporter wraps the Zipkin reporter
			exporter := zipkin.NewExporter(reporter, localEndpoint)
			trace.RegisterExporter(exporter)
		}
	})
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		log.Fatal(err)
	}
	view.RegisterExporter(exporter)
	injector.BindMap((*domain.Handler)(nil), "/metrics").ToInstance(exporter)
}

// DefaultConfig for opencensus module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"opencensus": config.Map{
			"jaeger.enable":   false,
			"jaeger.endpoint": "http://localhost:14268/api/traces",
			"zipkin.enable":   false,
			"zipkin.endpoint": "http://localhost:9411/api/v2/spans",
			"serviceName":     "flamingo",
			"serviceAddr":     ":13210",
			"tracing": config.Map{
				"sampler": config.Map{
					"whitelist":        config.Slice{},
					"blacklist":        config.Slice{},
					"allowParentTrace": true,
				},
			},
		},
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(systemendpoint.Module),
	}
}
