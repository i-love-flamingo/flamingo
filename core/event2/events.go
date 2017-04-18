package event2

import "fmt"

/**

 Usage

   EventDispatcher = new(DefaultEventDispatcher)

   // Subscribe
   Listener := //your listerner - Need to implement OnEvent
   EventDispatcher.Subscribe("login.sucess",Listener)

   // Throw a New Event
   Event := DefaultEvent{ "SomethingHappend" }
   DefaultEventDispatcher.Dispatch(Event)


 */
type (
	// Event can be everything that has a EventKey
	Event interface {
		GetEventKey() string
	}

	// Dispatcher routes events
	EventDispatcher interface {
		Dispatch(event Event)
		Subscribe(eventKey string, listener Listener)
	}

	// Listener is something that gets the Event and that can Subscribe to the EventDispatcher
	Listener interface {
		OnEvent(event Event)
	}


	DefaultEvent struct {
		Key string
	}

	DefaultEventDispatcher struct {
		subscriber map[string][]Listener
	}
)




func (d *DefaultEventDispatcher) Dispatch(event Event) {
	fmt.Printf("Incoiming event2 .. trying to dispatch %s",event)
	for _, l := range d.subscriber[event.GetEventKey()] {
		l.OnEvent(event)
	}
}


func (d *DefaultEventDispatcher) Subscribe(eventKey string, listener Listener) {
	if d.subscriber == nil {
		d.subscriber = make(map[string][]Listener)
	}
	d.subscriber[eventKey] = append(d.subscriber[eventKey], listener)
}



func (e DefaultEvent) GetEventKey() string {
	return e.Key
}
