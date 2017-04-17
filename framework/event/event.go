package event

type (
	// Event defines some event
	Event interface{}

	// Router routes events
	Router interface {
		Dispatch(event Event)
	}

	// Subscriber is notified of an event
	Subscriber interface {
		Notify(Event)
	}
)
