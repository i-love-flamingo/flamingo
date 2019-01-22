package security

import (
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/core/security/application/role/provider"
	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/core/security/interface/controller"
	"flamingo.me/flamingo/v3/core/security/interface/middleware"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
)

type (
	Module struct{}

	routes struct {
		dataController *controller.DataController
	}
)

func (r *routes) Inject(c *controller.DataController) {
	r.dataController = c
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleData("security.isLoggedIn", r.dataController.IsLoggedIn)
	registry.HandleData("security.isLoggedOut", r.dataController.IsLoggedOut)
	registry.HandleData("security.isGranted", r.dataController.IsGranted)
}

func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, &routes{})

	injector.BindMulti((*provider.RoleProvider)(nil)).To(provider.AuthRoleProvider{})
	injector.BindMulti((*voter.SecurityVoter)(nil)).To(voter.IsLoggedInVoter{})
	injector.BindMulti((*voter.SecurityVoter)(nil)).To(voter.RoleVoter{})
	injector.Bind((*role.Service)(nil)).To(role.ServiceImpl{})
	injector.Bind((*application.SecurityService)(nil)).To(application.SecurityServiceImpl{})
	injector.Bind((*middleware.RedirectUrlMaker)(nil)).To(authApplication.AuthManager{})
}

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
				"hierarchy": config.Map{
					domain.RoleUser.Permission(): config.Slice{},
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
