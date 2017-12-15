package csrfPreventionFilter

import (
	"github.com/satori/go.uuid"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// CsrfFunc is exported as a template function
	CsrfFunc struct {
		Generator NonceGenerator `inject:""`
	}
	// NonceGenerator is an interface to generate a nonce
	NonceGenerator interface {
		GenerateNonce() string
	}

	uuidGenerator struct{}
)

const (
	maxNonces  = 10
	csrfNonces = "csrf_nonces"
)

// Name alias for use in template
func (c *CsrfFunc) Name() string {
	return "csrftoken"
}

// Func returns the CSRF nonce
func (c *CsrfFunc) Func(ctx web.Context) interface{} {
	return func() interface{} {
		nonce := c.Generator.GenerateNonce()

		if ns, ok := ctx.Session().Values[csrfNonces]; ok {
			if list, ok := ns.([]string); ok {
				ctx.Session().Values[csrfNonces] = appendNonceToList(list, nonce)
			} else {
				ctx.Session().Values[csrfNonces] = []string{nonce}
			}
		} else {
			ctx.Session().Values[csrfNonces] = []string{nonce}
		}

		return nonce
	}
}

func appendNonceToList(list []string, nonce string) []string {
	if len(list) > maxNonces-1 {
		diff := len(list) - maxNonces
		list = list[diff+1:]
	}
	return append(list, nonce)
}

// generateNonce generates a nonce
func (*uuidGenerator) generateNonce() string {
	return uuid.NewV4().String()
}
