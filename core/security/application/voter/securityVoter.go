package voter

type (
	// AccessDecision defines access decision type which represents voter result
	AccessDecision int

	// SecurityVoter defines a common interface for voters who vote on security decisions
	SecurityVoter interface {
		Vote(allAssignedPermissions []string, desiredPermission string, forObject interface{}) AccessDecision
	}
)

const (
	// AccessAbstained defines access decision in case voter is not responsible for permission
	AccessAbstained AccessDecision = iota
	// AccessGranted defines access decision in case when voter grants an access
	AccessGranted AccessDecision = iota
	// AccessDenied defines access decision in case when voter denies an access
	AccessDenied AccessDecision = iota
)
