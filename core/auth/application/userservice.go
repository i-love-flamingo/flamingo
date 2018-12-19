package application

import (
	"context"

	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
	"github.com/gorilla/sessions"
)

type (
	// userService helps to use the authenticated user information
	UserService struct {
		authManager                 *AuthManager
		mappingService              *domain.UserMappingService
		synchronizer                Synchronizer
		preventSimultaneousSessions bool
	}

	// UserServiceInterface to mock in tests
	UserServiceInterface interface {
		InitUser(ctx context.Context, session *sessions.Session) error
		GetUser(ctx context.Context, session *sessions.Session) *domain.User
		IsLoggedIn(ctx context.Context, session *sessions.Session) bool
	}
)

func (us *UserService) Inject(manager *AuthManager, ums *domain.UserMappingService, s Synchronizer, cfg *struct {
	PreventSimultaneousSessions bool `inject:"config:auth.preventSimultaneousSessions"`
}) {
	us.authManager = manager
	us.mappingService = ums
	us.synchronizer = s
	us.preventSimultaneousSessions = cfg.PreventSimultaneousSessions
}

func (us *UserService) InitUser(c context.Context, session *sessions.Session) error {
	user := us.getUser(c, session)

	if us.preventSimultaneousSessions && user != nil && user.Type != domain.GUEST {
		return us.synchronizer.Insert(*user, session)
	}

	return nil
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

	r, _ := web.FromContext(c)
	user, err := us.mappingService.UserFromIDToken(id, r.Session())
	if user == nil || err != nil {
		return domain.Guest
	}

	return user
}

func (us *UserService) syncCheck(user *domain.User, session *sessions.Session) *domain.User {
	if us.preventSimultaneousSessions && user != nil && user.Type != domain.GUEST {
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
