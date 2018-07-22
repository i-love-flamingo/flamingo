package opencensus

import (
	"log"
	"net/http"
	"sync"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var registerOnce = new(sync.Once)

type Module struct {
	Endpoint    string `inject:"config:opencensus.jaeger.endpoint"`
	ServiceName string `inject:"config:opencensus.serviceName"`
}

func (m *Module) Configure(*dingo.Injector) {
	registerOnce.Do(func() {
		// For demoing purposes, always sample.
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		http.DefaultTransport = &ochttp.Transport{Base: http.DefaultTransport}

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
	})
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"opencensus": config.Map{
			"jaeger.endpoint": "http://localhost:14268",
			"serviceName":     "flamingo",
		},
	}
}
