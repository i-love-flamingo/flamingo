# Hello World Flamingo

## Development setup

Check out flamingo into `$GOPATH/src/flamingo`.

Install `glide` (e.g. via `brew install glide`).

Run `glide i` to fetch dependencies.

Your entrypoint is `akl/akl.go`, this is where the application is started.

Please make sure to set your on-save setting to `go imports` in Gogland! (Preferences > Languages & Frameworks > Go > On Save)

## Module overview

A module in Flamingo is usually in one of five possible locations:

- **akl**: This is the place where project modules live
- **framework**: This is the Flamingo framework core
- **core**: This is the Flamingo core, possibly open-sourced one day, and contains general Flamingo modules
- **om3**: This is the place where OM3 specific modules go which are generic enough for multiple projects, but not intended for core
- **_thirdparty_**: essentially everything from somewhere else such as github :)

A module always consists of a `Module` struct, usually in a file called `module.go`.

This struct defines the basic module dependencies, such as `RouterRegistry`.

The `Module` struct implements `dingo.Module`:

```go
Module interface {
  Configure(injector *Injector)
}
```

The `Configure` method is responsible for the dependency injection and module registration.

You can read more in [Dependency Injection](dependency-injection.md).

## Our first module

We start our first `helloworld` module in `akl` by creating a folder `akl/src/helloworld`.

In there we place a file `module.go`, and enter the following content:

```go
package helloworld

import "flamingo/framework/dingo"

type Module struct {}

func (m *Module) Configure(injector *dingo.Injector){}
```

Now we register the modules usage by adding it in `akl/akl.go`

```go
//original list ...
[]dingo.Module{
  // ...
  new(auth.Module),
  new(AKL),
  new(helloworld.Module),  // hello world
},
```

Now we have our first module setup in Flamingo :)

## Controller

Let's try to get some life into it. Controlling works with two parts, _Routes_ and _Handlers_.

A _Route_ defines an URL path which is mapped to a controller key, e.g. `/helloworld` to `helloworld.view`.

A _Handler_ defines a controller which handles a request to a controller key, e.g. `helloworld.view` handled by `controllers.IndexController`.

Abstracting these allows us to rewrite URLs for different locales and easily replace controllers. 

To create our IndexController we first need the controller. A controller can implement multiple interfaces:

- `router.GETController`, called for `GET` requests:
```go
GETController interface {
  Get(web.Context) web.Response
}
```
- `router.POSTController`, called for `POST` requests:
```go
POSTController interface {
  Post(web.Context) web.Response
}
```
- `func(web.Context) web.Response`, called for any request
- `http.Handler`, called for any request

We start by creating our controller in the file `akl/src/helloworld/interfaces/controller/index.go`

```go
package controller

import (
	"flamingo/framework/web"
	"strings"
)

// IndexController to handle /helloworld
type IndexController struct{}

// Get handles our GET requests
func (controller *IndexController) Get(ctx web.Context) web.Response {
	return &web.ContentResponse{
		Body: strings.NewReader("Hello World!"),
	}
}
```

Our controller returns a `ContentResponse` with the `Body` set to `Hello World!`.

Now we need to tell Flamingo where to find the controller. We update our `Module`,
add the `RouterRegistry` as a dependency and create a _Route_ and a _Handler_:

```go
package helloworld

import (
	"flamingo/akl/src/helloworld/interfaces/controller"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Route("/helloworld", "helloworld.index")
	m.RouterRegistry.Handle("helloworld.index", new(controller.IndexController))
}
```

Now we start flamingo by running, in the `akl/` folder, `go run akl.go server` and open [http://localhost:3210/de/helloworld](http://localhost:3210/de/helloworld)

We should see our `Hello World!` response.

## Render a template

The controller is fine, but we want some fancier content. For this we need a template and tell our Controller to render this template.

Currently we use the `pug_template` module, but obviously this can be anything :)

Templating is a topic on it's own, for now we use a simple template `helloworld.pug` and place it in `akl/frontend/src/templates/pages/helloworld.pug`:

```pug
extends /layouts/default

block content
  h1 Hello #{name}
```

Run the frontend pipeline and compile everything, `yarn run build`.

Now it's time to render the template. Flamingo has a couple of Responders-helpers, such as:

- `RenderAware`
- `RedirectAware`
- `ErrorAware`
- `JSONAware`

These make the controller _aware_ of a certain response behaviour, such as "being aware of JSON responses". 

We make our controller `RenderAware` by injecting the corresponding helper into our `IndexController`:

```go
import ""flamingo/framework/web/responder""
// ...
type IndexController struct{
    *responder.RenderAware `inject:""`
}
```

The `IndexController` is now `RenderAware`, which means it got a new method `Render(context web.Context, tpl string, data interface{}) *web.ContentResponse`

The `tpl` variable is the name of the template, in our case `pages/helloworld`. `data` is optional Template data.

We modify our `IndexController` accordingly:

```go
func (controller *IndexController) Get(ctx web.Context) web.Response {
	return controller.Render(ctx, "pages/helloworld", struct{Name string}{Name: "World"})
}
```

## Path parameters

Now we want the "World" to be taken from the URL.

First, we change our route definition like this:

```go
m.RouterRegistry.Route("/hello/:world", "helloworld.index(world)")
```

Now `world` is a parameter available to our controller. If we omit the list of parameters in the brackets we get all path parameters.
If we have parameters in the list which are not part of the route Flamingo will use GET values to fill them up.

Now it's time to change our controller to get the request parameter via the request context:

```go
func (controller *IndexController) Get(ctx web.Context) web.Response {
	return controller.Render(ctx, "pages/helloworld", struct{Name string}{Name: ctx.MustParam1("world")})
}
```

Now open [http://localhost:3210/de/hello/world](http://localhost:3210/de/hello/world) and compare to [http://localhost:3210/de/hello/you](http://localhost:3210/de/hello/you)

When we open our page now, we have a fancy rendered template :)

## FAQ

### I can't update my dependencies

Delete the `vendor/` folder.

### om3? core? framework? Where to put my code?

Depends on where it is supposed to live, use your experience in programming to decide :)
