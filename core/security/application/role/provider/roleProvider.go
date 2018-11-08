package provider

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/auth/application"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	securityDomain "flamingo.me/flamingo/core/security/domain"
)

type (
	RoleProvider interface {
		All(context.Context, *sessions.Session) []securityDomain.Role
	}

	DefaultRoleProvider struct {
		userService application.UserServiceInterface
	}
)

func (p *DefaultRoleProvider) Inject(us application.UserServiceInterface) {
	p.userService = us
}

func (p *DefaultRoleProvider) All(ctx context.Context, session *sessions.Session) []securityDomain.Role {
	user := p.userService.GetUser(ctx, session)
	if user == nil || user.Type != authDomain.USER {
		return []securityDomain.Role{
			securityDomain.RoleAnonymous,
		}
	}
	return []securityDomain.Role{
		securityDomain.RoleUser,
	}
}
