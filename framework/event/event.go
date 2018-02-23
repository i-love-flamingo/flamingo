package event

import (
	"context"
)

type (
	// Event defines some event
	Event interface{}

	// Router routes events
	Router interface {
		Dispatch(ctx context.Context, event Event)
	}

	// Subscriber is notified of an event
	Subscriber interface {
		Notify(Event)
	}

	// Subscriber is notified of an event
	SubscriberWithContext interface {
		NotifyWithContext(ctx context.Context, event Event)
	}
)
