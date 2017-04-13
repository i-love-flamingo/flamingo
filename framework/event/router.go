package event

type (
	// Event defines something that dispatches itself to a subscriber
	Event interface {
		Dispatch(subscriber interface{})
	}

	// Router routes events
	Router interface {
		Dispatch(key interface{}, dispatcher Event)
		Subscribe(key, callback interface{})
	}

	// Subscriber is something that accepts a list of events
	Subscriber interface {
		Events() []interface{}
	}
)
