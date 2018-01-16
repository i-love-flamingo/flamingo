package application

import (
	"go.aoe.com/flamingo/core/auth/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	UserService struct {
		AuthManager *AuthManager `inject:""`
	}
)

func (us *UserService) GetUser(c web.Context) *domain.User {
	id, err := us.AuthManager.IDToken(c)
	if err != nil {
		return domain.Guest
	}

	return domain.UserFromIDToken(id)
}

func (us *UserService) IsLoggedIn(c web.Context) bool {
	user := us.GetUser(c)
	return user.Type == domain.USER
}
