package event

import "context"

type (
	subscriberProvider            func() []Subscriber
	subscriberWithContextProvider func() []SubscriberWithContext

	// DefaultRouter is a default event routing implementation
	DefaultRouter struct {
		Subscriber            subscriberProvider            `inject:",optional"`
		SubscriberWithContext subscriberWithContextProvider `inject:",optional"`
	}
)

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultRouter) Dispatch(ctx context.Context, event Event) {
	for _, s := range d.Subscriber() {
		s.Notify(event)
	}
	for _, s := range d.SubscriberWithContext() {
		s.NotifyWithContext(ctx, event)
	}
}
