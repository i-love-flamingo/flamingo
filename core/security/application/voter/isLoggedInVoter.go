package voter

import (
	"context"

	"flamingo.me/flamingo/core/security/application/role"
	"flamingo.me/flamingo/core/security/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	IsLoggedInVoter struct {
		roleService role.Service
	}
)

func (v *IsLoggedInVoter) Inject(rs role.Service) {
	v.roleService = rs
}

func (v *IsLoggedInVoter) Vote(ctx context.Context, session *web.Session, permission string, _ interface{}) int {
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
