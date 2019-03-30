package security

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/core/security/interface/controller"
	"flamingo.me/flamingo/v3/core/security/interface/middleware"
	"flamingo.me/flamingo/v3/framework/config"
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

// DefaultConfig for core security module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"security": config.Map{
			"loginPath": config.Map{
				"handler":          "auth.login",
				"redirectStrategy": middleware.ReferrerRedirectStrategy,
				"redirectPath":     "/",
			},
			"authenticatedHomepage": config.Map{
				"strategy": middleware.ReferrerRedirectStrategy,
				"path":     "/",
			},
			"roles": config.Map{
				"permissionHierarchy": config.Map{
					domain.PermissionAuthorized: config.Slice{},
				},
				"voters": config.Map{
					"strategy":          application.VoterStrategyAffirmative,
					"allowIfAllAbstain": false,
				},
			},
			"eventLogging": false,
		},
	}
}
