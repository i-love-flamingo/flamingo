package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/role"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	IsLoggedInVoter struct {
		roleService role.Service
	}
)

func (v *IsLoggedInVoter) Inject(rs role.Service) {
	v.roleService = rs
}

func (v *IsLoggedInVoter) Vote(ctx context.Context, session *sessions.Session, permission string, _ interface{}) int {
	if permission != domain.RoleUser.Permission() {
		return AccessAbstained
	}

	roles := v.roleService.All(ctx, session)
	for index := range roles {
		if roles[index].Permission() == domain.RoleUser.Permission() {
			return AccessGranted
		}
	}

	return AccessDenied
}
