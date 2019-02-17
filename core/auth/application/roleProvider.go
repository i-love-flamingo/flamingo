package application

import (
	"context"

	authDomain "flamingo.me/flamingo/v3/core/auth/domain"
	securityDomain "flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// AuthRoleProvider implements the RoleProvider interface for authenticated users
	AuthRoleProvider struct {
		userService UserServiceInterface
	}
)

// Inject userService dependency
func (p *AuthRoleProvider) Inject(us UserServiceInterface) {
	p.userService = us
}

// All return all associated roles
func (p *AuthRoleProvider) All(ctx context.Context, session *web.Session) []securityDomain.Role {
	var roles []securityDomain.Role

	user := p.userService.GetUser(ctx, session)
	if user != nil && user.Type == authDomain.USER {
		roles = append(roles, authDomain.OAuthRoleUser)
	}

	return roles
}
