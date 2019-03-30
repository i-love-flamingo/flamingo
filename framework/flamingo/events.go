package flamingo

import (
	"context"

	"flamingo.me/dingo"
)

type (
	// Event defines some event
	Event interface{}

	// EventRouter routes events
	EventRouter interface {
		Dispatch(ctx context.Context, event Event)
	}

	// eventSubscriber is notified of an event, and gets the current ctx passed
	eventSubscriber interface {
		Notify(ctx context.Context, event Event)
	}

	// Flamingo default lifecycle events
	StartupEvent        struct{} // dispatched when the application starts
	ServerStartEvent    struct{} // dispatched when a server is started (not for CLI commands)
	ServerShutdownEvent struct{} // dispatched when a server is stopped (not for CLI commands)
	ShutdownEvent       struct{} // dispatched when the application shuts down

	eventSubscriberProvider func() []eventSubscriber

	// DefaultEventRouter is a default event routing implementation
	DefaultEventRouter struct {
		provider eventSubscriberProvider
	}
)

// Inject eventSubscriberProvider dependency
func (d *DefaultEventRouter) Inject(provider eventSubscriberProvider) {
	d.provider = provider
}

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultEventRouter) Dispatch(ctx context.Context, event Event) {
	if d.provider == nil {
		return
	}

	for _, s := range d.provider() {
		s.Notify(ctx, event)
	}
}

// BindEventSubscriber is a helper to bind a private event Subscriber via Dingo
func BindEventSubscriber(injector *dingo.Injector) *dingo.Binding {
	return injector.BindMulti(new(eventSubscriber))
}
