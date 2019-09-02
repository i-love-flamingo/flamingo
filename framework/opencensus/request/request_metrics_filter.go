package request

import (
	"context"
	"net/http"
	"strconv"

	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

type (
	metricsFilter struct {
	}

	responseWriterMetrics struct {
		ctx        context.Context
		rw         http.ResponseWriter
		statusCode int
	}

	responseMetrics struct {
		result             web.Result
		trackResponseCount bool
		callback           func(rw *responseWriterMetrics)
	}
)

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
func (r *metricsFilter) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	response := chain.Next(ctx, req, w)

	return &responseMetrics{
		result: response,
		callback: func(rw *responseWriterMetrics) {
			go recordResponseStatus(ctx, rw.statusCode)
		},
	}
}

func recordResponseStatus(ctx context.Context, status int) {
	c, _ := tag.New(
		ctx,
		tag.Insert(KeyHTTPStatus, strconv.Itoa(status/100*100)),
	)
	stats.Record(c, HTTPResponseCount.M(1))
}
