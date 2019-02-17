package domain

type (
	// PermissionSet is a list of permissions
	PermissionSet interface {
		Permissions() []string
	}

	// Role is an interface for available role
	Role interface {
		Label() string
		PermissionSet
	}

	// StringRole is simple representation of Role interface
	StringRole string

	complexRole struct {
		label       string
		permissions []string
	}
)

// NewRole creates instance of Role interface
func NewRole(label string, permissions []string) Role {
	return complexRole{
		label:       label,
		permissions: permissions,
	}
}

// PermissionAuthorized is the default permission for authorized users
var PermissionAuthorized = string("PermissionAuthorized")

// Label returns the role's label
func (r StringRole) Label() string {
	return string(r)
}

// Permissions returns the list of role's permissions
func (r StringRole) Permissions() []string {
	return []string{string(r)}
}

// Label returns the role's label
func (r complexRole) Label() string {
	return r.label
}

// Permissions returns the list of role's permissions
func (r complexRole) Permissions() []string {
	return r.permissions
}
