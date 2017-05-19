# Packages

Own packages can be added to a flamingo project. 
Check the prefered [package structure](structure.md).


## Routing concept

The URL structure in flamingo consists of a baseurl + the rest of the url.
The baseurl is matched against the configured contexts and determines the context that should be loaded.

For the part of the url behind the baseurl (=path) a standard routing concept is applied, where at the end the URL is matched to a controller.

* a route is assigned to a handle. (A handle is a string that represents a unique name). A handle can be something like "cms.page.view"
* for a handle a Controller can be registered. The indirection through handles allows us to regsiter different controllers for certain handlers in different contexts.
* Basti: Allowed to register handle twice - or already validated?

Routes can be configured the following ways:
* In Register func of your package
* As part of the project configuration. This again allows us to have different routing paths configured for differen contexts.


## packagename/Register.go

A package can/should be registered in the relevant context(s).

Therefore it should provide a function "Register" in the package root. 

Here you can normaly do:
* Register own routes with own routing handles
* Register a Controller Implementation (Assign a Controller to a Handle)
* Provide WidgetControllers (Basti?)
* Register an Implementation of a named Interface (Secondary Port in Ports and Adapters)

### Example:
(Basti - why with lazy func - only to inject Router?)
```

// Register adds handles for cms page routes.
func Register(c *di.Container) {
	c.Register(func(r *router.Router) {
		// default handlers
		r.Handle("cms.page.view", new(PageController))

		// default routes
		r.Route("/page/{name}", "cms.page.view")
	}, "router.register")
}

```
