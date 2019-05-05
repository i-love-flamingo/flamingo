# Systemendpoint module

This package provides a second endpoint intended for non public use.

You can register simple `http.Handler` to a desired route via Dingo map binding:

```go
injector.BindMap((*domain.Handler)(nil), "/my/route").To(&myHandler{})
```

This module will then bring up an HTTP Server at the configured address `systemendpoint.serviceAddr` 
which defaults to `:13210` serving all bound routes.

The server will be started on `flamingo.ServerStartEvent` and shut down on `flamingo.ServerShutdownEvent`.
