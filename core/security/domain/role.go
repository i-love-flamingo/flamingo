package domain

type (
	// RoleSet is a list of roles
	RoleSet interface {
		Roles() []Role
	}

	// Role is an available role
	Role string
)

// RoleUser is the default Role for users
var RoleUser = Role("RoleUser")

// todo check if still needed, then name ID
// Id returns the role identification
//func (r Role) Id() string {
//	return string(r)
//}

// Label returns the role's label
func (r Role) Label() string {
	return string(r)
}

// Permission returns the role's permission identifier
func (r Role) Permission() string {
	return string(r)
}
