package domain

type (
	RoleSet interface {
		Roles() []Role
	}

	Role string
)

var (
	RoleUser = Role("RoleUser")
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
