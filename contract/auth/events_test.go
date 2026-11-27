package auth_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/testapp/oidcfake"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type eventCounter struct {
	logins  atomic.Int32
	logouts atomic.Int32
}

// Notify implements the flamingo event subscriber.
func (c *eventCounter) Notify(_ context.Context, event flamingo.Event) {
	switch event.(type) {
	case *auth.WebLoginEvent:
		c.logins.Add(1)
	case *auth.WebLogoutEvent:
		c.logouts.Add(1)
	}
}

func TestEvents_LoginAndLogoutFireOnce(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	server := oidcfake.Start(t)
	counter := new(eventCounter)

	module := dingo.ModuleFunc(func(injector *dingo.Injector) {
		flamingo.BindEventSubscriber(injector).ToInstance(counter)
	})

	app := bootAuth(t, server.BrokerYAML("shop"), []dingo.Module{module})

	login(t, app, server)

	assert.EqualValues(t, 1, counter.logins.Load())
	assert.EqualValues(t, 0, counter.logouts.Load())

	logout(t, app, server)

	assert.EqualValues(t, 1, counter.logins.Load())
	assert.EqualValues(t, 1, counter.logouts.Load())
}
