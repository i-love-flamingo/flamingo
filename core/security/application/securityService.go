package application

import (
	"context"
	"flamingo.me/flamingo/v3/core/security/application/role"
	"fmt"

	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

const (
	// VoterStrategyAffirmative allows access if there is a positive vote
	VoterStrategyAffirmative = "affirmative"
	// VoterStrategyConsensus allows access if there are more positive votes
	VoterStrategyConsensus = "consensus"
	// VoterStrategyUnanimous allows access if there are no negative votes
	VoterStrategyUnanimous = "unanimous"
)

type (
	// SecurityService decides if a user is logged in/out, or granted a certain permission
	// todo name arguments
	SecurityService interface {
		IsLoggedIn(context.Context, *web.Session) bool
		IsLoggedOut(context.Context, *web.Session) bool
		IsGranted(context.Context, *web.Session, string, interface{}) bool
	}

	// SecurityServiceImpl default implementation of the SecurityService
	SecurityServiceImpl struct {
		voters            []voter.SecurityVoter
		roleService       role.Service
		voterStrategy     string
		allowIfAllAbstain bool
	}
)

var _ SecurityService = &SecurityServiceImpl{}

// Inject dependencies
func (s *SecurityServiceImpl) Inject(v []voter.SecurityVoter, r role.Service, cfg *struct {
	VoterStrategy     string `inject:"config:security.roles.voters.strategy"`
	AllowIfAllAbstain bool   `inject:"config:security.roles.voters.allowIfAllAbstain"`
}) {
	s.voters = v
	s.roleService = r
	s.voterStrategy = cfg.VoterStrategy
	s.allowIfAllAbstain = cfg.AllowIfAllAbstain
}

// IsLoggedIn checks if the user is granted login permission
func (s *SecurityServiceImpl) IsLoggedIn(ctx context.Context, session *web.Session) bool {
	return s.IsGranted(ctx, session, domain.PermissionAuthorized, nil)
}

// IsLoggedOut checks if the user is not granted login permission
func (s *SecurityServiceImpl) IsLoggedOut(ctx context.Context, session *web.Session) bool {
	return !s.IsGranted(ctx, session, domain.PermissionAuthorized, nil)
}

// IsGranted checks for a specific permission of the user
func (s *SecurityServiceImpl) IsGranted(ctx context.Context, session *web.Session, desiredPermission string, object interface{}) bool {
	allPermissions := s.roleService.AllPermissions(ctx, session)

	var results []voter.AccessDecision
	for index := range s.voters {
		results = append(results, s.voters[index].Vote(allPermissions, desiredPermission, object))
	}

	return s.decide(results)
}

func (s *SecurityServiceImpl) decide(results []voter.AccessDecision) bool {
	granted := 0
	denied := 0

	for _, result := range results {
		switch result {
		case voter.AccessGranted:
			granted++
		case voter.AccessDenied:
			denied++
		}
	}

	switch s.voterStrategy {
	case VoterStrategyAffirmative:
		return s.decideAffirmative(granted, denied)
	case VoterStrategyConsensus:
		return s.decideConsensus(granted, denied)
	case VoterStrategyUnanimous:
		return s.decideUnanimous(granted, denied)
	}
	panic(fmt.Sprintf("unrecognized voter strategy: %s", s.voterStrategy))
}

func (s *SecurityServiceImpl) decideAffirmative(granted int, denied int) bool {
	if granted > 0 {
		return true
	} else if denied > 0 {
		return false
	}
	return s.allowIfAllAbstain
}

func (s *SecurityServiceImpl) decideConsensus(granted int, denied int) bool {
	if granted > denied {
		return true
	} else if denied > granted {
		return false
	}
	return s.allowIfAllAbstain
}

func (s *SecurityServiceImpl) decideUnanimous(granted int, denied int) bool {
	if denied > 0 {
		return false
	} else if granted > 0 {
		return true
	}
	return s.allowIfAllAbstain
}
