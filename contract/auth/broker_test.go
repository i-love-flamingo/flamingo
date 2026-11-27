package auth_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

// staticIdentity is the identity of the broker "plain".
type staticIdentity struct{}

func (staticIdentity) Subject() string { return "static-user" }

func (staticIdentity) Broker() string { return "plain" }

// staticBroker always identifies a request.
type staticBroker struct{}

func (staticBroker) Broker() string { return "plain" }

func (staticBroker) Identify(context.Context, *web.Request) (auth.Identity, error) {
	return staticIdentity{}, nil
}

// brandOptioner adds the parameter "brand" to the authorize URL.
type brandOptioner struct{}

func (brandOptioner) Options(context.Context, string, *web.Request) []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("brand", "shop")}
}

func TestBroker_FactoryMapAddsType(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	var (
		mu       sync.Mutex
		received []config.Map
	)

	factory := auth.RequestIdentifierFactory(func(entry config.Map) (auth.RequestIdentifier, error) {
		mu.Lock()
		defer mu.Unlock()

		received = append(received, entry)

		return staticBroker{}, nil
	})

	module := dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.BindMap(new(auth.RequestIdentifierFactory), "static").ToInstance(factory)
	})

	app := bootAuth(t, "core.auth.web.broker:\n  - broker: plain\n    typ: static\n", []dingo.Module{module})

	assert.JSONEq(t, `{"loggedIn":true}`, loggedIn(t, app))

	mu.Lock()
	defer mu.Unlock()

	require.NotEmpty(t, received)

	for _, entry := range received {
		assert.Equal(t, "plain", entry["broker"])
		assert.Equal(t, "static", entry["typ"])
	}
}

func TestBroker_AuthCodeOptionerAddsParameters(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)

	module := dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.BindMulti(new(oauth.AuthCodeOptioner)).ToInstance(brandOptioner{})
	})

	app := bootAuth(t, server.BrokerYAML("shop"), []dingo.Module{module})

	authorize := login(t, app, server)

	assert.Equal(t, "shop", authorize.Query().Get("brand"))
}
