package config

type (
	// ConfigTemplateFunc allows to retrieve config variables
	ConfigTemplateFunc struct {
		area *Area
	}
)

func (c *ConfigTemplateFunc) Inject(area *Area) {
	c.area = area
}

// Func as implementation of url method
func (c *ConfigTemplateFunc) Func() interface{} {
	return func(what string) interface{} {
		val, _ := c.area.Config(what)
		return val
	}
}
