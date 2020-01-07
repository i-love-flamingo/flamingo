package http

import (
	"context"
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

func identifierFactory(cfg config.Map) auth.Identifier {
	i := new(basicAuthIdentifier)
	config.Map(cfg["users"].(map[string]interface{})).MapInto(&i.users)
	i.realm = cfg["realm"].(string)
	i.broker = cfg["broker"].(string)
	return i
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

func (i *BasicAuthIdentity) Broker() string {
	return i.broker
}

// Identify a user and match against the configured users
func (i *basicAuthIdentifier) Identify(ctx context.Context, request *web.Request) auth.Identity {
	user, pass, ok := request.Request().BasicAuth()
	if !ok {
		return nil
	}

	if userpass, ok := i.users[user]; ok && pass == userpass {
		return &BasicAuthIdentity{User: user, broker: i.broker}
	}

	return nil
}

func (i *basicAuthIdentifier) Broker() string {
	return i.broker
}

func (i *basicAuthIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	return &web.Response{
		Status: http.StatusUnauthorized,
		Header: http.Header{"WWW-Authenticate": []string{fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, i.realm)}},
	}
}
