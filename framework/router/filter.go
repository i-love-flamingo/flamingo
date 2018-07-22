package router

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/framework/web"
)

type (
	// Filter is an interface which can filter requests
	Filter interface {
		Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, fc *FilterChain) web.Response
	}

	// FilterChain defines the chain which contains all filters which will be worked off
	FilterChain struct {
		Filters    []Filter
		Controller Controller
	}

	lastFilter func(ctx context.Context, r *web.Request, w http.ResponseWriter) web.Response
)

func (fnc lastFilter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *FilterChain) web.Response {
	return fnc(ctx, r, w)
}

// Next calls the next filter and deletes it of the chain
func (fc *FilterChain) Next(ctx context.Context, r *web.Request, w http.ResponseWriter) web.Response {
	next := fc.Filters[0]

	//ctx, span := trace.StartSpan(ctx, "filter")
	//defer span.End()

	fc.Filters = fc.Filters[1:]
	return next.Filter(ctx, r, w, fc)
}
