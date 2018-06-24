package domain

import (
	"flamingo.me/flamingo/framework/web"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// UserType such as guest or user
	UserType string

	// User is a basic authenticated user
	User struct {
		Sub   string
		Name  string
		Email string
		Type  UserType
	}

	// LoginEvent
	LoginEvent struct {
		Context web.Context
	}

	// LogoutEvent
	LogoutEvent struct {
		Context web.Context
	}

	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)

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

// UserFromIDToken fills the user struct with the token information
func UserFromIDToken(idtoken *oidc.IDToken) *User {
	var claim struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	idtoken.Claims(&claim)

	return &User{
		Sub:   claim.Sub,
		Name:  claim.Name,
		Email: claim.Email,
		Type:  USER,
	}
}
