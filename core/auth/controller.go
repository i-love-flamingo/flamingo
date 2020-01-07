package auth

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

//type Middleware struct {
//	wis *WebIdentityService
//}
//
//func (wd *Middleware) RequireIdentified(action web.Action) web.Action {
//	return func(ctx context.Context, req *web.Request) web.Result {
//		if id := wd.wis.Identify(ctx, req); id == nil {
//			_, res := wd.wis.Authenticate(ctx, req)
//			return res
//		}
//		return action(ctx, req)
//	}
//}
//
//func (wd *Middleware) RequireIdentifiedFor(broker string, action web.Action) web.Action {
//	return func(ctx context.Context, req *web.Request) web.Result {
//		if wd.wis.IdentifyFor(broker, ctx, req) == nil {
//			return wd.wis.AuthenticateFor(broker, ctx, req)
//		}
//		return action(ctx, req)
//	}
//}

//func (wd *Middleware) RequirePermission(permission string, action web.Action) web.Action {
//	return func(ctx context.Context, req *web.Request) web.Result {
//		return action(ctx, req)
//	}
//}

type Controller struct {
	service *WebIdentityService
}

func (c *Controller) Inject(service *WebIdentityService) {
	c.service = service
}

func (c *Controller) Callback(ctx context.Context, request *web.Request) web.Result {
	return c.service.callback(ctx, request)
}
