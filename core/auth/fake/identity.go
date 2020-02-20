package fake

import (
	"encoding/gob"
)

type (
	// identity mocks auth.identity
	identity struct {
		subject string
		broker  string
	}

	// UserSessionData user session data stored upon successful authentication
	UserSessionData struct {
		Subject string
	}
)

func init() {
	gob.Register(UserSessionData{})
}

func (i *identity) Subject() string {
	return i.subject
}

func (i *identity) Broker() string {
	return i.broker
}
