package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// UserController uc
	UserController struct {
		userService application.UserServiceInterface
	}
)

// Inject UserController dependencies
func (u *UserController) Inject(service application.UserServiceInterface) {
	u.userService = service
}

// Data controller to return userinfo
func (u *UserController) Data(c context.Context, r *web.Request) interface{} {
	return u.userService.GetUser(c, r.Session().G())
}
