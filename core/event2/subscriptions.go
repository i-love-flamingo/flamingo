package event2

type (
	//TODO / Needed? Subscriber is something that decided/knows which listener to subscribe to what (wireing events to listerners)
	AddSubscriptions interface {
		AddSubscriptions()
	}
	//Just a little helper to include as anynous struct in your own struct
	Subscriber struct {
		EventDispatcher EventDispatcher `inject:""`
	}
)
