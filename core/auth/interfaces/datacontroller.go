package interfaces

import (
	"go.aoe.com/flamingo/core/auth/application"
	"go.aoe.com/flamingo/core/auth/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// UserController uc
	UserController struct {
		AuthManager *application.AuthManager `inject:""`
	}
)

// Data controller to return userinfo
func (u *UserController) Data(c web.Context) interface{} {
	id, err := u.AuthManager.IDToken(c)
	if err != nil {
		return domain.Guest
	}

	return domain.UserFromIDToken(id)
}
