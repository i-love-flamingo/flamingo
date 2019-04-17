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

	// StartupEvent is dispatched when the application starts
	StartupEvent        struct{}

	// ServerStartEvent is dispatched when a server is started (not for CLI commands)
	ServerStartEvent    struct{}

	// ServerShutdownEvent is dispatched when a server is stopped (not for CLI commands)
	ServerShutdownEvent struct{}

	// ShutdownEvent is  dispatched when the application shuts down
	ShutdownEvent       struct{}

	eventSubscriberProvider func() []eventSubscriber

	// DefaultEventRouter is a default event routing implementation
	DefaultEventRouter struct {
		provider eventSubscriberProvider
		logger   Logger
	}
)

// Inject eventSubscriberProvider dependency
func (d *DefaultEventRouter) Inject(provider eventSubscriberProvider, logger Logger) {
	d.provider = provider
	d.logger = logger
}

func catched(ctx context.Context, logger Logger, s eventSubscriber, e Event) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	s.Notify(ctx, e)
}

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultEventRouter) Dispatch(ctx context.Context, event Event) {
	if d.provider == nil {
		return
	}

	for _, s := range d.provider() {
		catched(ctx, d.logger, s, event)
	}
}

// BindEventSubscriber is a helper to bind a private event Subscriber via Dingo
func BindEventSubscriber(injector *dingo.Injector) *dingo.Binding {
	return injector.BindMulti(new(eventSubscriber))
}
