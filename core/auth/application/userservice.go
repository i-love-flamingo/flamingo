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
		synchronizer   Synchronizer
		onlyOneDevice  bool
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		InitUser(ctx context.Context, session *sessions.Session)
		GetUser(ctx context.Context, session *sessions.Session) *domain.User
		IsLoggedIn(ctx context.Context, session *sessions.Session) bool
	}
)

func (us *UserService) Inject(manager *AuthManager, ums *domain.UserMappingService, s Synchronizer, cfg *struct {
	OnlyOneDevice bool `inject:"config:auth.onlyOneDevice"`
}) {
	us.authManager = manager
	us.mappingService = ums
	us.synchronizer = s
	us.onlyOneDevice = cfg.OnlyOneDevice
}

func (us *UserService) InitUser(c context.Context, session *sessions.Session) {
	user := us.getUser(c, session)

	if us.onlyOneDevice && user != nil {
		us.synchronizer.Insert(user, session)
	}
}

// GetUser returns the current user information
func (us *UserService) GetUser(c context.Context, session *sessions.Session) *domain.User {
	user := us.getUser(c, session)
	user = us.syncCheck(user, session)

	return user
}

// IsLoggedIn determines the user's login status
func (us *UserService) IsLoggedIn(c context.Context, session *sessions.Session) bool {
	user := us.GetUser(c, session)
	user = us.syncCheck(user, session)
	return user.Type == domain.USER
}

func (us *UserService) getUser(c context.Context, session *sessions.Session) *domain.User {
	id, err := us.authManager.IDToken(c, session)
	if err != nil {
		return domain.Guest
	}

	user, err := us.mappingService.UserFromIDToken(id)
	if user == nil || err != nil {
		return domain.Guest
	}

	return user
}

func (us *UserService) syncCheck(user *domain.User, session *sessions.Session) *domain.User {
	if us.onlyOneDevice && user != nil {
		isActive, err := us.synchronizer.IsActive(*user, session)
		if !isActive || err != nil {
			delete(session.Values, KeyToken)
			delete(session.Values, KeyRawIDToken)
			delete(session.Values, KeyAuthstate)
			delete(session.Values, KeyTokenExtras)
			return domain.Guest
		}
	}

	return user
}
