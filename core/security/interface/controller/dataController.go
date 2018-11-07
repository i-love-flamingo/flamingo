package controller

import (
	"context"

	"flamingo.me/flamingo/core/security/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	DataController struct {
		securityService application.SecurityService
	}
)

func (c *DataController) Inject(s application.SecurityService) {
	c.securityService = s
}

func (c *DataController) IsLoggedIn(ctx context.Context, r *web.Request) interface{} {
	return c.securityService.IsLoggedIn(ctx, r.Session().G())
}

func (c *DataController) IsLoggedOut(ctx context.Context, r *web.Request) interface{} {
	return c.securityService.IsLoggedOut(ctx, r.Session().G())
}

func (c *DataController) IsGranted(ctx context.Context, r *web.Request) interface{} {
	role := r.MustParam1("role")
	return c.securityService.IsGranted(ctx, r.Session().G(), role, nil)
}
