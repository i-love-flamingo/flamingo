package role

import (
	"context"

	"flamingo.me/flamingo/core/security/application/role/provider"
	"flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/web"
)

type (
	Service interface {
		All(context.Context, *web.Session) []domain.Role
	}

	ServiceImpl struct {
		providers      []provider.RoleProvider
		rolesHierarchy config.Map
	}
)

func (s *ServiceImpl) Inject(p []provider.RoleProvider, cfg *struct {
	RolesHierarchy config.Map `inject:"config:security.roles.hierarchy"`
}) {
	s.providers = p
	s.rolesHierarchy = cfg.RolesHierarchy
}

func (s *ServiceImpl) All(ctx context.Context, session *web.Session) []domain.Role {
	rolesChan := make(chan []domain.Role)

	for index := range s.providers {
		go func(p provider.RoleProvider) {
			rolesChan <- p.All(ctx, session)
		}(s.providers[index])
	}

	var roles []domain.Role
	for range s.providers {
		fromChan := <-rolesChan

		var extracted []domain.Role
		for _, role := range fromChan {
			extracted = append(extracted, s.extractRoles(role)...)
		}
		roles = append(roles, extracted...)
	}

	return s.removeDuplicates(roles)
}

func (s *ServiceImpl) extractRoles(role domain.Role) []domain.Role {
	roles := []domain.Role{
		role,
	}

	var roleMap map[string][]string
	s.rolesHierarchy.MapInto(&roleMap)

	hierarchy, ok := roleMap[role.Permission()]
	if !ok {
		return roles
	}

	for index := range hierarchy {
		roles = append(roles, domain.Role(hierarchy[index]))
	}

	return roles
}

func (s *ServiceImpl) removeDuplicates(roles []domain.Role) []domain.Role {
	roleMap := map[string]bool{}
	var clean []domain.Role
	for _, role := range roles {
		if !roleMap[role.Permission()] {
			clean = append(clean, role)
		}
		roleMap[role.Permission()] = true
	}
	return clean
}
