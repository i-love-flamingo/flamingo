package healthcheck

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/core/auth/application"
)

// Auth healthcheck
type Auth struct {
	authManager *application.AuthManager
}

var _ Status = &Auth{}

// Inject auth manager dependency
func (s *Auth) Inject(authManager *application.AuthManager) {
	s.authManager = authManager
}

// Status checks the status
func (s *Auth) Status() (bool, string) {
	path := s.authManager.OAuth2Config(context.Background(), nil).AuthCodeURL("")
	_, err := http.Get(path)
	if err == nil {
		return true, "success"
	}

	return false, err.Error()
}
