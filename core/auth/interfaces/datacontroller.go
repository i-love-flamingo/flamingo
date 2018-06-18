package interfaces

import (
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	// UserController uc
	UserController struct {
		UserService *application.UserService `inject:""`
	}
)

// Data controller to return userinfo
func (u *UserController) Data(c web.Context) interface{} {
	return u.UserService.GetUser(c)
}
