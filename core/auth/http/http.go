package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

// basicAuthIdentifier identifies users based on HTTP Basic Authentication header
type basicAuthIdentifier struct {
	users  map[string]string
	realm  string
	broker string
}

func identifierFactory(cfg config.Map) (auth.RequestIdentifier, error) {
	i := new(basicAuthIdentifier)
	var conf struct {
		Realm  string            `json:"realm"`
		Broker string            `json:"broker"`
		Users  map[string]string `json:"users"`
	}
	if err := cfg.MapInto(&conf); err != nil {
		return nil, err
	}
	i.users = conf.Users
	i.realm = conf.Realm
	i.broker = conf.Broker
	return i, nil
}

// BasicAuthIdentity transports a user identity, currently just the username
type BasicAuthIdentity struct {
	User   string
	broker string
}

// Subject is the http basic auth user
func (i *BasicAuthIdentity) Subject() string {
	return i.User
}

// Broker identity
func (i *BasicAuthIdentity) Broker() string {
	return i.broker
}

// Identify a user and match against the configured users
func (i *basicAuthIdentifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	user, pass, ok := request.Request().BasicAuth()
	if !ok {
		return nil, errors.New("no basic auth given")
	}

	if userpass, ok := i.users[user]; ok && pass == userpass {
		return &BasicAuthIdentity{User: user, broker: i.broker}, nil
	}

	return nil, errors.New("invalid credentials")
}

// Broker identifies itself
func (i *basicAuthIdentifier) Broker() string {
	return i.broker
}

// Authenticate a request by send 401 and WWW-Authenticate
func (i *basicAuthIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	return &web.Response{
		Status: http.StatusUnauthorized,
		Header: http.Header{"WWW-Authenticate": []string{fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, i.realm)}},
	}
}
