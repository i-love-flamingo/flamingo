package provider

import (
	"context"

	"flamingo.me/flamingo/core/auth/application"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	securityDomain "flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	RoleProvider interface {
		All(context.Context, *web.Session) []securityDomain.Role
	}

	AuthRoleProvider struct {
		userService application.UserServiceInterface
	}
)

func (p *AuthRoleProvider) Inject(us application.UserServiceInterface) {
	p.userService = us
}

func (p *AuthRoleProvider) All(ctx context.Context, session *web.Session) []securityDomain.Role {
	var roles []securityDomain.Role

	user := p.userService.GetUser(ctx, session.G())
	if user != nil && user.Type == authDomain.USER {
		roles = append(roles, securityDomain.RoleUser)
	}

	return roles
}
