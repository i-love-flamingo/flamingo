package voter

import (
	"flamingo.me/flamingo/v3/core/security/domain"
)

// PermissionVoter votes on specific permission
type PermissionVoter struct {}

// Vote for permission
func (v *PermissionVoter) Vote(allAssignedPermissions []string, desiredPermission string, forObject interface{}) AccessDecision {
	if desiredPermission == domain.PermissionAuthorized {
		return AccessAbstained
	}

	permissionSet, ok := forObject.(domain.PermissionSet)
	if ok && !v.hasPermission(permissionSet.Permissions(), desiredPermission) {
		return AccessDenied
	}

	if !v.hasPermission(allAssignedPermissions, desiredPermission) {
		return AccessDenied
	}
	return AccessGranted
}

func (v *PermissionVoter) hasPermission(allAssignedPermissions []string, desiredPermission string) bool {
	for _, permission := range allAssignedPermissions {
		if permission == desiredPermission {
			return true
		}
	}

	return false
}
