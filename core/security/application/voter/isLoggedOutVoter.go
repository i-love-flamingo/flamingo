package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/role"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	IsLoggedOutVoter struct {
		roleService role.Service
	}
)

func (v *IsLoggedOutVoter) Inject(rs role.Service) {
	v.roleService = rs
}

func (v *IsLoggedOutVoter) Vote(ctx context.Context, session *sessions.Session, permission string, _ interface{}) int {
	if permission != domain.RoleAnonymous.Permission() {
		return AccessAbstained
	}

	roles := v.roleService.All(ctx, session)
	for index := range roles {
		if roles[index].Permission() == domain.RoleAnonymous.Permission() {
			return AccessGranted
		}
	}

	return AccessDenied
}
