package application

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// UserService helps to use the authenticated user information
	UserService struct {
		authManager    *AuthManager
		mappingService *domain.UserMappingService
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		GetUser(ctx context.Context, session *web.Session) *domain.User
		IsLoggedIn(ctx context.Context, session *web.Session) bool
	}
)

// Inject dependencies
func (us *UserService) Inject(manager *AuthManager, ums *domain.UserMappingService) {
	us.authManager = manager
	us.mappingService = ums
}

// GetUser returns the current user information
func (us *UserService) GetUser(c context.Context, session *web.Session) *domain.User {
	user := us.getUser(c, session)

	return user
}

// IsLoggedIn determines the user's login status
func (us *UserService) IsLoggedIn(c context.Context, session *web.Session) bool {
	user := us.getUser(c, session)
	return user.Type == domain.USER
}

func (us *UserService) getUser(c context.Context, session *web.Session) *domain.User {
	id, err := us.authManager.IDToken(c, session)
	if err != nil {
		return domain.Guest
	}

	r := web.RequestFromContext(c)
	user, err := us.mappingService.UserFromIDToken(id, r.Session())
	if user == nil || err != nil {
		return domain.Guest
	}

	return user
}
