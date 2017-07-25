# Routing

## Routing concept

The URL structure in flamingo consists of a baseurl + the rest of the url.
The baseurl is matched against the configured contexts and determines the context that should be loaded.

For the part of the url behind the baseurl (=path) a standard routing concept is applied, where at the end the URL is matched to a controller.

* a route is assigned to a handle. (A handle is a string that represents a unique name). A handle can be something like "cms.page.view"
* for a handle a Controller can be registered. The indirection through handles allows us to regsiter different controllers for certain handlers in different contexts.

Routes can be configured the following ways:
* Via `router.Registry` in `module.go`
* As part of the project configuration. This again allows us to have different routing paths configured for different contexts.

## Route

A route defines a mapping from a path to an handler identifier.

The handler identifier is used to easily support reverse routing and rewriting mechanisms.

## Handler

A handler is defined by mapping a path to a `http.Handler` or a controller.

### Controller types

- `http.Handler`: will call `ServeHTTP` for requests matching the route
- `GETController`: will call `Get(web.Context) web.Response` for `GET` Requests
- `POSTController`: will call `Post(web.Context) web.Response` for `POST` Requests
- `HEADController`: will call `Head(web.Context) web.Response` for `HEAD` Requests
- `DELETEController`: will call `Delete(web.Context) web.Response` for `DELETE` Requests
- `PUTController`: will call `Put(web.Context) web.Response` for `PUT` Requests
- `func(web.Context) web.Response`: called for any Request
- `DataController`: will call `Data(web.Context) interface{}` for Data requests
- `func(web.Context) interface{}`: for Data requests

## Data Controller

Views can request arbitrary data via the `get` template function.

Flamingo exposes these data controllers via their logical name at `/_flamingo/json?name=...`.
This is a default feature for Ajax-based cache holepunching etc.

Data Controller usually don't have a route, but can be mapped to a dedicated route, that makes Flamingo return data as JSON content type.

# Examples

As always an example illustrates the routing concept best, so here we have it:
(`module.go`)

```go
func (m *Module) Configure(injector *dingo.Injector) {
    // Register the controller
    m.RouterRegistry.Handle("search.view", new(interfaces.ViewController))
    
    // Map `/search` to ViewController with `type` set to `product`
    m.RouterRegistry.Route("/search", `search.view(type="product")`)
    
    // Map `/search/:type` to ViewController with `type` retrieved from the path
    m.RouterRegistry.Route("/search/:type", `search.view(type)`)
    
    // Map a controller action to a router (no METHOD specific handling)
    m.RouterRegistry.Handle("flamingo.redirect", (*controller.Redirect).Redirect)
}
```

# Route Format

The route format is based on the format the [Play Framework](https://www.playframework.com/documentation/2.5.x/ScalaRouting) is using.

Essentially there are 4 types of parts, of which the route is constructed

## Static

A piece which is just static, such as `/foo/bar/asd`.

## Parameter

A part with a named parameter, `/foo/:param/` which spans the request up to the next `/` or `.` (e.g. `.html`).

## Regex

A (optionally named) regex parameter such as `/foo/$param<[0-9]+>` which captures everything the regex captures, where `param` in this example is the name of the parameter.

## Wildcard

A wildcard which captures everything, such as `/foo/bar/*param`. Note that slashes are not escaped here!

## Router Target

The target of a route is a controller name and optional attributes.

## Parameters

Parameters are comma-separated identifiers.

If no parameters are specified and not brackets are used every route parameter is available as a parameter.

- `controller.view` Get's all available parameters
- `controller.view(param1, param2)` param1 and param2 will be set
- `controller.view(param1 ?= "foo", param2 = "bar")` param1 is optional, if not specified it is set to "foo". param2 is always set to "bar".

If specified parameters don't have a value or optional value and are not part of the path, then they are taken from GET parameters.

# Registering routes and controllers

Routes and controllers are registered at the `router.Registry`.

## Map a name to a controller
```go
registry.Handle("controller.name", new(controllers.ControllerName))
```

## Map a name to a controller action

This is necessary if you don't want a controller with Method-based matching, and instead register an action for all requests to this controller.
Flamingo takes care of injecting dependencies in the new-generated controller instance.
```go
registry.Handle("controller.Name", (*controllers.ControllerName).Action)
```

## Map a route to a controller

```go
registry.Route("/path/to/something", "controller.name")
```

The name will be `path.to.something` (outer slashes are stripped, then slashes will be converted to dots).

# Default Controller

Currently Flamingo registers the following controller:

- `flamingo.redirect(to, ...)` Redirects to `to`. All other parameters (but `to`) are passed on as URL parameters 
- `flamingo.redirectUrl(url)` Redirects to `url` 
- `flamingo.redirectPermanent(to, ...)` Redirects permanently to `to`. All other parameters (but `to`) are passed on as URL parameters 
- `flamingo.redirectPermanentUrl(url)` Redirects permanently to `url` 

# Context routes

Beside registering routes in the code it is also possible to register them in your context.yml.

The root node `routes` consists of an array of objects with:

- `controller`: must name a controller to execute
- `path`: optional path where this is accessable
- `name`: optional name where this will be available for reverse routing

Context routes always take precedence over normal routes!

## Example

```yml
routes:
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
