package auth

type (
	// Identity donates an authentication object which at least identifies the authenticated subject
	Identity interface {
		Subject() string
		Broker() string
	}

	// HasRoles adds the ability to provide roles to an identity
	HasRoles interface {
		Roles() []string
	}
)
