# Systemendpoint module

This package provides a second endpoint intended for internal use.

You can register simple `http.Handler` to a desired route via dingo map binding:

```go
injector.BindMap((*domain.Handler)(nil), "/my/route").To(&myHandler{})
```

This module will then bring up an HTTP Server at the configured address `systemendpoint.serviceAddr` 
which defaults to `:13210` serving all bound routes.

The server will be started on `flamingo.AppStartupEvent` and shut down on `flamingo.AppShutdownEvent`.
