package voter

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	RoleVoter struct {
		roleService role.Service
	}
)

func (v *RoleVoter) Inject(rs role.Service) {
	v.roleService = rs
}

func (v *RoleVoter) Vote(ctx context.Context, session *web.Session, permission string, object interface{}) int {
	if permission == domain.RoleUser.Permission() {
		return AccessAbstained
	}

	roleSet, ok := object.(domain.RoleSet)
	if ok && !v.hasPermission(roleSet.Roles(), permission) {
		return AccessDenied
	}

	roles := v.roleService.All(ctx, session)
	if !v.hasPermission(roles, permission) {
		return AccessDenied
	}
	return AccessGranted
}

func (v *RoleVoter) hasPermission(roles []domain.Role, permission string) bool {
	for index := range roles {
		if roles[index].Permission() == permission {
			return true
		}
	}

	return false
}
