package domain

const (
	PermissionAnonymous = "RoleAnonymous"
	PermissionUser      = "RoleUser"
)

type (
	Role interface {
		Id() string
		Label() string
		Permission() string
	}

	RoleSet interface {
		Roles() []Role
	}

	DefaultRole string
)

var (
	_ Role = DefaultRole("")

	RoleAnonymous = DefaultRole(PermissionAnonymous)
	RoleUser      = DefaultRole(PermissionUser)
)

func (r DefaultRole) Id() string {
	return string(r)
}

func (r DefaultRole) Label() string {
	return string(r)
}

func (r DefaultRole) Permission() string {
	return string(r)
}
