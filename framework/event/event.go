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

	// SubscriberWithContext is notified of an event, and gets the current ctx passed
	SubscriberWithContext interface {
		NotifyWithContext(ctx context.Context, event Event)
	}
)
