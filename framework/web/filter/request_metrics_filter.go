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
	MetricsFilter struct{}

	responseWriterMetrics struct {
		rw     http.ResponseWriter
		status int
		bytes  int64
	}

	responseMetrics struct {
		result web.Result
	}
)

var (
	// responseMeasure counts different HTTP responses
	responseMeasure = stats.Int64("flamingo/response/bytes", "Count of http responses by status code", stats.UnitBytes)

	// keyHTTPStatus defines response http status code
	keyHTTPStatus, _ = tag.NewKey("status_code")
)

func init() {
	if err := opencensus.View("flamingo/response/bytes_count", responseMeasure, view.Count(), keyHTTPStatus); err != nil {
		panic(err)
	}
	if err := opencensus.View("flamingo/response/bytes_sum", responseMeasure, view.Sum(), keyHTTPStatus); err != nil {
		panic(err)
	}
}

// Header to access the response writers Header
func (r *responseWriterMetrics) Header() http.Header {
	return r.rw.Header()
}

// Write to the response writer
func (r *responseWriterMetrics) Write(b []byte) (int, error) {
	written, err := r.rw.Write(b)
	r.bytes += int64(written)
	return written, err
}

// WriteHeader records the status
func (r *responseWriterMetrics) WriteHeader(statusCode int) {
	r.status = statusCode
	r.rw.WriteHeader(statusCode)
}

// Flush the inner response writer if supported, noop otherwise
func (r *responseWriterMetrics) Flush() {
	if f, ok := r.rw.(http.Flusher); ok {
		f.Flush()
	}
}

// Apply metricsFilter to request
func (r responseMetrics) Apply(ctx context.Context, rw http.ResponseWriter) error {
	var err error

	// http.StatusOK is the default case
	responseWriter := &responseWriterMetrics{rw: rw, status: http.StatusOK}

	if r.result != nil {
		err = r.result.Apply(ctx, responseWriter)
	}

	c, _ := tag.New(ctx, tag.Upsert(keyHTTPStatus, strconv.Itoa(responseWriter.status/100)+"xx"))
	stats.Record(c, responseMeasure.M(responseWriter.bytes))

	return err
}

// Filter a web request
func (r *MetricsFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	return &responseMetrics{result: chain.Next(ctx, req, w)}
}
