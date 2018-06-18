package csrfPreventionFilter

import (
	"net/http"

	"github.com/pkg/errors"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	csrfFilter struct {
		responder.ErrorAware `inject:""`
	}
)

const (
	// Ignore is an option which can be set to ignore the csrfFilter
	Ignore router.ControllerOption = "csrf.ignore"
)

// Filter protects the system of CSRF attacks.
// It compares the nonce of the request to the nonce of the session.
// If they don't match it will return an error. A nonce could only be used once.
// That's why after filtering the request the nonce will be deleted of the session
// If the controller implements the ControllerOptionAware and the "csrf.ignore"
// option is set, this filter will be skipped.
func (f *csrfFilter) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	if ctx.Request().Method == "POST" {

		// checks if controller doesn't want to check csrf (for example the profiler)
		if options, ok := chain.Controller.(router.ControllerOptionAware); ok {
			if options.CheckOption(Ignore) {
				return chain.Next(ctx, w)
			}
		}

		// session list of csrfNonces
		list, err := getNonceList(ctx)
		if err != nil {
			return f.Error(ctx, err)
		}

		// nonce in request
		nonce, err := ctx.Form1("csrf_token")
		if err != nil {
			return f.Error(ctx, err)
		}

		// compare request nonce to session nonce
		if !contains(list, nonce) {
			return f.Error(ctx, errors.New("session doesn't contain the csrf-nonce of the request"))
		}
		deleteNonceInSession(nonce, ctx)
	} else {
		nonce, err := ctx.Query1("csrf_token")
		if err != nil {
			deleteNonceInSession(nonce, ctx)
			ctx.Request().URL.Query().Del("csrf_token")
		}
	}

	return chain.Next(ctx, w)
}

func getNonceList(ctx web.Context) ([]string, error) {
	if ns, ok := ctx.Session().Values[csrfNonces]; ok {
		if list, ok := ns.([]string); ok {
			return list, nil
		}
		return nil, errors.New(`the session key "csrfNonces" isn't a list'"`)
	}
	return nil, errors.New(`session hasn't the key "csrfNonces"`)
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
	ctx.Session().Values[csrfNonces] = list
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
