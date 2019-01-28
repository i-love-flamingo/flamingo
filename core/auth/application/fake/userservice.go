package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/core/auth/domain"
)

const (
	UserSessionKey = "auth.fake.user"
)

type (
	UserService struct {}
)

func (us *UserService) InitUser(c context.Context, session *web.Session) error {
	return nil
}

func (u *UserService) GetUser(ctx context.Context, session *web.Session) *domain.User {
	value, _ := session.Load(UserSessionKey)
	user, ok := value.(domain.User)
	if !ok {
		return domain.Guest
	}

	return &user
}

func (u *UserService) IsLoggedIn(c context.Context, session *web.Session) bool {
	user := u.GetUser(c, session)
	return user.Type == domain.USER
}
