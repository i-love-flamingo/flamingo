package voter

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/application/role"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// IsLoggedInVoter votes for users who have authenticated
type IsLoggedInVoter struct {
	roleService role.Service
}

// Inject roleService dependency
func (v *IsLoggedInVoter) Inject(rs role.Service) {
	v.roleService = rs
}

// Vote for the authentication request
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
