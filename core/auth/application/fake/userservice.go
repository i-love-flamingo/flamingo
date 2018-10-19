package fake

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/auth/domain"
)

const (
	UserSessionKey = "auth.fake.user"
)

type (
	UserService struct {

	}
)

func (u *UserService) GetUser(ctx context.Context, session *sessions.Session) *domain.User {
	value := session.Values[UserSessionKey]
	user, ok := value.(domain.User)
	if !ok {
		return domain.Guest
	}

	return &user
}

func (u *UserService) IsLoggedIn(c context.Context, session *sessions.Session) bool {
	user := u.GetUser(c, session)
	return user.Type == domain.USER
}
