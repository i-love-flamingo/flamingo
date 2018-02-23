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

## Firing events

### On the Context

An Event is fired on the request via `context.EventRouter().Dispatch(ctx, event)`.

```go
type (
    IndexController struct {}
    
    MyEvent struct {
        Data string
    }
)

func (controller *IndexController) Get(ctx web.Context) web.Response {
    ctx.EventRouter().Dispatch(ctx, &MyEvent{Data: "Hello"})
}
```


## Subscribing to events

To listen to Events you need to create a "Subscriber". 
A Subscriber will get all Events and need to decide which Events it want to handle:

```go
type (
    type EventSubscriber struct {}
)

//Implement Subscriber Interface
func (subscriber *EventSubscriber) Notify(event event.Event) {
    switch event := event.(type) {
    case *MyEvent:
        subscriber.OnMyEvent(event)  // call event handler
    }
}
```

Currently Flamingo uses Dingo Multibindings to register Event Subscriber

```go
func (m *Module) Configure(injector *dingo.Injector) {
    injector.BindMulti((*event.Subscriber)(nil)).To(EventSubscriber{})
}
```
