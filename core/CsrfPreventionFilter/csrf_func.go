package CsrfPreventionFilter

import (
	"github.com/satori/go.uuid"
	"go.aoe.com/flamingo/framework/web"
)

type (
	CsrfFunc struct {
		denyStatus string
	}
)

const maxNonces = 10
const nonces = "nonces"

// Name alias for use in template
func (c *CsrfFunc) Name() string {
	return "csrftoken"
}

func (c *CsrfFunc) Func(ctx web.Context) interface{} {
	return func() interface{} {
		nonce := generateNonce()

		if ns, ok := ctx.Session().Values[nonces]; ok {
			if list, ok := ns.([]string); ok {
				ctx.Session().Values[nonces] = appendNonceToList(list, nonce)
			} else {
				ctx.Session().Values[nonces] = []string{nonce}
			}
		} else {
			ctx.Session().Values[nonces] = []string{nonce}
		}

		return nonce
	}
}

// generateNonce generates a nonce
func generateNonce() string {
	return uuid.NewV4().String()
}

func appendNonceToList(list []string, nonce string) []string {
	if len(list) > maxNonces-1 {
		diff := len(list) - maxNonces
		list = list[diff:]
	}
	return append(list, nonce)
}
