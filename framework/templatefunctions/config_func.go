package templatefunctions

import "go.aoe.com/flamingo/framework/config"

type (
	// ConfigFunc allows to retrieve config variables
	ConfigFunc struct {
		Area *config.Area `inject:""`
	}
)

// Name alias for use in template
func (c ConfigFunc) Name() string {
	return "config"
}

// Func as implementation of url method
func (c *ConfigFunc) Func() interface{} {
	return func(what string) interface{} {
		val, _ := c.Area.Config(what)
		return val
	}
}
