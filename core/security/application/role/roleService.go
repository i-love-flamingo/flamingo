package role

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Provider interface
	Provider interface {
		All(context.Context, *web.Session) []domain.Role
	}

	// Service implements behaviour for retrieving permissions by using all role providers
	Service interface {
		AllPermissions(context.Context, *web.Session) []string
	}

	// The ServiceImpl is the default Service implementation
	ServiceImpl struct {
		providers           []Provider
		permissionHierarchy map[string][]string
	}
)

// Inject dependencies
func (s *ServiceImpl) Inject(p []Provider, cfg *struct {
	PermissionHierarchy config.Map `inject:"config:security.roles.permissionHierarchy"`
}) {
	s.providers = p

	var permissionHierarchy map[string][]string
	err := cfg.PermissionHierarchy.MapInto(&permissionHierarchy)
	if err != nil {
		panic(err)
	}

	s.permissionHierarchy = permissionHierarchy
}

// AllPermissions returns all available permissions, based on their hierarchy
func (s *ServiceImpl) AllPermissions(ctx context.Context, session *web.Session) []string {
	rolesChan := make(chan []domain.Role)

	for index := range s.providers {
		go func(p Provider) {
			rolesChan <- p.All(ctx, session)
		}(s.providers[index])
	}

	var permissions []string
	for range s.providers {
		fromChan := <-rolesChan

		var extracted []string
		for _, role := range fromChan {
			extracted = append(extracted, s.extractPermissions(role)...)
		}
		permissions = append(permissions, extracted...)
	}

	return s.removeDuplicates(permissions)
}

func (s *ServiceImpl) extractPermissions(role domain.Role) []string {
	permissions := role.Permissions()

	for _, permission := range role.Permissions() {
		hierarchy, ok := s.permissionHierarchy[permission]
		if !ok {
			continue
		}

		permissions = append(permissions, hierarchy...)
	}

	return permissions
}

func (s *ServiceImpl) removeDuplicates(permissions []string) []string {
	permissionMap := map[string]bool{}
	var cleanPermissions []string

	for _, permission := range permissions {
		if !permissionMap[permission] {
			cleanPermissions = append(cleanPermissions, permission)
		}
		permissionMap[permission] = true
	}

	return cleanPermissions
}
