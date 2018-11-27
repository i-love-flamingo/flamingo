package domain

const (
	PermissionAnonymous = "RoleAnonymous"
	PermissionUser      = "RoleUser"
)

type (
	RoleSet interface {
		Roles() []Role
	}

	Role string
)

var (
	RoleAnonymous = Role(PermissionAnonymous)
	RoleUser      = Role(PermissionUser)
)

func (r Role) Id() string {
	return string(r)
}

func (r Role) Label() string {
	return string(r)
}

func (r Role) Permission() string {
	return string(r)
}
