package web

import (
	"context"
	"net/http"
)

type (
	// Filter is an interface which can filter requests
	Filter interface {
		Filter(ctx context.Context, req *Request, w http.ResponseWriter, fc *FilterChain) Result
	}

	// FilterChain defines the chain which contains all filters which will be worked off
	FilterChain struct {
		filters   []Filter
		final     lastFilter // special case for the final controller
		postApply []func(err error, result Result)
	}

	lastFilter func(ctx context.Context, req *Request, w http.ResponseWriter) Result
)

func (fnc lastFilter) Filter(ctx context.Context, req *Request, w http.ResponseWriter, chain *FilterChain) Result {
	return fnc(ctx, req, w)
}

// Next calls the next filter and deletes it of the chain
func (fc *FilterChain) Next(ctx context.Context, req *Request, w http.ResponseWriter) Result {
	if len(fc.filters) == 0 {
		// filter chain ended
		return fc.final(ctx, req, w)
	}

	next := fc.filters[0]
	fc.filters = fc.filters[1:]
	return next.Filter(ctx, req, w, fc)
}

// AddPostApply adds a callback to be called after the response has been applied to the responsewriter
func (fc *FilterChain) AddPostApply(callback func(err error, result Result)) {
	fc.postApply = append(fc.postApply, callback)
}
