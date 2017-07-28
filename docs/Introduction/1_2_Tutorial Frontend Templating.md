# Flamingo Frontend Development

In this document you will find anything necessary to get started in developing the flamingo frontend.

## Tutorial Step 1 "Mock Time Rendering"

In this tutorial you will first learn how to work on templated - independently from the backend implementation in flamingo.
Our example shows a "message" and the current time:
* The Message should later be passed as variable from the Controller
* The Time should be retrieved by a template function

You will learn how to pass data to the templates from mock files.


### Example “What time is it?”

Create a new page `whattime.pug` inside `akl/frontend/src/templates/pages`:

```pug
extends /layouts/default
block content
  - message = message || 'It’s time!'
  h1= __('Current time: %s - %s', get('whattime').now, message)
```

Then create a new mock file `whattime.json` for that page inside `akl/frontend/src/mock`:

```json
{
  "now": "Thu Jun 29 2017 17:34:04 GMT+0200 (CEST)"
}
```

Add a new mock-page-mapping in `akl/frontend/src/mock/_mockmap.js`:

```js
module.exports = {
  'pages/whattime.pug': {
    'whattime': 'whattime'
  }
}
```

Compile the templates:

```sh
yarn run build:templates
```

Start the dev server:

```sh
yarn run dev
```

Now navigate to the frontend devserver [http://localhost:1337/whattime.html](http://localhost:1337/whattime.html)

You should see a page which prints “Current time: Thu Jun 29 2017 17:34:04 GMT+0200 (CEST) - It’s time!”.

The mock data is currently read from the mock files directly and servered by the Node dev server.
If you reload the page, the time stays the same.
And this is actually all you need for frontend development.

### Background Informations

Just to get an idea of what’s going on above the dotted line, let’s get the mock data from the mock server.

Register the route for our whattime view in `akl/config/context.yml`:

```yaml
routes:
  # default index page
  - path: /
    controller: flamingo.render(tpl="pages/home")
    name: home
  # whattime page
  - path: /whattime
    controller: flamingo.render(tpl="pages/whattime")
    name: whattime
```

Restart the dev server with the environment variable `USE_MOCK_SERVER` set to `true` using the following command:

```sh
yarn run dev:mockserver
```

The mock data will now be requested from the mock server, which we start up inside the `akl/` directory:

```sh
export CONTEXT="dev" && go run akl.go serve
```

Now navigate to [http://localhost:3210/en/whattime](http://localhost:3210/en/whattime)

You should see a page which prints “Current time: Thu Jun 29 2017 17:34:04 GMT+0200 (CEST) - It’s time!”.
Yes, it also renders the mocked data with the default message.

Remove the route for our whattime view in `akl/config/context.yml` again to follow the Step 2 of the tutorial.

## Tutorial Step 2 "Flamingo Rendering"

Next, we want to have a controller, which overrides the mocked data with “real” data.

We need 2 Controllers:
* One that returns the complete template (rendered)
* A "Datacontroller" that registeres for the `get('whattime)` template function, so that we replace the mocked data.

### Creating the flamingo "whattime" package

Create a `whattime` module in `akl` by creating a folder `akl/src/whattime` similarly as in the Hello World tutorial.

The new index controller in the whattime module will now do that job.
Here is how we structure the whattime module:

```
├─ whattime
|   ├─ interfaces
|   |   └─ controller
|   |       ├─ data.go
|   |       └─ index.go
|   └─ module.go
```

In module.go we configure the index controller and the data controller. It should look like this:

```go
package whattime

import (
	"flamingo/akl/src/whattime/interfaces/controller"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
)

type Module struct {
	RouterRegistry *router.Registry `inject:""`
}

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Route("/whattime", "whattime.index")
	m.RouterRegistry.Handle("whattime.index", new(controller.IndexController))
	m.RouterRegistry.Handle("whattime", new(controller.DataController))
}
```

In `data.go` we implement the property now which will be accessed via `get('whattime').now`. It should look like this:

```go
package controller

import (
	"flamingo/framework/web"
	"time"
)

type (
	DataController struct {}
)

func (controller *DataController) Data(ctx web.Context) interface {} {
	return struct{Now string}{Now: time.Now().String()}
}
```

Our index controller (`index.go`) will pass a random message to our `whattime` template, it should look like this:

```go
package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"math/rand"
)

// IndexController to handle /whattime
type (
	IndexController struct {
		responder.RenderAware `inject:""`
	}
)

// Get handles our GET requests
func (controller *IndexController) Get(ctx web.Context) web.Response {
	if rand.Intn(10)%2 == 0 {
		return controller.Render(ctx, "pages/whattime", struct{ Message string }{Message: "Hurry!"})
	}
	return controller.Render(ctx, "pages/whattime", struct{Message string}{Message: "Relax!"})
}
```

Restart the server:

```sh
export CONTEXT="dev" && go run akl.go serve
```

Now navigate to [http://localhost:3210/en/whattime](http://localhost:3210/en/whattime) and reload the page a few times.

You should see a page which prints the current time and either “Hurry!” or “Relax!” randomly on each page reload.

As a final step, let’s implement a whattime Pug mixin which we use inside the whattime page template.
Create the file `whattime.pug` inside the folder `akl/frontend/src/templates/mixins/` with the following content:

```pug
mixin whattime(msg)
  - msg = msg || 'It’s time!'
  h1= __('Current time: %s - %s', get('whattime').now, msg)
```

Include it in `akl/frontend/src/templates/mixins/index.pug`:

```pug
include whattime
```

Change the content of `akl/frontend/src/templates/pages/whattime.pug` so that it looks like this:

```pug
extends /layouts/default
block content
  +whattime(message)
  hr
  +whattime(message)
```

Again, navigate to [http://localhost:1337/whattime.html](http://localhost:1337/whattime.html)

You should see a page which prints “Current time: Thu Jun 29 2017 17:34:04 GMT+0200 (CEST) - It’s time!” twice.

Compile the templates:

```sh
yarn run build:templates
```

Make sure the flamingo server still runs and navigate to [http://localhost:3210/whattime](http://localhost:3210/whattime)
Reload the page a few times.

You should see a page which prints the current time twice with minor differences
and either 2x “Hurry!” or 2x “Relax!” randomly on each page reload.

Done.
