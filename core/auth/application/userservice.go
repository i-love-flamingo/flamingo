package application

import (
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	UserService struct {
		AuthManager *AuthManager `inject:""`
	}
)

type UserServiceInterface interface {
	GetUser(web.Context) *domain.User
	IsLoggedIn(web.Context) bool
}

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
