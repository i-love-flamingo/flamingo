# Flamingo Basics

## Events

Flamingo uses a builtin event router, which routes events for each request.

This means that events are request-scoped, so you can assume that fired events should not
cross request boundaries (and are bound to the same go routine).

### Event interfaces

An Event can be everything, usually a struct with a few fields.

```go
LoginSucessEvent struct {
    UserId string
}
```

Events should not have the current context in them!

### Firing events

An Event is fired using the `EventRouter` 

```go
type (
	IndexController struct {
		responder   *web.Responder
		eventRouter flamingo.EventRouter
	}

	MyEvent struct {
		Data string
	}
)

// Inject dependencies
func (controller *IndexController) Inject(
	eventRouter flamingo.EventRouter,
	responder *web.Responder,
) *IndexController {
	controller.responder = responder
	controller.eventRouter = eventRouter

	return controller
}

// Get the data
func (controller *IndexController) Get(ctx context.Context, r *web.Request) web.Result {
	controller.eventRouter.Dispatch(ctx, &MyEvent{Data: "Hello"})

	return controller.responder.TODO()
}
```

### Subscribing to events

To listen to events you need to create a "Subscriber". 
A Subscriber will get all events and need to decide which events it wants to handle:

```go
type (
	EventSubscriber struct{}
)

// Notify should get called by flamingo event logic
func (subscriber *EventSubscriber) Notify(ctx context.Context, event flamingo.Event) {
	if e, ok := event.(*MyEvent); ok {
		subscriber.OnMyEvent(e) // call event handler and do something
	}
}
```

Flamingo uses Dingo multibindings internally to register an event subscriber. In your module's `Configure`,
you can just call `flamingo.BindEventSubscriber` to register your subscriber.

```go
// Configure DI
func (m *MyModule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).To(new(EventSubscriber))
}
```

## Sessions

### General session usage

Session handling in Flamingo is done by the `web.Session` struct. You can get the current session by just
request `*web.Session` in your `Inject()` function.

The inner handling is done by a `gorilla.Session` object.

Sessions have a `Values` map of type `map[string]interface{}`, which can be used to store arbitrary data.

However, it is important to know that underlying `gob` is used, so it might be necessary to register
your custom types via `gob.Register(MyStruct{})` in your module's `Configure` method if you
want to make sure it is properly persisted.

Persistence is done automatically if you use `Values`.

#### Session Configuration

Flamingo expects a `session.Store` dingo binding, this is currently handled via the `session.backend` config parameter.

Flamingo comes with 3 persistence implementations for sessions: `redis`, `file` and `memory`. 
The redis backend uses the config param `session.backend.redis.host` to find the redis, e.g. `redis.host:6379`.

You can create your own one by implementing the `sessions.Store` interface and replace the default `flamingo.SessionModule`.

### Authentication

Flamingo's [`core/auth`](../3. Core Modules/Auth (OpenId Connect).md) module provides basic OpenID connect authentication.

Given that the module is used in your project (that means registered) you can inject
the `application.AuthManager` in your controller, and use that to retrieve
User information from the context.

Please note: the auth package needs a proper session backend like redis, the cookie
backend does not provide enough space for jwt tokens.

```go
import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type Controller struct {
	authManager *application.AuthManager
	userService application.UserServiceInterface
	session     *web.Session
}

// Inject dependencies
func (c *Controller) Inject(
	authManager *application.AuthManager,
	userService application.UserServiceInterface,
	session *web.Session,
) *Controller {
	c.authManager = authManager
	c.userService = userService
	c.session = session

	return c
}

func (c *Controller) Get(ctx context.Context, r *web.Request) web.Result {
	token, err := c.authManager.IDToken(ctx, c.session)
	// ...
	user := c.userService.GetUser(ctx, c.session)

	/*
	 * client is of type http.Client, and provides
	 * a basic http client functionality.
	 * If the context belongs to a logged in user
	 * then all requests done via this client will have
	 * automatically the current OAuth2 Access Token assigned
	 */
	client, err := c.authManager.HTTPClient(ctx, c.session)
}

```
