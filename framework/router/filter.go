package router

import (
	"net/http"

	"go.aoe.com/flamingo/framework/web"
)

type (
	Filter interface {
		Filter(ctx web.Context, w http.ResponseWriter, fc *FilterChain) web.Response
	}

	FilterChain struct {
		filters []Filter
	}

	lastFilter func(ctx web.Context, w http.ResponseWriter) web.Response
)

func (fnc lastFilter) Filter(ctx web.Context, w http.ResponseWriter, chain *FilterChain) web.Response {
	return fnc(ctx, w)
}

func (fc *FilterChain) Next(ctx web.Context, w http.ResponseWriter) web.Response {
	next := fc.filters[0]
	fc.filters = fc.filters[1:]
	return next.Filter(ctx, w, fc)
}
