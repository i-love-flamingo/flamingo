# Sessions

## General session usage

Session handling in Flamingo is bound to the `web.Context`.

`web.Context` has a `Session()` method returning a `gorilla.Session` object, which
the programmer can assume is properly persisted and handled.

Sessions have a `Values` map of type `map[string]interface{}`, which can be used to store arbitrary data.

However, it is important to know that underlying `gob` is used, so it might be necessary to register
your custom types via `gob.Register(MyStruct{})` in your module's `Configure` method if you
want to make sure it is properly persisted.

Persistence is done automatically if you use `Values`.

## Authentication

Flamingo's `core/auth` module provides basic OpenID connect authentication.

Given that the module is used in your project (that means registered) you can inject
the `application.AuthManager` in your controller, and use that to retrieve
User information from the context.

Please note: the auth package needs a proper session backend like redis, the cookie
backend does not provide enough space for jwt tokens.

```go
import (
    "go.aoe.com/flamingo/core/auth/application"
    "go.aoe.com/flamingo/core/auth/domain"
)

type Controller struct {
    AuthManager *application.AuthManager `inject:""`
}

func (c *Controller) Get(ctx web.Context) web.Response {
    token, err := c.AuthManager.IdToken(ctx)
    // ...
    user := domain.UserFromIDToken(token)  // user contains the User information obtained from the ID token
    
    client, err := c.AuthManager.HttpClient(ctx)
    /*
     * client is of type http.Client, and provides
     * a basic http client functionality.
     * If the context belongs to a logged in user
     * then all requests done via this client will have
     * automatically the current OAuth2 Access Token assigned
     */
}
```

## Configuration

Flamingo expects a `session.Store` dingo binding, this is currently handled via the `session.backend` config parameter.

Possible values are currently `file` for temporary file storage and `redis` for a redis backend.

The redis backend uses the config param `session.backend.redis.host` to find the redis, e.g. `redis.host:6379`.
