package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/provider"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	RoleVoter struct {
		roleProvider provider.RoleProvider
	}
)

func (v *RoleVoter) Inject(rp provider.RoleProvider) {
	v.roleProvider = rp
}

func (v *RoleVoter) Vote(ctx context.Context, session *sessions.Session, role string, object interface{}) int {
	if role == domain.RoleAnonymous || role == domain.RoleUser {
		return AccessAbstained
	}

	roleSet, ok := object.(domain.RoleSet)
	if ok && !v.hasRole(roleSet.Roles(), role) {
		return AccessDenied
	}

	roles := v.roleProvider.All(ctx, session)
	if !v.hasRole(roles, role) {
		return AccessDenied
	}
	return AccessGranted
}

func (v *RoleVoter) hasRole(roles []domain.Role, role string) bool {
	for index := range roles {
		if roles[index].Role() == role {
			return true
		}
	}

	return false
}
