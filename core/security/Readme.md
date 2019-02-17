# Security module

Used to handle permissions for handlers defined in controllers.

## Security Middleware

To add security middleware, simply inject it into routes struct and use it
as a wrapper for any handler.
There are three possibilities:
* HandleIfLoggedOut - will forward to handler only if user is not logged in. Otherwise,
it will return redirect response either to homepage for authenticated users, either to
http referrer, depending on configuration (see bellow).
* HandleIfLoggedIn - will forward to handler only if user is logged in. Otherwise,
it will return redirect response to login page with redirect url either to specific path, 
either to requested path, depending on configuration (see bellow).
* HandleIfGranted - will forward to handler only if specific permission is granted to user. 
Otherwise, it will return 403 page.

```
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

Security service provides more detailed security checks. Beside checking if user is 
logged in or not, it's possible to check for permission.
If IsGranted method is called with no object, it will check only permissions added to
user in session. Otherwise, if item implements security domain.PermissionSet interface, it will
check if permission is attached to item itself.

```
type SomeController struct {
  securityService  application.SecurityService
  repository       *SomeRepository
}

func (c *SomeController) Handle(ctx context.Context, r *web.Request) web.Response {
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
```
type PermissionSet interface {
  Permissions() []Role
}
```

## Security Voters

Security voters are used inside Security Service to decide on final access permissions for user
and on object. By default, there are three Security Voters: IsLoggedInVoter, IsLoggedOutVoter, RoleVoter.
Voters can return three different statuses (defined as constants):
* AccessAbstained - vote is not relevant
* AccessDeny - voter declares permission for resource as denied
* AccessGranted - voter declares permission for resource as granted

After all voters vote, final decision is made on voters strategy (see configuration bellow):
* affirmative - grants access if there is at least one voter which grants access
* consensus - grants access if there is at more voters which grant access then ones which deny
* unanimous - grants access if there is at least one voter which grants access and none which denies
If there is no clear decision between voters, final mediator is defined by configuration
parameter "allowIfAllAbstain"

To add new voter simply provide struct which implements SecurityVoter interface and inject
it into list of Security Voters:

```
func (m *Module) Configure(injector *dingo.Injector) {
  injector.BindMulti(new(voter.SecurityVoter))).To(voter.CustomVoter{})
}
```

SecurityVoter interface:
```
type SecurityVoter interface {
  Vote(allPermissions []string, desiredPermisssion string, forObject interface{}) int
}
```

## Roles Providers

Role providers are used to fetch all roles granted for user in session.
To provide more roles it's possible to define additional Role Provider:

```
func (m *Module) Configure(injector *dingo.Injector) {
  injector.BindMulti(new(role.Provider))).To(provider.CustomProvider{})
}
```

Role Provider interface:
```
type Provider interface {
  All(context.Context, *sessions.Session) []domain.Role
}
```

## Configuration example

Permissions hierarchy provides automatic inclusion of child permissions into list of permissions if 
their parent is fetched via Role Providers.

```
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
