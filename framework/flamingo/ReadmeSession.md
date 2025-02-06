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

Flamingo expects a `session.Store` dingo binding, this is currently handled via the `flamingo.session.backend` config parameter.
Flamingo comes with 3 implementations for sessions: `redis`, `file` and `memory`. 

##### memory
Stores the sessions in memory only. 
There is no persistence across application restarts. 
This option should only be used for testing or local development.

##### file
Stores the sessions in the directory specified by the `flamingo.session.file` parameter (default is `/sessions`).
The files will be persisted across application restarts. 

##### redis
Stores the session in an external instance of Redis-compatible cache ([Redis](https://redis.io/), [Valkey](https://valkey.io)). 
Use the following parameters to configure the connection to redis. 

```yaml
flamingo.session.redis.url: redis://:my-secret-password@my-redis/0    # full URL (can be used instead of host, password, database)
flamingo.session.redis.host: my-redis                                 # hostname
flamingo.session.redis.password: my-secret-password                   # password
flamingo.session.redis.database: 0                                    # database
flamingo.session.redis.idle.connections: 10                           # maximum number of idle connections
flamingo.session.redis.tls: true                                      # enable tls for connections
flamingo.session.redis.clusterMode: false                             # for redis servers running in cluster mode
flamingo.session.redis.timeout: 5s                                    # timeout for establishing the connection (as time.Duration string)
flamingo.session.redis.keyPrefix: "my-session-key-prefix:"            # optional: prefix to be used for session keys
```

##### custom
You can create your own session backend by implementing the gorilla `sessions.Store` interface. 
There is a list of existing implementation in the [gorilla/sessions repository](https://github.com/gorilla/sessions/#store-implementations).
To use them, just replace the default `flamingo.SessionModule` and bind your implementation to the `session.Store` interface via dingo.

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
