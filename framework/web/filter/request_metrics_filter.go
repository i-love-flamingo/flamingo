package filter

import (
	"context"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/metric"

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
	// responseBytesCount counts the total number of bytes served by the application
	responseBytesCount metric.Int64Counter
	// responseCount count the number of responses served by the application
	responsesCount metric.Int64Counter
	// keyHTTPStatus defines response http status code
	keyHTTPStatus attribute.Key = "status_code"
)

func init() {
	var err error
	responseBytesCount, err = otel.Meter("flamingo.me/opentelemetry").Int64Counter("flamingo/response/bytes_count",
		metric.WithDescription("Count of responses number of bytes"), metric.WithUnit("By"))
	if err != nil {
		panic(err)
	}
	responsesCount, err = otel.Meter("flamingo.me/opentelemetry").Int64Counter("flamingo/response/count",
		metric.WithDescription("Count of number of responses"))
	if err != nil {
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

// Apply metricsFilter to request
func (r responseMetrics) Apply(ctx context.Context, rw http.ResponseWriter) error {
	var err error

	// http.StatusOK is the default case
	responseWriter := &responseWriterMetrics{rw: rw, status: http.StatusOK}

	if r.result != nil {
		err = r.result.Apply(ctx, responseWriter)
	}

	statusAttribute := keyHTTPStatus.String(strconv.Itoa(responseWriter.status/100) + "xx")
	responseBytesCount.Add(ctx, responseWriter.bytes, metric.WithAttributes(statusAttribute))
	responsesCount.Add(ctx, 1, metric.WithAttributes(statusAttribute))

	return err
}

// Filter a web request
func (r *MetricsFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	return &responseMetrics{result: chain.Next(ctx, req, w)}
}
