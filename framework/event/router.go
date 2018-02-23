package event

import "context"

type (
	// DefaultRouter is a default event routing implementation
	DefaultRouter struct {
		Subscriber            []Subscriber            `inject:",optional"`
		SubscriberWithContext []SubscriberWithContext `inject:",optional"`
	}
)

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultRouter) Dispatch(ctx context.Context, event Event) {
	for _, s := range d.Subscriber {
		s.Notify(event)
	}
	for _, s := range d.SubscriberWithContext {
		s.NotifyWithContext(ctx, event)
	}
}
