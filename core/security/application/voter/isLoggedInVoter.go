package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/provider"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	IsLoggedInVoter struct {
		roleProvider provider.RoleProvider
	}
)

func (v *IsLoggedInVoter) Inject(rp provider.RoleProvider) {
	v.roleProvider = rp
}

func (v *IsLoggedInVoter) Vote(ctx context.Context, session *sessions.Session, role string, _ interface{}) int {
	if role != domain.RoleUser {
		return AccessAbstained
	}

	roles := v.roleProvider.All(ctx, session)
	for index := range roles {
		if roles[index].Role() == domain.RoleUser {
			return AccessGranted
		}
	}

	return AccessDenied
}
