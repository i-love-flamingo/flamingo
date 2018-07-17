package templatefunctions

import "flamingo.me/flamingo/framework/config"

type (
	// ConfigFunc allows to retrieve config variables
	ConfigFunc struct {
		area *config.Area
	}
)

func (c *ConfigFunc) Inject(area *config.Area) {
	c.area = area
}

// Name alias for use in template
func (c ConfigFunc) Name() string {
	return "config"
}

// Func as implementation of url method
func (c *ConfigFunc) Func() interface{} {
	return func(what string) interface{} {
		val, _ := c.area.Config(what)
		return val
	}
}
