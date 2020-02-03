package domain

type (
	// Identity mocks auth.Identity
	Identity struct {
		subject string
		broker  string
	}
)

func NewIdentity(subject string, broker string) *Identity {
	return &Identity{subject: subject, broker: broker}
}

// Subject getter
func (i *Identity) Subject() string {
	return i.subject
}

// Broker getter
func (i *Identity) Broker() string {
	return i.broker
}
