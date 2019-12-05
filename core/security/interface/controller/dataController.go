package controller

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/framework/web"
)

// DataController returns helper for checking access
type DataController struct {
	securityService application.SecurityService
}

// Inject security service dependency
func (c *DataController) Inject(s application.SecurityService) {
	c.securityService = s
}

// IsLoggedIn check
func (c *DataController) IsLoggedIn(ctx context.Context, r *web.Request, _ web.RequestParams) interface{} {
	return c.securityService.IsLoggedIn(ctx, r.Session())
}

// IsLoggedOut check
func (c *DataController) IsLoggedOut(ctx context.Context, r *web.Request, _ web.RequestParams) interface{} {
	return c.securityService.IsLoggedOut(ctx, r.Session())
}

// IsGranted permission check
func (c *DataController) IsGranted(ctx context.Context, r *web.Request, params web.RequestParams) interface{} {
	permission := params["permission"]
	return c.securityService.IsGranted(ctx, r.Session(), permission, nil)
}
