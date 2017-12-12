package CsrfPreventionFilter

import (
	"log"

	"net/http"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	csrfFilter struct {
		responder.ErrorAware `inject:""`
	}
)

func (f *csrfFilter) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	if ctx.Request().Method == "POST" {

		// session list of nonces
		list, err := getNonceList(ctx)
		if err != nil {
			log.Println("CSRF ERROR: list doesn't exist")
			return f.Error(ctx, err)
		}

		// nonce in request
		nonce, err := ctx.Form1("csrf_token")
		if err != nil {
			log.Println("CSRF ERROR: no nonce in request")
			return f.Error(ctx, err)
		}

		// compare request nonce to session nonce
		if !contains(list, nonce) {
			log.Println("CSRF ERROR: not same nonce")
			return f.Error(ctx, errors.New("session doesn't contain the csrf-nonce of the request"))
		}
		log.Printf("It woooooorks :D", nonce)
	}

	return chain.Next(ctx, w)
}

func getNonceList(ctx web.Context) ([]string, error) {
	if ns, ok := ctx.Session().Values["nonces"]; ok {
		if list, ok := ns.([]string); ok {
			return list, nil
		} else {
			return nil, errors.New("CSRF token has no list of nonces")
		}
	} else {
		return nil, errors.New("CSRF token doesn't exist")
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
