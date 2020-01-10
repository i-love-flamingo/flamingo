package mock

import (
	"context"
	"encoding/json"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// Identifier mocks request identification
	Identifier struct {
		broker   string
		identity auth.Identity
	}

	// Identity mocks auth.Identity
	Identity struct {
		Sub    string
		broker string
	}

	// OIDCIdentity mocks oauth.OpenIDIdentity
	OIDCIdentity struct {
		Identity
		idToken  *oidc.IDToken
		idclaims []byte
		atclaims []byte
	}
)

var (
	_ auth.RequestIdentifier = new(Identifier)
	_ auth.Identity          = new(Identity)
	_ oauth.OpenIDIdentity   = new(OIDCIdentity)
)

// Broker getter
func (m *Identifier) Broker() string {
	if m.broker == "" {
		return "mock"
	}
	return m.broker
}

// Identify an incoming request
func (m *Identifier) Identify(context.Context, *web.Request) (auth.Identity, error) {
	if m.identity != nil {
		return m.identity, nil
	}
	return &Identity{
		broker: m.broker,
		Sub:    "mocked",
	}, nil
}

// SetIdentity forces an identity to be returned
func (m *Identifier) SetIdentity(identity auth.Identity) {
	m.identity = identity
}

// SetBroker identity for the identifier
func (m *Identifier) SetBroker(broker string) {
	m.broker = broker
}

// Subject getter
func (i *Identity) Subject() string {
	return i.Sub
}

// Broker getter
func (i *Identity) Broker() string {
	return i.broker
}

// TokenSource panic TODO
func (i *OIDCIdentity) TokenSource() oauth2.TokenSource {
	panic("implement me")
}

// IDToken getter
func (i *OIDCIdentity) IDToken() *oidc.IDToken {
	return i.idToken
}

// SetIDTokenClaims marshals the given claims
func (i *OIDCIdentity) SetIDTokenClaims(claims interface{}) (err error) {
	i.idclaims, err = json.Marshal(claims)
	return
}

// IDTokenClaims unmarshals the given claims
func (i *OIDCIdentity) IDTokenClaims(into interface{}) error {
	return json.Unmarshal(i.idclaims, into)
}

// SetAccessTokenClaims marshals the given claims
func (i *OIDCIdentity) SetAccessTokenClaims(claims interface{}) (err error) {
	i.atclaims, err = json.Marshal(claims)
	return
}

// AccessTokenClaims unmarshals the given claims
func (i *OIDCIdentity) AccessTokenClaims(into interface{}) error {
	return json.Unmarshal(i.atclaims, into)
}
