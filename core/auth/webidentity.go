package auth

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	IdentifierFactory func(config config.Map) Identifier

	Storage interface {
		Load(key string) (data interface{}, ok bool)
		Store(key string, value interface{})
	}

	sessionStorage struct {
		session *web.Session
	}

	Identifier interface {
		Broker() string
	}

	StorageIdentifier interface {
		Identifier
		Identify(ctx context.Context, storage Storage) Identity
	}

	RequestIdentifier interface {
		Identifier
		Identify(ctx context.Context, request *web.Request) Identity
	}

	WebAuthenticater interface {
		Authenticate(ctx context.Context, request *web.Request) web.Result
	}

	WebCallbacker interface {
		Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result
	}

	// WebIdentityService calls one or more identifier to get all possible identities of a user
	WebIdentityService struct {
		identityProviders []Identifier
		reverseRouter     web.ReverseRouter
	}
)

func (s sessionStorage) Load(key string) (data interface{}, ok bool) {
	return s.session.Load(key)
}

func (s sessionStorage) Store(key string, value interface{}) {
	s.session.Store(key, value)
}

func (s *WebIdentityService) Inject(identityProviders []Identifier, reverseRouter web.ReverseRouter) {
	s.identityProviders = identityProviders
	s.reverseRouter = reverseRouter
}

func identify(identifier Identifier, ctx context.Context, request *web.Request) Identity {
	switch identifier := identifier.(type) {
	case StorageIdentifier:
		return identifier.Identify(ctx, sessionStorage{session: request.Session()})
	case RequestIdentifier:
		return identifier.Identify(ctx, request)
	}
	return nil
}

// Identify the user, if any identity is found
func (s *WebIdentityService) Identify(ctx context.Context, request *web.Request) Identity {
	for _, provider := range s.identityProviders {
		if identity := identify(provider, ctx, request); identity != nil {
			return identity
		}
	}

	return nil
}

// Identify the user, if any identity is found
func (s *WebIdentityService) IdentifyFor(broker string, ctx context.Context, request *web.Request) Identity {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			return identify(provider, ctx, request)
		}
	}

	return nil
}

// IdentifyAll collects all possible user identites, in case multiple are available
func (s *WebIdentityService) IdentifyAll(ctx context.Context, request *web.Request) []Identity {
	var identities []Identity

	for _, provider := range s.identityProviders {
		if identity := identify(provider, ctx, request); identity != nil {
			identities = append(identities, identity)
		}
	}

	return identities
}

// Authenticate finds the first available (enforced) authentication result
func (s *WebIdentityService) Authenticate(ctx context.Context, request *web.Request) (string, web.Result) {
	for _, provider := range s.identityProviders {
		if authenticator, ok := provider.(WebAuthenticater); ok {
			if result := authenticator.Authenticate(ctx, request); result != nil {
				return provider.Broker(), result
			}
		}
	}
	return "", nil
}

// Authenticate finds the first available (enforced) authentication result
func (s *WebIdentityService) AuthenticateFor(broker string, ctx context.Context, request *web.Request) web.Result {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			if authenticator, ok := provider.(WebAuthenticater); ok {
				return authenticator.Authenticate(ctx, request)
			}
			return nil
		}
	}
	return nil
}
func (s *WebIdentityService) callback(ctx context.Context, request *web.Request) web.Result {
	broker := request.Params["broker"]
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			return provider.(WebCallbacker).Callback(ctx, request, func(request *web.Request) *url.URL {
				u, err := s.reverseRouter.Absolute(request, "core.auth.debug", nil)
				if err != nil {
					panic(err)
				}
				return u
			})
		}
	}
	return nil
}
