package context

type Context struct {
	Name    string
	BaseUrl string

	Parent *Context

	Configuration map[string]string

	Routes map[string]string

	Handler map[string]interface{}
}

func (c *Context) Flat() *Context {
	var res *Context

	*res = *c

	for path, name := range c.Parent.Flat().Routes {
		if _, ok := res.Routes[path]; !ok {
			res.Routes[path] = name
		}
	}

	for name, handler := range c.Parent.Flat().Handler {
		if _, ok := res.Handler[name]; !ok {
			res.Handler[name] = handler
		}
	}

	return res
}
