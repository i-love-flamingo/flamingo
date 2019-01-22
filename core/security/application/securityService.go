package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/security/application/voter"
	"flamingo.me/flamingo/v3/core/security/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

const (
	VoterStrategyAffirmative = "affirmative"
	VoterStrategyConsensus   = "consensus"
	VoterStrategyUnanimous   = "unanimous"
)

type (
	SecurityService interface {
		IsLoggedIn(context.Context, *web.Session) bool
		IsLoggedOut(context.Context, *web.Session) bool
		IsGranted(context.Context, *web.Session, string, interface{}) bool
	}

	SecurityServiceImpl struct {
		voters            []voter.SecurityVoter
		voterStrategy     string
		allowIfAllAbstain bool
	}
)

var (
	_ SecurityService = &SecurityServiceImpl{}
)

func (s *SecurityServiceImpl) Inject(v []voter.SecurityVoter, cfg *struct {
	VoterStrategy     string `inject:"config:security.roles.voters.strategy"`
	AllowIfAllAbstain bool   `inject:"config:security.roles.voters.allowIfAllAbstain"`
}) {
	s.voters = v
	s.voterStrategy = cfg.VoterStrategy
	s.allowIfAllAbstain = cfg.AllowIfAllAbstain
}

func (s *SecurityServiceImpl) IsLoggedIn(ctx context.Context, session *web.Session) bool {
	return s.IsGranted(ctx, session, domain.RoleUser.Permission(), nil)
}

func (s *SecurityServiceImpl) IsLoggedOut(ctx context.Context, session *web.Session) bool {
	return !s.IsGranted(ctx, session, domain.RoleUser.Permission(), nil)
}

func (s *SecurityServiceImpl) IsGranted(ctx context.Context, session *web.Session, permission string, object interface{}) bool {
	var results []int
	for index := range s.voters {
		results = append(results, s.voters[index].Vote(ctx, session, permission, object))
	}

	return s.decide(results)
}

func (s *SecurityServiceImpl) decide(results []int) bool {
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
