package web

import (
	"context"
	"net/http"
)

type (
	// Filter is an interface which can filter requests
	Filter interface {
		Filter(ctx context.Context, r *Request, w http.ResponseWriter, fc *FilterChain) Result
	}

	// FilterChain defines the chain which contains all filters which will be worked off
	FilterChain struct {
		Filters []Filter
	}

	lastFilter func(ctx context.Context, r *Request, w http.ResponseWriter) Result
)

func (fnc lastFilter) Filter(ctx context.Context, r *Request, w http.ResponseWriter, chain *FilterChain) Result {
	return fnc(ctx, r, w)
}

// Next calls the next filter and deletes it of the chain
func (fc *FilterChain) Next(ctx context.Context, r *Request, w http.ResponseWriter) Result {
	next := fc.Filters[0]

	fc.Filters = fc.Filters[1:]
	return next.Filter(ctx, r, w, fc)
}
