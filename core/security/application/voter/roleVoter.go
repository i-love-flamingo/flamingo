package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/role"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	RoleVoter struct {
		roleService role.Service
	}
)

func (v *RoleVoter) Inject(rs role.Service) {
	v.roleService = rs
}

func (v *RoleVoter) Vote(ctx context.Context, session *sessions.Session, permission string, object interface{}) int {
	if permission == domain.RoleAnonymous.Permission() || permission == domain.RoleUser.Permission() {
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
