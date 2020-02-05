package domain

import (
	"encoding/gob"

	"flamingo.me/flamingo/v3/core/auth/mock"
)

type (
	// Identity mocks auth.Identity
	Identity struct {
		mock.Identity
		mock.OIDCIdentity
	}

	// UserSessionData user session data stored upon successful authentication
	UserSessionData struct {
		Subject string
	}
)

func init() {
	gob.Register(UserSessionData{})
}

// NewIdentity provider
func NewIdentity(subject string, broker string) *Identity {
	id := Identity{
		Identity: mock.Identity{Sub: subject},
	}
	id.SetBroker(broker)

	return &id
}
