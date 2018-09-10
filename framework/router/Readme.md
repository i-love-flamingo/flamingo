# Router Module
Provides basic routing functionality for flamingo

You can use this package either standalone, or together with the `prefixrouter` 

## Basic Routing concept

For the path of an url a standard routing concept is applied, where at the end the URL is matched to a controller.

* a route is assigned to a handle. (A handle is a string that represents a unique name). A handle can be something like "cms.page.view"
* for a handle a Controller can be registered. The indirection through handles allows us to register different controllers for certain handlers in different contexts.

Routes can be configured the following ways (See Basic Usage below):

* Via `router.Registry` in your modules initialisation (typical in `module.go`)
* As part of the project configuration. This again allows us to have different routing paths configured for different contexts.

## Basic Usage:

### Registering Routes in Modules
Routes are registered normaly during initialisation of your flamingo module.

In order to register new Routes you need to bind a new Router Module
```
router.Bind(injector, new(routes))
```

routes need to implement the type:
```
// Module defines a router Module, which is able to register routes
Module interface {
  Routes(registry *Registry)
}
```

Insides the `Routes` method you can then use the `Registry` to register new Routes.

For example:
```go
func (r *routes) Routes(registry *router.Registry) {
	registry.Route("/hello", "hello")
	registry.HandleGet("hello", r.helloController.Get)
}
```

### Registering Routes via Configuration
Add a `routes.yml` in your config folder like this:

```yaml
- path: /
  name: index
  controller: flamingo.render(tpl="index")

- path: /anotherPath
  controller: flamingo.render(tpl="index")

- path: /redirect
  controller: flamingo.redirect(to="index")
```

You can use the flamingo default controllers (see below)

## Routing Details


### Route

A route defines a mapping from a path to a "handler identifier".

The handler identifier is used to easily support reverse routing and rewriting mechanisms.

### Handler

A "handler identifier" can be mapped to one or more `Action`s, e.g.:
```go
registry.HandleGet("hello", r.helloController.Get)
registry.HandlePost("hello", r.helloController.Get)
```


### Data Controller

Views can request arbitrary data via the `data` template function.


### Route Format

The route format is based on the format the [Play Framework](https://www.playframework.com/documentation/2.5.x/ScalaRouting) is using.

Essentially there are 4 types of parts, of which the route is constructed

#### Static

A piece which is just static, such as `/foo/bar/asd`.

#### Parameter

A part with a named parameter, `/foo/:param/` which spans the request up to the next `/` or `.` (e.g. `.html`).

#### Regex

A (optionally named) regex parameter such as `/foo/$param<[0-9]+>` which captures everything the regex captures, where `param` in this example is the name of the parameter.

#### Wildcard

A wildcard which captures everything, such as `/foo/bar/*param`. Note that slashes are not escaped here!

#### Router Target

The target of a route is a controller name and optional attributes.

#### Parameters

Parameters are comma-separated identifiers.

If no parameters are specified and not brackets are used every route parameter is available as a parameter.

- `controller.view` Get's all available parameters
- `controller.view(param1, param2)` param1 and param2 will be set
- `controller.view(param1 ?= "foo", param2 = "bar")` param1 is optional, if not specified it is set to "foo". param2 is always set to "bar".

If specified parameters don't have a value or optional value and are not part of the path, then they are taken from GET parameters.

#### Catchall

It is possible to specify a catchall address, which gets all parameters and applies all "leftover" as GET parameters, use `*` to indicate a catchall.

Example:

`controller.view(param1, *)`

This is quite helpful for reverse-routing.


## Default Controller

Currently Flamingo registers the following controller:

- `flamingo.redirect(to, ...)` Redirects to `to`. All other parameters (but `to`) are passed on as URL parameters 
- `flamingo.redirectUrl(url)` Redirects to `url` 
- `flamingo.redirectPermanent(to, ...)` Redirects permanently to `to`. All other parameters (but `to`) are passed on as URL parameters 
- `flamingo.redirectPermanentUrl(url)` Redirects permanently to `url` 

## Configured routes

Beside registering routes in the code it is also possible to register them in your routes.yml.

The root node consists of an array of objects with:

- `controller`: must name a controller to execute
- `path`: optional path where this is accessable
- `name`: optional name where this will be available for reverse routing

Context routes always take precedence over normal routes!

### Example

```yml
- path: /
  controller: flamingo.redirect(to="cms.page.view", name="home")
  name: home
- path: /home
  controller: cms.page.view(name="home")
- path: /special
  controller: cms.page.view(name?="special")
```

This will result in the following accessable routes:

- `/`: Redirects to `/home` (because there is a route for `cms.page.view` with `name` set to `home`. Otherwise this would go to `/cms/home`)
- `/home`: Shows `cms.page.view(name="home")`
- `/special`: Shows `cms.page.view(name="special")`
- `/special?name=foo`: Shows `cms.page.view(name="foo")` (optional argument retrieved from GET)

The `/` route is now also available as a controller named `home`, which is just an alias for calling the `flamingo.redirect` controller with the parameters `to="cms.page.view"` and `name="home"`.

## Router filter

Router filters can be used as middleware in the dispatching process. The filters are executed before the controller action.
A router filter can be registered via dingo injection in `module.go`'s Configure function:

```go
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(router.Filter)).To(myFilter{})
}
```

A Filter must implement the `router.Filter` interface by providing a Filter function: `Filter(ctx web.Context, w http.ResponseWriter, fc *FilterChain) web.Response`.

The filters are handled in order of `dingo.Modules` as defined in `flamingo.App()` call.
You will have to return `chain.Next(ctx, w)` in your `Filter` function to call the next filter. If you return something else,
the chain will be aborted and the actual controller action will not be executed.
