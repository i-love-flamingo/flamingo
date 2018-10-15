package application

import (
	"context"

	"flamingo.me/flamingo/core/auth/domain"
	"github.com/gorilla/sessions"
)

type (
	// userService helps to use the authenticated user information
	UserService struct {
		authManager    *AuthManager
		mappingService *domain.UserMappingService
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		GetUser(ctx context.Context, session *sessions.Session) *domain.User
		IsLoggedIn(ctx context.Context, session *sessions.Session) bool
	}
)

func (us *UserService) Inject(manager *AuthManager, ums *domain.UserMappingService) {
	us.authManager = manager
	us.mappingService = ums
}

// GetUser returns the current user information
func (us *UserService) GetUser(c context.Context, session *sessions.Session) *domain.User {
	id, err := us.authManager.IDToken(c, session)
	if err != nil {
		return domain.Guest
	}

	user, err := us.mappingService.UserFromIDToken(id)
	if err != nil {
		return domain.Guest
	}

	return user
}

// IsLoggedIn determines the user's login status
func (us *UserService) IsLoggedIn(c context.Context, session *sessions.Session) bool {
	user := us.GetUser(c, session)
	return user.Type == domain.USER
}
