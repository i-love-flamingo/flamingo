package template_functions

import (
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/core/pug_template/pugast"
	"reflect"
)

type (
	// GetFunc allows templates to access the router's `get` method
	GetFunc struct {
		Router *router.Router `inject:""`
	}
)

// Name alias for use in template
func (g GetFunc) Name() string {
	return "get"
}

// Func as implementation of get method
func (g *GetFunc) Func(ctx web.Context) interface{} {
	return func(what string) interface{} {
		return fixtype(g.Router.Get(what, ctx))
	}
}

func fixtype(val interface{}) interface{} {
	if reflect.TypeOf(val).Kind() == reflect.Slice {
		for i, e := range val.([]interface{}) {
			val.([]interface{})[i] = fixtype(e)
		}
		val = pugast.Array(val.([]interface{}))
	} else if reflect.TypeOf(val).Kind() == reflect.Map {
		for k, v := range val.(map[string]interface{}) {
			val.(map[string]interface{})[k] = fixtype(v)
		}
	}
	return val
}
