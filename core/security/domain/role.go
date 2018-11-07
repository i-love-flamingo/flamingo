package domain

const (
	RoleAnonymous = "RoleAnonymous"
	RoleUser      = "RoleUser"
)

type (
	Role interface {
		Id() string
		Label() string
		Role() string
	}

	RoleSet interface {
		Roles() []Role
	}

	DefaultRole string
)

var (
	_ Role = DefaultRole("")
)

func (r DefaultRole) Id() string {
	return string(r)
}

func (r DefaultRole) Label() string {
	return string(r)
}

func (r DefaultRole) Role() string {
	return string(r)
}
