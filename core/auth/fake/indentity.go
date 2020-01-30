package fake

type (
	// Identity mocks auth.Identity
	Identity struct {
		subject string
		broker  string
	}
)

// Subject getter
func (i *Identity) Subject() string {
	return i.subject
}

// Broker getter
func (i *Identity) Broker() string {
	return i.broker
}

// SetBroker Setter
func (i *Identity) SetBroker(broker string) {
	i.broker = broker
}
