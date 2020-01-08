package auth

import (
	"context"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// IdentifierFactory creates RequestIdentifier
	IdentifierFactory func(config config.Map) (RequestIdentifier, error)

	// RequestIdentifier identifies an request and returns a matching identity
	RequestIdentifier interface {
		Broker() string
		Identify(ctx context.Context, request *web.Request) (Identity, error)
	}

	// WebAuthenticater allows to request an authentication
	WebAuthenticater interface {
		Authenticate(ctx context.Context, request *web.Request) web.Result
	}

	// WebCallbacker is called for callbacks to that identity broker
	WebCallbacker interface {
		Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result
	}

	// WebLogouter logs user out
	WebLogouter interface {
		Logout(ctx context.Context, request *web.Request)
	}

	// WebIdentityService calls one or more identifier to get all possible identities of a user
	WebIdentityService struct {
		identityProviders []RequestIdentifier
		reverseRouter     web.ReverseRouter
		eventRouter       flamingo.EventRouter
	}

	// WebLoginEvent for the current request
	WebLoginEvent struct {
		Request  *web.Request
		Broker   string
		Identity Identity
	}

	// WebLogoutEvent for the current request
	WebLogoutEvent struct {
		Request *web.Request
		Broker  string
	}
)

// Inject dependencies
func (s *WebIdentityService) Inject(
	identityProviders []RequestIdentifier,
	reverseRouter web.ReverseRouter,
	eventRouter flamingo.EventRouter,
) {
	s.identityProviders = identityProviders
	s.reverseRouter = reverseRouter
	s.eventRouter = eventRouter
}

// Identify the user, if any identity is found
func (s *WebIdentityService) Identify(ctx context.Context, request *web.Request) Identity {
	for _, provider := range s.identityProviders {
		if identity, _ := provider.Identify(ctx, request); identity != nil {
			return identity
		}
	}

	return nil
}

// IdentifyFor the user with a given broker
func (s *WebIdentityService) IdentifyFor(ctx context.Context, broker string, request *web.Request) (Identity, error) {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			return provider.Identify(ctx, request)
		}
	}

	return nil, fmt.Errorf("no broker with code %q found", broker)
}

// IdentifyAll collects all possible user identites, in case multiple are available
func (s *WebIdentityService) IdentifyAll(ctx context.Context, request *web.Request) []Identity {
	var identities []Identity

	for _, provider := range s.identityProviders {
		if identity, _ := provider.Identify(ctx, request); identity != nil {
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

// AuthenticateFor starts the authentication for a given broker
func (s *WebIdentityService) AuthenticateFor(ctx context.Context, broker string, request *web.Request) web.Result {
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

// Logout logs all user out
func (s *WebIdentityService) Logout(ctx context.Context, request *web.Request) {
	for _, provider := range s.identityProviders {
		if authenticator, ok := provider.(WebLogouter); ok {
			authenticator.Logout(ctx, request)
			s.eventRouter.Dispatch(ctx, &WebLogoutEvent{Request: request, Broker: provider.Broker()})
		}
	}
}

// LogoutFor logs a specific broker out
func (s *WebIdentityService) LogoutFor(ctx context.Context, broker string, request *web.Request) {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			if authenticator, ok := provider.(WebLogouter); ok {
				authenticator.Logout(ctx, request)
				s.eventRouter.Dispatch(ctx, &WebLogoutEvent{Request: request, Broker: provider.Broker()})
			}
		}
	}
}

// RequestIdentifier returns the given request identifier
func (s *WebIdentityService) RequestIdentifier(broker string) RequestIdentifier {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			return provider
		}
	}
	return nil
}
