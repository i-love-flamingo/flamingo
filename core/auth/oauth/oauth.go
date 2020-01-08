package oauth

import (
	"encoding/gob"

	"golang.org/x/oauth2"
)

type (
	// TokenSourcer defines a TokenSource which is can be used to get an AccessToken vor OAuth2 flows
	TokenSourcer interface {
		TokenSource() oauth2.TokenSource
	}

	token struct {
		tokenSource oauth2.TokenSource
	}
)

func init() {
	gob.Register(oauth2.Token{})
}

// TokenSource getter
func (i token) TokenSource() oauth2.TokenSource {
	return i.tokenSource
}
