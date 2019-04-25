# Flamingo Events

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