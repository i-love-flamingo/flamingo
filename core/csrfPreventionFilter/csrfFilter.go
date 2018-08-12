package csrfPreventionFilter

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
	"github.com/pkg/errors"
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

var _ router.Filter = (*csrfFilter)(nil)

// Filter protects the system of CSRF attacks.
// It compares the nonce of the request to the nonce of the session.
// If they don't match it will return an error. A nonce could only be used once.
// That's why after filtering the request the nonce will be deleted of the session
// If the controller implements the ControllerOptionAware and the "csrf.ignore"
// option is set, this filter will be skipped.
func (f *csrfFilter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	if r.Request().Method == "POST" {

		// checks if controller doesn't want to check csrf (for example the profiler)
		//if options, ok := chain.Controller.(router.ControllerOptionAware); ok {
		//	if options.CheckOption(Ignore) {
		//		return chain.Next(ctx, r, w)
		//	}
		//}

		// session list of csrfNonces
		list, err := getNonceList(ctx, r)
		if err != nil {
			return f.Error(ctx, err)
		}

		// nonce in request
		nonce, ok := r.Form1("csrf_token")
		if !ok {
			return f.Error(ctx, errors.New("no csrf_token parameter"))
		}

		// compare request nonce to session nonce
		if !contains(list, nonce) {
			return f.Error(ctx, errors.New("session doesn't contain the csrf-nonce of the request"))
		}
		deleteNonceInSession(nonce, ctx, r)
	} else {
		nonce, ok := r.Query1("csrf_token")
		if !ok {
			deleteNonceInSession(nonce, ctx, r)
			r.Request().URL.Query().Del("csrf_token")
		}
	}

	return chain.Next(ctx, r, w)
}

func getNonceList(ctx context.Context, r *web.Request) ([]string, error) {
	if ns, ok := r.Session().Values[csrfNonces]; ok {
		if list, ok := ns.([]string); ok {
			return list, nil
		}
		return nil, errors.New(`the session key "csrfNonces" isn't a list'"`)
	}
	return nil, errors.New(`session hasn't the key "csrfNonces"`)
}

func deleteNonceInSession(nonce string, ctx context.Context, r *web.Request) error {
	list, err := getNonceList(ctx, r)
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
	r.Session().Values[csrfNonces] = list
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
