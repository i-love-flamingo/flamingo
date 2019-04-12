package interfaces

import (
	"context"
	"errors"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"net/http"
	"net/url"
)

type (
	//Authservice - basic auth authservice impplementation
	Authservice struct {
		responder *web.Responder
		realm string
	}

	//Idendity - is a http basic auth Idendity using plantext Username and Password
	Idendity struct {
		Username string
		Password string
		role string
	}
)

var (
	_ domain.Idendity = &Idendity{}
)

// Inject for Authservice
func (o *Authservice) Inject(responder *web.Responder, config *struct {
	Users config.Slice `inject:"config:basicauth.users,optional"`
	realm string `inject:"config:basicauth.realm,optional"`
}) {
	o.responder = responder
}

func (a *Authservice) Authenticate(ctx context.Context, returnURL *url.URL) (web.Result, error) {
	redirectResponse := a.responder.URLRedirect(returnURL)
	redirectResponse.Header.Add("WWW-Authenticate", `Basic realm="`+a.realm+`", charset="UTF-8"`)
	redirectResponse.Status = http.StatusUnauthorized
	return redirectResponse,nil
}

func (a *Authservice) IsAuthenticated(ctx context.Context, r *web.Request) bool {
	_, _, ok := r.Request().BasicAuth()
	//TODO check user pw
	return ok
}

func (a *Authservice) GetIdendity(ctx context.Context, r *web.Request) (domain.Idendity, error) {
	user, pw, ok := r.Request().BasicAuth()
	if !ok {
		return nil, errors.New("No idendity")
	}

	return &Idendity{
		Password: pw,
		Username: user,
		role: "",
	}, nil
}



func (i *Idendity) User() domain.User {
	return &domain.SimpleUser{
		SubjectVal: i.Username,
		DefaultRole: "role",
		NameVal: &i.Username,
	}
}