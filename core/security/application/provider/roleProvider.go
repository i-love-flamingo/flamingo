package provider

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/auth/application"
	authDomain "flamingo.me/flamingo/core/auth/domain"
	securityDomain "flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/config"
)

type (
	RoleProvider interface {
		All(context.Context, *sessions.Session) []securityDomain.Role
	}

	DefaultRoleProvider struct {
		userService    application.UserServiceInterface
		rolesHierarchy config.Map
	}
)

func (p *DefaultRoleProvider) Inject(us application.UserServiceInterface, cfg *struct {
	RolesHierarchy config.Map `inject:"config:security.roles.hierarchy"`
}) {
	p.userService = us
	p.rolesHierarchy = cfg.RolesHierarchy
}

func (p *DefaultRoleProvider) All(ctx context.Context, session *sessions.Session) []securityDomain.Role {
	user := p.userService.GetUser(ctx, session)
	if user == nil || user.Type != authDomain.USER {
		return p.extractRoles(securityDomain.RoleAnonymous)
	}
	return p.extractRoles(securityDomain.RoleUser)
}

func (p *DefaultRoleProvider) extractRoles(role string) []securityDomain.Role {
	roles := []securityDomain.Role{
		securityDomain.DefaultRole(role),
	}

	var roleMap map[string][]string
	p.rolesHierarchy.MapInto(&roleMap)

	hierarchy, ok := roleMap[role]
	if !ok {
		return roles
	}

	for index := range hierarchy {
		roles = append(roles, securityDomain.DefaultRole(hierarchy[index]))
	}

	return p.removeDuplicates(roles)
}

func (p *DefaultRoleProvider) removeDuplicates(roles []securityDomain.Role) []securityDomain.Role {
	roleMap := map[string]bool{}
	var clean []securityDomain.Role
	for _, role := range roles {
		if !roleMap[role.Role()] {
			clean = append(clean, role)
		}
		roleMap[role.Role()] = true
	}
	return clean
}
