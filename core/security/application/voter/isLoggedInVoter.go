package voter

import (
	"flamingo.me/flamingo/v3/core/security/domain"
)

// IsLoggedInVoter votes for users who have authenticated
type IsLoggedInVoter struct{}

// Vote for the authentication request
func (v *IsLoggedInVoter) Vote(allAssignedPermissions []string, desiredPermission string, _ interface{}) AccessDecision {
	if desiredPermission != domain.PermissionAuthorized {
		return AccessAbstained
	}

	for _, permission := range allAssignedPermissions {
		if permission == domain.PermissionAuthorized {
			return AccessGranted
		}
	}

	return AccessDenied
}
