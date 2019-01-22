package fake

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

const (
	// UserSessionKey for setting fake users
	UserSessionKey = "auth.fake.user"
)

type (
	// UserService is a fake type to support integration tests
	UserService struct{}
)

// InitUser satisfies the interface but does not initialize anything
func (us *UserService) InitUser(c context.Context, session *web.Session) error {
	return nil
}

// GetUser returns the user form the session
func (us *UserService) GetUser(ctx context.Context, session *web.Session) *domain.User {
	value, _ := session.Load(UserSessionKey)
	user, ok := value.(domain.User)
	if !ok {
		return domain.Guest
	}

	return &user
}

// IsLoggedIn returns true if there is a User set
func (us *UserService) IsLoggedIn(c context.Context, session *web.Session) bool {
	user := us.GetUser(c, session)
	return user.Type == domain.USER
}
