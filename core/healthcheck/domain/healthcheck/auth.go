package healthcheck

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/core/auth/application"
)

type (
	// Auth is the healthcheck for auth module
	Auth struct {
		authManager *application.AuthManager
	}
)

var (
	_ Status = &Auth{}
)

// Inject dependencies
func (s *Auth) Inject(authManager *application.AuthManager) {
	s.authManager = authManager
}

// Status returns the health state of auth manager
func (s *Auth) Status() (bool, string) {
	path := s.authManager.OAuth2Config(context.Background()).AuthCodeURL("")
	_, err := http.Get(path)
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}
