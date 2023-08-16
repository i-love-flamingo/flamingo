# Security module

The security module provides a way to check and handle permissions for actions inside a web application that used a regular session.

The implemented security concept uses "Roles". Roles are returned by "RoleProviders".
A "Role" need to return "Permission". Based on the "Permissions" Security Voters decide if access should be granted or not.

That module can be used with the provided middleware during binding of routes.
It can of course also used by the provided "SecurityService".

The "auth" module provides a RoleProvider - so this security module can be used together with the auth nodule.

You can register your very own "RoleProvider" if required.

## Security Middleware

To add security middleware, simply inject it into routes struct and use it
as a wrapper for any handler.

There are three possibilities:

* HandleIfLoggedOut - will forward to handler only if the user is not logged in. Otherwise,
it will return a redirect response either to the homepage for authenticated users or to
http referrer, depending on configuration (see below).
* HandleIfLoggedIn - will forward to handler only if the user is logged in. Otherwise,
it will return a redirect response to the login page with redirect url either to specific path
or to requested path, depending on configuration (see below).
* HandleIfGranted - will forward to handler only if a specific permission is granted to the user. 
Otherwise, it will return a 403 page.

```go
type routes struct {
  someController      *controller.SomeController
  securityMiddleware  *middleware.SecurityMiddleware
}

func (r *routes) Routes(registry *router.Registry) {
  registry.HandleGet("register", r.securityMiddleware.HandleIfLoggedOut(r.someController.Register))
  registry.HandleGet("my.account", r.securityMiddleware.HandleIfLoggedIn(r.someController.MyAccount))
  registry.HandleGet("users.list", r.securityMiddleware.HandleIfGranted(r.someController.Users, "PermissionAdmin"))
  registry.HandleGet("users.list", r.securityMiddleware.HandleIfNotGranted(r.someController.Users, "PermissionSuperAdmin"))
}
```

## Security Service

Security service provides more detailed security checks. Beside checking if the user is 
logged in or not, it's possible to check for a permission.
If `IsGranted` is called with no object, it will check only permissions added to the
user in session. Otherwise, if the item implements security `domain.PermissionSet` interface, it will
check if the permission is attached to the item itself.

```go
type SomeController struct {
  securityService  application.SecurityService
  repository       *SomeRepository
}

func (c *SomeController) Handle(ctx context.Context, r *web.Request) web.Result {
  if c.securityService.IsLoggedIn(ctx, r.Session().G()) {
    ...
  }
  if c.securityService.IsLoggedOut(ctx, r.Session().G()) {
    ...
  }
  item := c.repository.Get(1)
  if c.securityService.IsGranted(ctx, r.Session().G(), "PermissionEdit", item) {
    ...
  }
  if c.securityService.IsGranted(ctx, r.Session().G(), "PermissionSave", nil) {
    ...
  }
}
```

PermissionSet interface:
```go
type PermissionSet interface {
  Permissions() []Role
}
```

## Security Voters

Security voters are used inside security service to decide on final access permissions for the user
and on the object. By default, there are three security voters: `IsLoggedInVoter`, `IsLoggedOutVoter`, `RoleVoter`.
Voters can return three different status (defined as constants):
* `AccessAbstained` - vote is not relevant
* `AccessDeny` - voter declares permission for resource as denied
* `AccessGranted` - voter declares permission for resource as granted

After all voters vote, final decision is made on the voters strategy (see configuration below):
* `affirmative` - grants access if there is at least one voter which grants access
* `consensus` - grants access if there are more voters which grant access then voters which deny
* `unanimous` - grants access if there is at least one voter which grants access and none which denies
If there is no clear decision between voters, a final mediator is defined by configuration parameter "allowIfAllAbstain"

To add a new voter simply implement the `SecurityVoter` interface and inject it into the list of Security Voters:

```go
func (m *Module) Configure(injector *dingo.Injector) {
  injector.BindMulti(new(voter.SecurityVoter))).To(voter.CustomVoter{})
}
```

SecurityVoter interface:
```go
type SecurityVoter interface {
  Vote(allPermissions []string, desiredPermisssion string, forObject interface{}) int
}
```

## Roles Providers

Role providers are used to fetch all roles granted for the user in a session.
To provide more roles it's possible to define additional role providers:

```go
func (m *Module) Configure(injector *dingo.Injector) {
  injector.BindMulti(new(role.Provider))).To(provider.CustomProvider{})
}
```

Role Provider interface:
```go
type Provider interface {
  All(context.Context, *sessions.Session) []domain.Role
}
```

## Configuration example

Permissions hierarchy provides automatic inclusion of child permissions into list of permissions if 
their parent is fetched via role providers.

```yaml
security: 
    login:
        handler: "auth.log"
        redirectStrategy: "path" # possible referrer|path
        redirectPath: "/" #only if strategy is "path"
    authenticatedHomepage:
        strategy: "path" # possible referrer|path
        path: "/" #only if strategy is "path"
    roles:
        permissionHierarchy:
            PermissionAuthorized:
                - PermissionView
            PermissionAdmin:
                - PermissionView
                - PermissionEdit
            PermissionSuperAdmin:
                - PermissionView
                - PermissionEdit
                - PermissionDelete
        voters:
            strategy: "unanimous" # possible unanimous|affirmative|consensus
            allowIfAllAbstain: false
```
