package router

import (
	"net/http"

	"go.aoe.com/flamingo/framework/web"
)

type (
	// Filter is an interface which can filter requests
	Filter interface {
		Filter(ctx web.Context, w http.ResponseWriter, fc *FilterChain) web.Response
	}

	// FilterChain defines the chain which contains all filters which will be worked off
	FilterChain struct {
		Filters    []Filter
		Controller Controller
	}

	lastFilter func(ctx web.Context, w http.ResponseWriter) web.Response
)

func (fnc lastFilter) Filter(ctx web.Context, w http.ResponseWriter, chain *FilterChain) web.Response {
	return fnc(ctx, w)
}

// Next calls the next filter and deletes it of the chain
func (fc *FilterChain) Next(ctx web.Context, w http.ResponseWriter) web.Response {
	next := fc.Filters[0]
	fc.Filters = fc.Filters[1:]
	return next.Filter(ctx, w, fc)
}
