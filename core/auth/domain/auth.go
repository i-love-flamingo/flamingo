package domain

import (
	"encoding/gob"
	securityDomain "flamingo.me/flamingo/v3/core/security/domain"
)

type (
	// Idendity information
	Idendity interface {
		User() User
	}

	// User is a basic authenticated user
	User interface {
		Subject() string
		Email() (string,bool)
		Name() (string,bool)
		CustomField(string) (string,bool)
		Roles() []securityDomain.Role
	}

	//SimpleUser - implements User and can be used
	SimpleUser struct {
		SubjectVal string
		EmailVal *string
		NameVal *string
		CustomFieldVal map[string]string
		DefaultRole string
	}
)

var (
	_ User = &SimpleUser{}
)

func init() {
	gob.Register(SimpleUser{})
}

// Subject - returns the users subject
func (u *SimpleUser) Subject() string {
	return u.SubjectVal
}


// Email - returns email if found
func (u *SimpleUser) Email() (string,bool) {
	if u.EmailVal == nil {
		return "",false
	}
	return *u.EmailVal, true
}


// Name - returns Name of user
func (u *SimpleUser) Name() (string,bool) {
	if u.NameVal == nil {
		return "",false
	}
	return *u.NameVal, true
}


// CustomField - get a customfield by key
func (u *SimpleUser) CustomField(key string) (string,bool) {
	if val, ok := u.CustomFieldVal[key]; ok {
		return val, true
	}
	return "", false
}


// Roles return the roles
func (u *SimpleUser) Roles() []securityDomain.Role {
	role := securityDomain.NewRole(u.DefaultRole,[]string{securityDomain.PermissionAuthorized})
	return []securityDomain.Role{role}
}