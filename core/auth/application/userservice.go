package application

import (
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// UserService helps to use the authenticated user information
	UserService struct {
		AuthManager *AuthManager `inject:""`
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		GetUser(web.Context) *domain.User
		IsLoggedIn(web.Context) bool
	}
)

// GetUser returns the current user information
func (us *UserService) GetUser(c web.Context) *domain.User {
	id, err := us.AuthManager.IDToken(c)
	if err != nil {
		return domain.Guest
	}

	return domain.UserFromIDToken(id)
}

// IsLoggedIn determines the user's login status
func (us *UserService) IsLoggedIn(c web.Context) bool {
	user := us.GetUser(c)
	return user.Type == domain.USER
}
