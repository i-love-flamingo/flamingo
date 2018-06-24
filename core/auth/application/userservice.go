package application

import (
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// userService helps to use the authenticated user information
	UserService struct {
		authManager *AuthManager
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		GetUser(web.Context) *domain.User
		IsLoggedIn(web.Context) bool
	}
)

func (us *UserService) Inject(manager *AuthManager) {
	us.authManager = manager
}

// GetUser returns the current user information
func (us *UserService) GetUser(c web.Context) *domain.User {
	id, err := us.authManager.IDToken(c)
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
