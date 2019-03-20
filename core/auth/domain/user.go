package domain

import (
	"encoding/gob"

	"flamingo.me/flamingo/v3/core/security/domain"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// UserType such as guest or user
	UserType string

	// User is a basic authenticated user
	User struct {
		Sub          string
		Name         string
		Email        string
		Salutation   string
		FirstName    string
		LastName     string
		Street       string
		ZipCode      string
		City         string
		DateOfBirth  string
		Country      string
		CustomFields map[string]string
		Type         UserType
		Groups       []string
	}

	// LoginEvent for the current session
	LoginEvent struct {
		Session *web.Session
	}

	// LogoutEvent for the current session
	LogoutEvent struct {
		Session *web.Session
	}

	// Auth information
	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)

func init() {
	gob.Register(User{})
}

// Get a custom field by the name
func (u User) Get(name string) string {
	if u.CustomFields == nil {
		return ""
	}
	return u.CustomFields[name]
}

const (
	// GUEST user
	GUEST UserType = "guest"

	// USER is an authenticated user
	USER UserType = "user"
)

// Guest is our default guest user
var Guest = &User{
	Name: "Guest",
	Type: GUEST,
}

// OAuthRoleUser represents basic role for authorized user
var OAuthRoleUser = domain.StringRole(domain.PermissionAuthorized)
