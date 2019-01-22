package role

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/application/role/provider"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Service implements behaviour for retrieving roles
	Service interface {
		All(context.Context, *web.Session) []domain.Role
	}

	// The ServiceImpl is the default Service implementation
	ServiceImpl struct {
		providers      []provider.RoleProvider
		rolesHierarchy config.Map
	}
)

// Inject dependencies
func (s *ServiceImpl) Inject(p []provider.RoleProvider, cfg *struct {
	RolesHierarchy config.Map `inject:"config:security.roles.hierarchy"`
}) {
	s.providers = p
	s.rolesHierarchy = cfg.RolesHierarchy
}

// All returns all available roles, based on their hierarchy
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
