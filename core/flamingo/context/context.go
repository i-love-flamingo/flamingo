// Context's are used for multi-site setups
package context

import "flamingo/core/flamingo/service_container"

type (
	// Context defines a configuration context for multi-site setups
	Context struct {
		Name    string
		BaseUrl string

		Parent           *Context `json:"-"`
		Childs           []*Context
		ServiceContainer *service_container.ServiceContainer `json:"-"`

		Routes        []Route           `yaml:"routes"`
		Configuration map[string]string `yaml:"config"`
		Contexts      map[string]string `yaml:"contexts"`
	}

	Route struct {
		Path       string
		Controller string
		Args       map[string]string
	}
)

func New(name string, sc *service_container.ServiceContainer, childs ...*Context) *Context {
	ctx := &Context{
		Name:             name,
		ServiceContainer: sc,
		Childs:           childs,
	}

	for _, c := range childs {
		c.Parent = ctx
	}

	return ctx
}

func (c *Context) GetFlatContexts() map[string]*Context {
	res := make(map[string]*Context)
	flat := c.Flat()
	for baseurl, name := range c.Contexts {
		res[name] = flat[c.Name+`/`+name]
		res[name].BaseUrl = baseurl
		res[name].Childs = nil
		res[name].Contexts = nil
		res[name].Name = name
	}
	return res
}

func (c *Context) Flat() map[string]*Context {
	res := make(map[string]*Context)
	res[c.Name] = c

	for _, child := range c.Childs {
		for cn, flatchild := range child.Flat() {
			res[c.Name+`/`+cn] = flatchild.MergeFrom(*c)
		}
	}

	return res
}

func (c Context) MergeFrom(from Context) *Context {
	if c.Configuration == nil {
		c.Configuration = make(map[string]string)
	}

	for k, v := range from.Configuration {
		if _, ok := c.Configuration[k]; !ok {
			c.Configuration[k] = v
		}
	}

	knownhandler := make(map[string]bool)
	for _, route := range c.Routes {
		knownhandler[route.Controller] = true
	}

	for _, route := range from.Routes {
		if !knownhandler[route.Controller] {
			c.Routes = append(c.Routes, route)
		}
	}

	if c.ServiceContainer == nil {
		c.ServiceContainer = service_container.New()
	}

	/*
		keep := *c.ServiceContainer

		*c.ServiceContainer = *from.ServiceContainer

		for

		// remove known services of the same type
		for _, u := range from.ServiceContainer.Unnamed() {
			c.ServiceContainer.Remove(u.Value)
			c.ServiceContainer.Register(u.Value)
		}
	*/

	return &c
}
