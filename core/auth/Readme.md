# Auth Modul

The purpose of the auth modul is to provide a standard way to receive and check for an Idendity in the current request (or its session).

This is provided as a secondary port and therefore it is up to the implementation to decide what kind of Idendity is returned and also how the Authentication process should work.
 
## Usage

Any other module that need to have an Idendity can simple inject am "Authinfo" application service. And use the Idendity

e.g.:

```go

if !authinfo.IsAuthenticated(ctx,r) {
  return nil, errors.New("no idendity given")
  //Or start auth flow  return authinfo.Authenticate(ctx,myUrl)
}

idendity, err := authinfo.GetIdendity(ctx,r)
if err != nil {
	return nil, err
}
myCustomerService.GetCustomerForIdendity(ctx, idendity)
```
