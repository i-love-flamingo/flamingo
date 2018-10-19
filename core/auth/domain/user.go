package domain

import (
	"encoding/gob"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
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
		customFields map[string]string
		Type         UserType
	}

	// LoginEvent
	LoginEvent struct {
		Session *sessions.Session
	}

	// LogoutEvent
	LogoutEvent struct {
		Session *sessions.Session
	}

	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)

func init() {
	gob.Register(User{})
}

func (u User) Get(name string) string {
	if u.customFields == nil {
		return ""
	}
	return u.customFields[name]
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
