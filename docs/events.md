# Events


The framework comes with a concept to publish and subscribe to events

To trigger Events somewhere:

```
type (
 SomeType struct {
		EventRouter event.Router `inject:""`
 }
 	LoginSucessEvent struct {
 		UserId string
 	}
)

// ... - throws an event
func (s *SomeType) SomeAction()  {


	s.EventRouter.Dispatch(
		LoginSucessEvent{"U123213"},
	)

```


To listen to Events you need to create a "Subscriber". 
A Subscriber will get all Events and need to decide which Events it want to handle:

```

type (
	EventOrchestration struct {
	}
)

//Implement Subscriber Interface
func (s *EventOrchestration) Notify(ev event.Event) {
	fmt.Printf("Event disoatched to Cartservice %s", ev)

	switch ev := ev.(type) {
	case domain.LoginSucessEvent:
		.. your action
	}
}

```

The regsitration of a Subscriber happens normaly in your packages initialisation Module like this:

```

type (
	// Module registers our profiler
	Module struct {
		EventRouter    event.Router           `inject:""`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
		m.EventRouter.AddSubscriber(new(application.EventOrchestration))
}

```

