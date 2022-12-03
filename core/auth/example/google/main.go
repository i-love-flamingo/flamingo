package main

import (
	"context"
	"net/http"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/web"
)

func main() {
	flamingo.App([]dingo.Module{
		new(oauth.Module),
		new(testModule),
	})
}

type testModule struct{}

func (*testModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	testController *testController
}

func (r *routes) Inject(controller *testController) *routes {
	r.testController = controller
	return r
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/", "index")
	registry.HandleAny("index", r.testController.Index)
}

type testController struct {
	responder          *web.Responder
	webIdentityService *auth.WebIdentityService
}

// Inject dependencies
func (controller *testController) Inject(responder *web.Responder, webIdentityService *auth.WebIdentityService) *testController {
	controller.responder = responder
	controller.webIdentityService = webIdentityService

	return controller
}

func (controller *testController) Index(ctx context.Context, req *web.Request) web.Result {
	identity := controller.webIdentityService.Identify(ctx, req)
	body := ""
	if identity == nil {
		body = "Hello Guest"
	} else {
		oidcIdentity, _ := identity.(oauth.OpenIDIdentity)
		body = "Hello " + oidcIdentity.Subject() + " IDToken Expired: " + oidcIdentity.IDToken().Expiry.GoString()
	}

	return controller.responder.HTTP(
		http.StatusOK,
		strings.NewReader(body),
	)

}
