package voter

import (
	"context"

	"github.com/gorilla/sessions"

	"flamingo.me/flamingo/core/security/application/provider"
	"flamingo.me/flamingo/core/security/domain"
)

type (
	IsLoggedOutVoter struct {
		roleProvider provider.RoleProvider
	}
)

func (v *IsLoggedOutVoter) Inject(rp provider.RoleProvider) {
	v.roleProvider = rp
}

func (v *IsLoggedOutVoter) Vote(ctx context.Context, session *sessions.Session, role string, _ interface{}) int {
	if role != domain.RoleAnonymous {
		return AccessAbstained
	}

	roles := v.roleProvider.All(ctx, session)
	for index := range roles {
		if roles[index].Role() == domain.RoleAnonymous {
			return AccessGranted
		}
	}

	return AccessDenied
}
