package domain

import "github.com/coreos/go-oidc"

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
)

const (
	// GUEST user
	GUEST UserType = "guest"

	// USER is an authenticated user
	USER = "user"
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
