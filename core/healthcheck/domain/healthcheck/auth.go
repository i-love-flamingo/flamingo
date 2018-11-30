package healthcheck

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/core/auth/application"
)

type (
	Auth struct {
		authManager *application.AuthManager
	}
)

var (
	_ Status = &Auth{}
)

func (s *Auth) Inject(authManager *application.AuthManager) {
	s.authManager = authManager
}

func (s *Auth) Status() (bool, string) {
	path := s.authManager.OAuth2Config(context.Background()).AuthCodeURL("")
	_, err := http.Get(path)
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}
