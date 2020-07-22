package mock

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"
)

// SingleIdentityIdentifyMethod returns an IdentifyMethod which will always return the specified identity
func SingleIdentityIdentifyMethod(identity auth.Identity) IdentifyMethod {
	return func(identifier *Identifier, ctx context.Context, request *web.Request) (auth.Identity, error) {
		return identity, nil
	}
}

// SingleIdentityIdentifier returns a mocked broker with a given identity
func SingleIdentityIdentifier(broker string, identity auth.Identity) *Identifier {
	return &Identifier{
		broker:         broker,
		identifyMethod: SingleIdentityIdentifyMethod(identity),
	}
}
