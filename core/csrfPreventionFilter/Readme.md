# CSRF Prevention Filter

This package provides middleware for CSRF security prevention (Cross-Site Request Forgery).

## Configuration

By default, there are three parameters supported by module:
* "all" - defines if all POST forms should be CSRF secured (default is false)
* "secret" - defines key for AES encryption (16, 24 or 32 bytes for AES-128, AES-192 or AES-256)
* "ttl" - defines max time (in seconds) validation for some token (default is 15 minutes)

```
csrf:
  all: false
  secret: 6368616e676520746869732070617373776f726420746f206120736563726574
  ttl: 900
```

## Specific form

In case when it's not required to secure all forms, it's possible just to put
middleware just for particular handler. In that case, only POST request for that handler will
be secured.

```go
type (
  routes struct {
    someController  *controller.SomeController
    csrfMiddleware  *interfaces.CsrfMiddleware
  }
)

func (r *routes) Routes(registry *router.Registry) {
  registry.HandlePost("some.handler", r.csrfMiddleware.Secured(r.someController.Handler))
}
```

## Template

To add hidden input token into template, use template function:

```
!= csrfInput()
```

To add just token into template, use template function:

```
input(type="hidden" name="csrftoken" value=csrfToken())
```
