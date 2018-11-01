package csrfPreventionFilter

import (
	"context"

	"flamingo.me/flamingo/framework/session"
	"github.com/satori/go.uuid"
)

type (
	// CsrfFunc is exported as a template function
	CsrfFunc struct {
		Generator  NonceGenerator `inject:""`
		TokenLimit int            `inject:"config:csrfPreventionFilter.tokenLimit"`
	}
	// NonceGenerator is an interface to generate a nonce
	NonceGenerator interface {
		GenerateNonce() string
	}

	UuidGenerator struct{}
)

const (
	csrfNonces = "csrf_nonces"
)

// Func returns the CSRF nonce
func (c *CsrfFunc) Func(ctx context.Context) interface{} {
	return func() interface{} {
		nonce := c.Generator.GenerateNonce()

		s, _ := session.FromContext(ctx)

		if ns, ok := s.G().Values[csrfNonces]; ok {
			if list, ok := ns.([]string); ok {
				s.G().Values[csrfNonces] = appendNonceToList(list, nonce, c.TokenLimit)
			} else {
				s.G().Values[csrfNonces] = []string{nonce}
			}
		} else {
			s.G().Values[csrfNonces] = []string{nonce}
		}

		return nonce
	}
}

func appendNonceToList(list []string, nonce string, tokenLimit int) []string {
	if len(list) > tokenLimit-1 {
		diff := len(list) - tokenLimit
		list = list[diff+1:]
	}
	return append(list, nonce)
}

// generateNonce generates a nonce
func (*UuidGenerator) GenerateNonce() string {
	return uuid.NewV4().String()
}
