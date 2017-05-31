package event

type (
	// DefaultRouter is a default event routing implementation
	DefaultRouter struct {
		Subscriber []Subscriber `inject:""`
	}
)

// Dispatch calls the event's Dispatch method on each subscriber
func (d *DefaultRouter) Dispatch(event Event) {
	for _, s := range d.Subscriber {
		s.Notify(event)
	}
}
