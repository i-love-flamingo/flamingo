package oauth

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type (
	// Identity defines a TokenSource which is can be used to get an AccessToken vor OAuth2 flows
	Identity interface {
		auth.Identity
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

var (
	// OAuthTypeChecker checks the Identity for OAuth Identity
	OAuthTypeChecker = func(identity auth.Identity) bool {
		_, ok := identity.(Identity)

		return ok
	}
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

type oauth2CallIdentifier struct {
	identifier string
	provider   *oidc.Provider
	clientID   string
}

func oauth2Factory(cfg config.Map) (auth.RequestIdentifier, error) {
	var oidcConfig oidcConfig

	if err := cfg.MapInto(&oidcConfig); err != nil {
		return nil, err
	}

	provider, err := oidc.NewProvider(context.Background(), oidcConfig.Endpoint)
	if err != nil {
		return nil, err
	}

	return &oauth2CallIdentifier{
		identifier: oidcConfig.Broker,
		provider:   provider,
		clientID:   oidcConfig.ClientID,
	}, nil
}

func (identifier *oauth2CallIdentifier) Broker() string {
	return identifier.identifier
}

func (identifier *oauth2CallIdentifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	verifier := identifier.provider.Verifier(&oidc.Config{
		ClientID: identifier.clientID,
	})

	var err error
	var token *oidc.IDToken

	for _, rawToken := range request.Request().Header.Values("Authorization") {
		if !strings.HasPrefix(rawToken, "Bearer ") {
			continue
		}
		token, err = verifier.Verify(ctx, rawToken[7:])
		if err == nil {
			return &oauth2Identity{
				identifier: identifier.identifier,
				token:      token,
				rawToken:   rawToken[7:],
			}, nil
		}
	}

	return nil, fmt.Errorf("can not identify call, last error: %#w", err)
}

type oauth2Identity struct {
	identifier string
	token      *oidc.IDToken
	rawToken   string
}

var _ Identity = new(oauth2Identity)

func (identity *oauth2Identity) Broker() string {
	return identity.identifier
}

func (identity *oauth2Identity) Subject() string {
	return identity.token.Subject
}

func (identity *oauth2Identity) TokenSource() oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: identity.rawToken,
		TokenType:   "Bearer",
		Expiry:      identity.token.Expiry,
	})
}

func (identity *oauth2Identity) AccessTokenClaims(into interface{}) error {
	return identity.token.Claims(into)
}
