package mock

import (
	"encoding/json"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// Identity mocks auth.Identity
	Identity struct {
		Sub    string
		broker string
	}

	// OIDCIdentity mocks oauth.OpenIDIdentity
	OIDCIdentity struct {
		Identity
		idToken     *oidc.IDToken
		tokenSource oauth2.TokenSource
		idclaims    []byte
		atclaims    []byte
	}
)

// Subject getter
func (i *Identity) Subject() string {
	return i.Sub
}

// Broker getter
func (i *Identity) Broker() string {
	return i.broker
}

// SetBroker Setter
func (i *Identity) SetBroker(broker string) {
	i.broker = broker
}

// TokenSource getter
func (i *OIDCIdentity) TokenSource() oauth2.TokenSource {
	return i.tokenSource
}

// SetTokenSource to specify a test/mock token source
func (i *OIDCIdentity) SetTokenSource(source oauth2.TokenSource) {
	i.tokenSource = source
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
