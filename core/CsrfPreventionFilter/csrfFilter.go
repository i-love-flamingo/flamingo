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
			log.Println("CSRF ERROR: session doesn't contain the nonce of the request")
			return f.Error(ctx, errors.New("session doesn't contain the csrf-nonce of the request"))
		}
		deleteNonceInSession(nonce, ctx)
		log.Println("CSRF token pass", nonce)
	}

	return chain.Next(ctx, w)
}

func getNonceList(ctx web.Context) ([]string, error) {
	if ns, ok := ctx.Session().Values[nonces]; ok {
		if list, ok := ns.([]string); ok {
			return list, nil
		} else {
			return nil, errors.New(`the session key "nonces" isn't a list'"`)
		}
	} else {
		return nil, errors.New(`session hasn't the key "nonces"`)
	}
}

func deleteNonceInSession(nonce string, ctx web.Context) error {
	list, err := getNonceList(ctx)
	if err != nil {
		return err
	}
	if !contains(list, nonce) {
		return errors.New("couldn't delete nonce of list because it doesn't exist in the list")
	}
	for i, e := range list {
		if e == nonce {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	ctx.Session().Values[nonces] = list
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
