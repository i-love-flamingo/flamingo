package security

import (
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/core/security/interface/controller"
	"flamingo.me/flamingo/v3/core/security/interface/middleware"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module is the security module entry point
	Module struct{}

	routes struct {
		dataController *controller.DataController
	}
)

// Inject dependencies
func (r *routes) Inject(c *controller.DataController) {
	r.dataController = c
}

// Routes registers security controller
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleData("security.isLoggedIn", r.dataController.IsLoggedIn)
	registry.HandleData("security.isLoggedOut", r.dataController.IsLoggedOut)
	registry.HandleData("security.isGranted", r.dataController.IsGranted)
}

// Configure security dependency injection
func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, &routes{})

	injector.BindMulti(new(voter.SecurityVoter)).To(voter.IsLoggedInVoter{})
	injector.BindMulti(new(voter.SecurityVoter)).To(voter.PermissionVoter{})
	injector.Bind(new(role.Service)).To(role.ServiceImpl{})
	injector.Bind(new(application.SecurityService)).To(application.SecurityServiceImpl{})
	injector.Bind(new(middleware.RedirectURLMaker)).To(middleware.RedirectURLMakerImpl{})
}

// CueConfig schema
func (*Module) CueConfig() string {
	return fmt.Sprintf(`
core: security: {
	loginPath: {
		handler: string | *"auth.login"
		redirectStrategy: string | *"%s"
		redirectPath: string | *"/"
	}
	authenticatedHomepage:{
		strategy: string | *"%s"
		path: string | *"/"
	}
	roles: {
		permissionHierarchy: {
			%s: [..._]
		}
		voters: {
			strategy: string | *"%s"
			allowIfAllAbstain: bool | *false
		}
	}
	eventLogging: bool | *false
}
`, middleware.ReferrerRedirectStrategy, middleware.ReferrerRedirectStrategy, domain.PermissionAuthorized, application.VoterStrategyAffirmative)
}

// FlamingoLegacyConfigAlias mapping for legacy settings
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	alias := make(map[string]string)
	for _, v := range []string{
		"security.loginPath.handler",
		"security.loginPath.redirectStrategy",
		"security.loginPath.redirectPath",
		"security.authenticatedHomepage.strategy",
		"security.authenticatedHomepage.path",
		"security.roles.permissionHierarchy",
		"security.roles.voters.strategy",
		"security.roles.voters.allowIfAllAbstain",
		"security.eventLogging",
	} {
		alias[v] = "core." + v
	}
	return alias
}
