package oauth

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo/v3/framework/web"
	"golang.org/x/oauth2"
)

type (
	// Identity defines a TokenSource which is can be used to get an AccessToken vor OAuth2 flows
	Identity interface {
		TokenSource() oauth2.TokenSource
		AccessTokenClaims(into interface{}) error
	}

	token struct {
		tokenSource oauth2.TokenSource
	}

	// AuthCodeOptioner returns an oauth2.AuthCodeOption for the broker
	AuthCodeOptioner interface {
		Options(ctx context.Context, broker string, request *web.Request) []oauth2.AuthCodeOption
	}

	authCodeOptionerProvider func() []AuthCodeOptioner

	oauth2AuthCodeOption struct{ authCodeOption oauth2.AuthCodeOption }
)

func init() {
	gob.Register(oauth2.Token{})
}

// TokenSource getter
func (i token) TokenSource() oauth2.TokenSource {
	return i.tokenSource
}

func (o oauth2AuthCodeOption) Options(context.Context, string, *web.Request) []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{o.authCodeOption}
}
