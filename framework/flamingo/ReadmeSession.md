# Flamingo Sessions

### General session usage

Session handling in Flamingo is done by the `web.Session` struct. You can get the current session by just
request `*web.Session` in your `Inject()` function.

The inner handling is done by a `gorilla.Session` object.

Sessions have a `Values` map of type `map[string]interface{}`, which can be used to store arbitrary data.

However, it is important to know that underlying `gob` is used, so it might be necessary to register
your custom types via `gob.Register(MyStruct{})` in your module's `Configure` method if you
want to make sure it is properly persisted.

Persistence is done automatically if you use `Values`.

#### Session Configuration

Flamingo expects a `session.Store` dingo binding, this is currently handled via the `session.backend` config parameter.

Flamingo comes with 3 persistence implementations for sessions: `redis`, `file` and `memory`. 
The redis backend uses the config param `session.backend.redis.host` to find the redis, e.g. `redis.host:6379`.

You can create your own one by implementing the `sessions.Store` interface and replace the default `flamingo.SessionModule`.

### Authentication

Flamingo's [`core/auth`](../3. Core Modules/OAuth.md) module provides basic OpenID connect authentication.

Given that the module is used in your project (that means registered) you can inject
the `application.AuthManager` in your controller, and use that to retrieve
User information from the context.

Please note: the auth package needs a proper session backend like redis, the cookie
backend does not provide enough space for jwt tokens.

```go
import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type Controller struct {
	authManager *application.AuthManager
	userService application.UserServiceInterface
	session     *web.Session
}

// Inject dependencies
func (c *Controller) Inject(
	authManager *application.AuthManager,
	userService application.UserServiceInterface,
	session *web.Session,
) *Controller {
	c.authManager = authManager
	c.userService = userService
	c.session = session

	return c
}

func (c *Controller) Get(ctx context.Context, r *web.Request) web.Result {
	token, err := c.authManager.IDToken(ctx, c.session)
	// ...
	user := c.userService.GetUser(ctx, c.session)

	/*
	 * client is of type http.Client, and provides
	 * a basic http client functionality.
	 * If the context belongs to a logged in user
	 * then all requests done via this client will have
	 * automatically the current OAuth2 Access Token assigned
	 */
	client, err := c.authManager.HTTPClient(ctx, c.session)
}

```
