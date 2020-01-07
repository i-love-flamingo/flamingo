package auth

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	IdentifierFactory func(config config.Map) RequestIdentifier

	Storage interface {
		Load(key string) (data interface{}, ok bool)
		Store(key string, value interface{})
		Delete(key string)
	}

	sessionStorage struct {
		session *web.Session
	}

	RequestIdentifier interface {
		Broker() string
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
		identityProviders []RequestIdentifier
		reverseRouter     web.ReverseRouter
	}
)

func (s sessionStorage) Load(key string) (data interface{}, ok bool) {
	return s.session.Load(key)
}

func (s sessionStorage) Store(key string, value interface{}) {
	s.session.Store(key, value)
}

func (s *WebIdentityService) Inject(identityProviders []RequestIdentifier, reverseRouter web.ReverseRouter) {
	s.identityProviders = identityProviders
	s.reverseRouter = reverseRouter
}

func identify(identifier RequestIdentifier, ctx context.Context, request *web.Request) Identity {
	return identifier.Identify(ctx, request)
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

func (s *WebIdentityService) storeRedirectURL(request *web.Request) {
	redirecturl, ok := request.Params["redirecturl"]
	if !ok || redirecturl == "" {
		redirecturl = request.Request().Referer()
	}

	absolute, _ := s.reverseRouter.Absolute(request, "", nil)

	refURL, err := url.Parse(redirecturl)
	if err != nil || (refURL.Host != request.Request().Host && refURL.Host != absolute.Host) {
		redirecturl = absolute.String()
	}

	request.Session().Store("core.auth.redirect", redirecturl)
}

func (s *WebIdentityService) getRedirectURL(request *web.Request) *url.URL {
	redirectURL, ok := request.Session().Load("core.auth.redirect")
	if !ok {
		rurl, _ := s.reverseRouter.Absolute(request, "", nil)
		return rurl
	}
	rurl, err := url.Parse(redirectURL.(string))
	if err != nil {
		rurl, _ = s.reverseRouter.Absolute(request, "", nil)
	}
	return rurl
}

// Authenticate finds the first available (enforced) authentication result
func (s *WebIdentityService) Authenticate(ctx context.Context, request *web.Request) (string, web.Result) {
	s.storeRedirectURL(request)
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
	s.storeRedirectURL(request)
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
			return provider.(WebCallbacker).Callback(ctx, request, s.getRedirectURL)
		}
	}
	return nil
}
