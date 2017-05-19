package template_functions

import "flamingo/framework/context"

type (
	// ConfigFunc allows to retrieve config variables
	ConfigFunc struct {
		Context *context.Context `inject:""`
	}
)

// Name alias for use in template
func (c ConfigFunc) Name() string {
	return "config"
}

// Func as implementation of url method
func (c *ConfigFunc) Func() interface{} {
	return func(what string) interface{} {
		return c.Context.Config(what)
	}
}
