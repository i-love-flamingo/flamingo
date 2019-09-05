package filter

import (
	"context"
	"net/http"
	"strconv"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// MetricsFilter collects status codes of HTTP responses
	MetricsFilter struct {
	}

	responseWriterMetrics struct {
		ctx        context.Context
		rw         http.ResponseWriter
		statusCode int
	}

	responseMetrics struct {
		result   web.Result
		callback func(rw *responseWriterMetrics)
	}
)

var (
	// hTTPResponseCount counts different HTTP responses
	hTTPResponseCount = stats.Int64("flamingo/request/http_response", "Count of http responses by status code", stats.UnitDimensionless)

	// keyHTTPStatus defines response http status code
	keyHTTPStatus, _ = tag.NewKey("status_code")
)

func init() {
	if err := opencensus.View("flamingo/request/http_response", hTTPResponseCount, view.Count(), keyHTTPStatus); err != nil {
		panic(err)
	}
}

func (r *responseWriterMetrics) Header() http.Header {
	return r.rw.Header()
}

func (r *responseWriterMetrics) Write(b []byte) (int, error) {
	return r.rw.Write(b)
}

func (r *responseWriterMetrics) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.rw.WriteHeader(statusCode)
}

// Apply metricsFilter to request
func (r responseMetrics) Apply(ctx context.Context, rw http.ResponseWriter) error {
	var err error
	var rWriter = &responseWriterMetrics{ctx: ctx, rw: rw}

	defer r.callback(rWriter)

	if r.result != nil {
		err = r.result.Apply(ctx, rWriter)
	}

	return err
}

// Filter a web request
func (r *MetricsFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	response := chain.Next(ctx, req, w)

	return &responseMetrics{
		result: response,
		callback: func(rw *responseWriterMetrics) {
			c, _ := tag.New(ctx, tag.Insert(keyHTTPStatus, strconv.Itoa(rw.statusCode/100*100)))
			stats.Record(c, hTTPResponseCount.M(1))
		},
	}
}
