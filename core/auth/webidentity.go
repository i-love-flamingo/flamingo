package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// RequestIdentifierFactory creates RequestIdentifier
	RequestIdentifierFactory func(config config.Map) (RequestIdentifier, error)

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

	// WebLogoutWithRedirect logs user out and redirects to an sso logout endpoint
	WebLogoutWithRedirect interface {
		Logout(ctx context.Context, request *web.Request) *url.URL
	}

	// WebIdentityRefresher refreshs an existing identity, e.g. by invalidating cached session data
	WebIdentityRefresher interface {
		RefreshIdentity(ctx context.Context, request *web.Request) error
	}

	// WebIdentityService calls one or more identifier to get all possible identities of a user
	WebIdentityService struct {
		identityProviders []RequestIdentifier
		reverseRouter     web.ReverseRouter
		eventRouter       flamingo.EventRouter
		responder         *web.Responder
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

	// IdentityTypeChecker for type asserting an Identity
	IdentityTypeChecker func(identity Identity) bool
)

// Inject dependencies
func (s *WebIdentityService) Inject(
	identityProviders []RequestIdentifier,
	reverseRouter web.ReverseRouter,
	eventRouter flamingo.EventRouter,
	responder *web.Responder,
) *WebIdentityService {
	s.identityProviders = identityProviders
	s.reverseRouter = reverseRouter
	s.eventRouter = eventRouter
	s.responder = responder
	return s
}

// Identify the user, if any identity is found
func (s *WebIdentityService) Identify(ctx context.Context, request *web.Request) Identity {
	if s == nil {
		return nil
	}

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

// IdentifyAs returns an identity for a given interface
// identity, err := s.IdentifyAs(ctx, request, OpenIDTypeChecker)
// identity.(oauth.OpenIDIdentity)
func (s *WebIdentityService) IdentifyAs(ctx context.Context, request *web.Request, checkType IdentityTypeChecker) (Identity, error) {
	if s == nil {
		return nil, fmt.Errorf("web identity service is nil")
	}

	for _, provider := range s.identityProviders {
		if identity, _ := provider.Identify(ctx, request); identity != nil {
			if checkType(identity) {
				return identity, nil
			}
		}
	}

	return nil, fmt.Errorf("no identity for type %T found", checkType)
}

func (s *WebIdentityService) storeRedirectURL(request *web.Request) {
	redirecturl, ok := request.Params["redirecturl"]
	if !ok || redirecturl == "" {
		u, err := s.reverseRouter.Absolute(request, "", nil)
		if err == nil {
			redirecturl = u.String()
		}
	}

	absolute, _ := s.reverseRouter.Absolute(request, "", nil)

	refURL, err := url.Parse(redirecturl)
	if err != nil || (redirecturl[0] != '/' && (refURL.Host != request.Request().Host && refURL.Host != absolute.Host)) {
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
	for _, provider := range s.identityProviders {
		if authenticator, ok := provider.(WebAuthenticater); ok {
			if result := authenticator.Authenticate(ctx, request); result != nil {
				s.storeRedirectURL(request)
				return provider.Broker(), result
			}
		}
	}
	return "", nil
}

// AuthenticateFor starts the authentication for a given broker
func (s *WebIdentityService) AuthenticateFor(ctx context.Context, broker string, request *web.Request) web.Result {
	for _, provider := range s.identityProviders {
		if provider.Broker() == broker {
			if authenticator, ok := provider.(WebAuthenticater); ok {
				if result := authenticator.Authenticate(ctx, request); result != nil {
					s.storeRedirectURL(request)
					return result
				}
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

func (s *WebIdentityService) logoutRedirect(request *web.Request, postLogoutRedirect *url.URL) web.Result {
	if len(s.getLogoutRedirects(request)) > 0 {
		if postLogoutRedirect != nil {
			request.Session().Store("core.auth.logoutredirect", postLogoutRedirect)
		} else {
			request.Session().Delete("core.auth.logoutredirect")
		}
		return s.responder.RouteRedirect("core.auth.logoutCallback", nil)
	}
	if postLogoutRedirect != nil {
		return s.responder.URLRedirect(postLogoutRedirect)
	}
	return s.responder.RouteRedirect("", nil)
}

func (s *WebIdentityService) logout(ctx context.Context, request *web.Request, postLogoutRedirect *url.URL, broker string, all bool) web.Result {
	s.storeLogoutRedirects(request, redirectURLlist{})

	for _, provider := range s.identityProviders {
		if provider.Broker() == broker || all {
			if authenticator, ok := provider.(WebLogouter); ok {
				authenticator.Logout(ctx, request)
				s.eventRouter.Dispatch(ctx, &WebLogoutEvent{Request: request, Broker: provider.Broker()})
			} else if authenticator, ok := provider.(WebLogoutWithRedirect); ok {
				s.addLogoutRedirect(request, authenticator.Logout(ctx, request))
				s.eventRouter.Dispatch(ctx, &WebLogoutEvent{Request: request, Broker: provider.Broker()})
			}
		}
	}

	return s.logoutRedirect(request, postLogoutRedirect)
}

// Logout logs all user out
func (s *WebIdentityService) Logout(ctx context.Context, request *web.Request, postLogoutRedirect *url.URL) web.Result {
	return s.logout(ctx, request, postLogoutRedirect, "", true)
}

// LogoutFor logs a specific broker out
func (s *WebIdentityService) LogoutFor(ctx context.Context, broker string, request *web.Request, postLogoutRedirect *url.URL) web.Result {
	return s.logout(ctx, request, postLogoutRedirect, broker, false)
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

type redirectURLlist []*url.URL

func init() {
	gob.Register(redirectURLlist{})
	gob.Register(new(url.URL))
}

func (s *WebIdentityService) addLogoutRedirect(request *web.Request, u *url.URL) {
	if u == nil {
		return
	}
	redirects := s.getLogoutRedirects(request)
	redirects = append(redirects, u)
	s.storeLogoutRedirects(request, redirects)
}

func (s *WebIdentityService) getLogoutRedirects(request *web.Request) []*url.URL {
	rawredirects, _ := request.Session().Load("core.auth.logoutredirects")
	redirects, _ := rawredirects.(redirectURLlist)
	return redirects
}

func (s *WebIdentityService) storeLogoutRedirects(request *web.Request, redirects []*url.URL) {
	request.Session().Store("core.auth.logoutredirects", redirectURLlist(redirects))
}
