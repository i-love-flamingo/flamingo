// Package oidcfake is a fake OpenID Connect provider for the contract tests.
// It runs mockoidc on a free local port and adds an end session endpoint, which mockoidc does not have.
package oidcfake

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/oauth2-proxy/mockoidc"
)

// Server is a running fake provider. It stops at the end of the test.
type Server struct {
	Issuer       string // the broker "endpoint"
	ClientID     string
	ClientSecret string
	EndSession   string // URL of the end_session endpoint

	oidc            *mockoidc.MockOIDC
	endSessionCalls atomic.Int32

	mu           sync.Mutex
	grantTypes   []string
	redirectURIs map[string]string // code to the redirect_uri of its authorize request
}

// Start starts the provider on 127.0.0.1 and stops it with the cleanup of the test.
// The end_session endpoint runs on a second server, because mockoidc runs its middleware only on its own paths.
func Start(tb testing.TB) *Server {
	tb.Helper()

	server := &Server{redirectURIs: map[string]string{}}

	endSession := httptest.NewServer(http.HandlerFunc(server.endSession))
	tb.Cleanup(endSession.Close)

	server.EndSession = endSession.URL + "/end_session"

	oidc, err := mockoidc.NewServer(nil)
	if err != nil {
		tb.Fatalf("oidcfake: new server: %v", err)
	}

	for _, middleware := range []func(http.Handler) http.Handler{server.discovery, server.recordGrantType, server.checkRedirectURI} {
		if err := oidc.AddMiddleware(middleware); err != nil {
			tb.Fatalf("oidcfake: add middleware: %v", err)
		}
	}

	listener, err := new(net.ListenConfig).Listen(tb.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("oidcfake: listen: %v", err)
	}

	if err := oidc.Start(listener, nil); err != nil {
		tb.Fatalf("oidcfake: start: %v", err)
	}

	tb.Cleanup(func() {
		if err := oidc.Shutdown(); err != nil {
			tb.Errorf("oidcfake: shutdown: %v", err)
		}
	})

	server.oidc = oidc
	server.Issuer = oidc.Issuer()
	server.ClientID = oidc.ClientID
	server.ClientSecret = oidc.ClientSecret

	return server
}

// EndSessionCalls gives the number of requests that the end_session endpoint got.
func (s *Server) EndSessionCalls() int {
	return int(s.endSessionCalls.Load())
}

// GrantTypes gives the grant_type of each request to the token endpoint that got a token, in sequence.
func (s *Server) GrantTypes() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]string{}, s.grantTypes...)
}

// BrokerYAML gives a config file part with one oidc entry in core.auth.web.broker.
// The entry maps claims of both tokens to keys with other names, so a swapped or unmapped claim shows.
// mockoidc does not know the scope "offline_access", so enableOfflineToken is false.
func (s *Server) BrokerYAML(name string) string {
	return fmt.Sprintf(`core.auth.web.broker:
  - broker: %[1]s
    typ: oidc
    endpoint: %[2]q
    clientID: %[3]q
    clientSecret: %[4]q
    scopes: [profile, email]
    enableOfflineToken: false
    enableEndSessionEndpoint: true
    claims:
      idToken:
        userID: sub
        mail: email
      accessToken:
        tokenSubject: sub
`, name, s.Issuer, s.ClientID, s.ClientSecret)
}

// QueueUser makes the next login a user with only a subject and a verified email.
// The default user of mockoidc has a name, a phone number and an address.
func (s *Server) QueueUser(subject, email string) {
	s.oidc.QueueUser(&mockoidc.MockUser{Subject: subject, Email: email, EmailVerified: true})
}

// discovery adds the end_session_endpoint to the discovery document.
func (s *Server) discovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		if r.URL.Path != mockoidc.DiscoveryEndpoint {
			next.ServeHTTP(writer, r)

			return
		}

		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		var document map[string]any
		if err := json.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)

			return
		}

		document["end_session_endpoint"] = s.EndSession

		body, err := json.Marshal(document)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)

			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(recorder.Code)
		_, _ = writer.Write(body)
	})
}

// grantWriter remembers the grant_type of a token request, when the response has the status 200.
// It does so before it sends the status, so a client that has the response sees the record.
type grantWriter struct {
	http.ResponseWriter

	server    *Server
	grantType string
}

func (w *grantWriter) WriteHeader(status int) {
	if status == http.StatusOK {
		w.server.mu.Lock()
		w.server.grantTypes = append(w.server.grantTypes, w.grantType)
		w.server.mu.Unlock()
	}

	w.ResponseWriter.WriteHeader(status)
}

// recordGrantType remembers the grant_type of each token request that the provider answers with 200.
// The oauth2 package first sends the client credentials in the Authorization header. mockoidc reads them
// only from the form and answers 400, then the package sends them again in the form.
func (s *Server) recordGrantType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == mockoidc.TokenEndpoint && r.ParseForm() == nil {
			w = &grantWriter{ResponseWriter: w, server: s, grantType: r.PostForm.Get("grant_type")}
		}

		next.ServeHTTP(w, r)
	})
}

// codeWriter remembers the code of the redirect of an authorize response, with the redirect_uri of the request.
type codeWriter struct {
	http.ResponseWriter

	server      *Server
	redirectURI string
}

func (w *codeWriter) WriteHeader(status int) {
	if target, err := url.Parse(w.Header().Get("Location")); err == nil && target.Query().Get("code") != "" {
		w.server.mu.Lock()
		w.server.redirectURIs[target.Query().Get("code")] = w.redirectURI
		w.server.mu.Unlock()
	}

	w.ResponseWriter.WriteHeader(status)
}

// checkRedirectURI makes the provider compare the redirect_uri of a code exchange with the one of the authorize
// request, as RFC 6749 section 4.1.3 requires. mockoidc does not compare them.
// The check answers 400 invalid_grant, as a provider does, when they differ or when the code is unknown.
func (s *Server) checkRedirectURI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		if req.ParseForm() != nil {
			next.ServeHTTP(writer, req)

			return
		}

		switch req.URL.Path {
		case mockoidc.AuthorizationEndpoint:
			writer = &codeWriter{ResponseWriter: writer, server: s, redirectURI: req.Form.Get("redirect_uri")}
		case mockoidc.TokenEndpoint:
			s.mu.Lock()
			want, known := s.redirectURIs[req.PostForm.Get("code")]
			s.mu.Unlock()

			if req.PostForm.Get("grant_type") == "authorization_code" && (!known || req.PostForm.Get("redirect_uri") != want) {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(`{"error":"invalid_grant","error_description":"unknown code, or redirect_uri does not match the authorize request"}`))

				return
			}
		}

		next.ServeHTTP(writer, req)
	})
}

// endSession counts the call and redirects to the post_logout_redirect_uri, as a provider does.
func (s *Server) endSession(w http.ResponseWriter, r *http.Request) {
	s.endSessionCalls.Add(1)

	target := r.URL.Query().Get("post_logout_redirect_uri")
	if target == "" || !strings.HasPrefix(target, "http") {
		w.WriteHeader(http.StatusOK)

		return
	}

	http.Redirect(w, r, target, http.StatusFound) //nolint:gosec // the fake sends the user where the app tells it, as a provider does
}
