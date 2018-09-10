# Events

Flamingo uses a builtin event router, which routes events for each request.

This means that events are request-scoped, so you can assume that fired events should not
cross request boundaries (and are bound to the same go routine).

## Event interfaces

An Event can be everything, usually a struct with a few fields.

```go
LoginSucessEvent struct {
    UserId string
}
```

Events should not have the current context in them!

## Firing events

An Event is fired using the `EventRouter` 

```go
type (
    IndexController struct {
       EventRouter    event.Router                 `inject:""`
    }
    
    MyEvent struct {
        Data string
    }
)

func (controller *IndexController) Get(ctx web.Context) web.Response {
    controller.EventRouter.Dispatch(ctx, &MyEvent{Data: "Hello"})
}
```


## Subscribing to events

To listen to Events you need to create a "Subscriber". 
A Subscriber will get all Events and need to decide which Events it want to handle:

```go
type (
    type EventSubscriber struct {}
)


//Notify should get called by flamingo Eventlogic
func (e *EventSubscriber) NotifyWithContext(ctx context.Context, event event.Event) {
    switch event := event.(type) {
    case *MyEvent:
        subscriber.OnMyEvent(event)  // call event handler and do something
    }
}
```

Currently Flamingo uses Dingo Multibindings to register Event Subscriber

```go
func (m *Module) Configure(injector *dingo.Injector) {
    injector.BindMulti((*event.SubscriberWithContext)(nil)).To(application.EventSubscriber{})
}
```

There is also the interface `SubscriberWithContext`as well as an interface `Subscriber` that you can use when you don't need the current request context. 
