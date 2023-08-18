package auth

import (
	"context"

	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Identity donates an authentication object which at least identifies the authenticated subject
	Identity interface {
		Subject() string
		Broker() string
	}

	securityRoleProvider struct {
		service *WebIdentityService
	}
)

func (p *securityRoleProvider) Inject(service *WebIdentityService) {
	p.service = service
}

func (p *securityRoleProvider) All(ctx context.Context, _ *web.Session) []domain.Role {
	request := web.RequestFromContext(ctx)
	if request == nil {
		return nil
	}

	var roles []domain.Role
	var identified bool
	for _, identity := range p.service.IdentifyAll(ctx, request) {
		_ = identity
		identified = true
		// if roler, ok := identity.(hasRoles); ok {
		//roles = append(roles, roler.Roles()...)
		//}
	}

	if identified {
		roles = append(roles, domain.StringRole(domain.PermissionAuthorized))
	}

	return roles
}
