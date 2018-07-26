package opencensus

import (
	"log"
	"net/http"
	"sync"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
)

var (
	registerOnce = new(sync.Once)
	KeyArea, _   = tag.NewKey("area")
)

func View(m stats.Measure, aggr *view.Aggregation, tagKeys ...tag.Key) {
	view.Register(&view.View{
		Measure:     m,
		Aggregation: aggr,
		TagKeys:     append([]tag.Key{KeyArea}, tagKeys...),
	})
}

type Module struct {
	Endpoint     string `inject:"config:opencensus.jaeger.endpoint"`
	ServiceName  string `inject:"config:opencensus.serviceName"`
	ServiceAddr  string `inject:"config:opencensus.serviceAddr"`
	JaegerEnable bool   `inject:"config:opencensus.jaeger.enable"`
}

func (m *Module) Configure(injector *dingo.Injector) {
	registerOnce.Do(func() {
		// For demoing purposes, always sample.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		http.DefaultTransport = &ochttp.Transport{Base: http.DefaultTransport}

		if m.JaegerEnable {
			// Register the Jaeger exporter to be able to retrieve
			// the collected spans.
			exporter, err := jaeger.NewExporter(jaeger.Options{
				Endpoint:    m.Endpoint,
				ServiceName: m.ServiceName,
			})
			if err != nil {
				log.Fatal(err)
			}
			trace.RegisterExporter(exporter)
		}

		{
			exporter, err := prometheus.NewExporter(prometheus.Options{})
			if err != nil {
				log.Fatal(err)
			}
			view.RegisterExporter(exporter)
			s := http.NewServeMux()
			s.Handle("/metrics", exporter)
			go http.ListenAndServe(m.ServiceAddr, s)
		}
	})
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"opencensus": config.Map{
			"jaeger.endpoint": "http://localhost:14268",
			"jaeger.enable":   false,
			"serviceName":     "flamingo",
			"serviceAddr":     ":13210",
		},
	}
}
