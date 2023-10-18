package oauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type mockRouter struct {
	broker string
}

var _ web.ReverseRouter = (*mockRouter)(nil)

// Relative mock action
func (m *mockRouter) Relative(_ string, _ map[string]string) (*url.URL, error) {
	panic("not implemented")
}

// Absolute mock action
func (m *mockRouter) Absolute(_ *web.Request, _ string, _ map[string]string) (*url.URL, error) {
	return url.Parse(strings.ReplaceAll("/core/auth/login/:broker", ":broker", m.broker))
}

type mockCallbackErrorHandler struct {
	Called                   bool
	SuppliedError            string
	SuppliedErrorDescription string
}

func (m *mockCallbackErrorHandler) Handle(_ context.Context, _ string, _ *web.Request, _ func(request *web.Request) *url.URL, err string, errDesc string) web.Result {
	m.Called = true
	m.SuppliedError = err
	m.SuppliedErrorDescription = errDesc

	target, _ := url.Parse("https://example.com/callback-error-handler")

	return &web.URLRedirectResponse{URL: target}
}

var _ CallbackErrorHandler = &mockCallbackErrorHandler{}

func TestParallelStateRaceConditions(t *testing.T) {
	t.Run("test states", func(t *testing.T) {
		t.Parallel()

		identifier := &openIDIdentifier{
			authCodeOptionerProvider: func() []AuthCodeOptioner { return nil },
			oauth2Config:             &oauth2.Config{},
			reverseRouter:            new(mockRouter),
			responder:                &web.Responder{},
		}

		session := web.EmptySession()

		resp := identifier.Authenticate(context.Background(), web.CreateRequest(nil, session))
		state1 := resp.(*web.URLRedirectResponse).URL.Query().Get("state")
		resp = identifier.Authenticate(context.Background(), web.CreateRequest(nil, session))
		state2 := resp.(*web.URLRedirectResponse).URL.Query().Get("state")

		request, err := http.NewRequest(http.MethodGet, "http://example.com/callback", nil)
		assert.NoError(t, err)

		request.URL.RawQuery = url.Values{"state": []string{"invalid-state"}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp := resp.(*web.ServerErrorResponse)
		assert.EqualError(t, errResp.Error, "state mismatch")

		request.URL.RawQuery = url.Values{"state": []string{state2}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp = resp.(*web.ServerErrorResponse)
		assert.EqualError(t, errResp.Error, "query value not found")

		request.URL.RawQuery = url.Values{"state": []string{state1}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp = resp.(*web.ServerErrorResponse)
		assert.EqualError(t, errResp.Error, "query value not found")

		request.URL.RawQuery = url.Values{"state": []string{state1}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp = resp.(*web.ServerErrorResponse)
		assert.EqualError(t, errResp.Error, "state mismatch")
	})

	t.Run("test default time shift", func(t *testing.T) {
		t.Parallel()

		identifier := &openIDIdentifier{
			authCodeOptionerProvider: func() []AuthCodeOptioner { return nil },
			oauth2Config:             &oauth2.Config{},
			reverseRouter:            new(mockRouter),
			responder:                &web.Responder{},
		}

		request, err := http.NewRequest(http.MethodGet, "http://example.com/callback", nil)
		assert.NoError(t, err)

		session := web.EmptySession()

		resp := identifier.Authenticate(context.Background(), web.CreateRequest(nil, session))
		state1 := resp.(*web.URLRedirectResponse).URL.Query().Get("state")

		now = func() time.Time {
			return time.Now().Add(35 * time.Minute)
		}

		request.URL.RawQuery = url.Values{"state": []string{state1}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp := resp.(*web.ServerErrorResponse)
		assert.EqualError(t, errResp.Error, "state mismatch")

		now = time.Now
	})

	t.Run("test custom time shift", func(t *testing.T) {
		t.Parallel()

		identifier := &openIDIdentifier{
			authCodeOptionerProvider: func() []AuthCodeOptioner { return nil },
			oauth2Config:             &oauth2.Config{},
			reverseRouter:            new(mockRouter),
			responder:                &web.Responder{},
		}

		request, err := http.NewRequest(http.MethodGet, "http://example.com/callback", nil)
		assert.NoError(t, err)

		session := web.EmptySession()

		resp := identifier.Authenticate(context.Background(), web.CreateRequest(nil, session))
		state1 := resp.(*web.URLRedirectResponse).URL.Query().Get("state")
		oneHour := time.Hour
		identifier.stateTimeout = &oneHour

		now = func() time.Time {
			return time.Now().Add(35 * time.Minute)
		}

		request.URL.RawQuery = url.Values{"state": []string{state1}}.Encode()
		resp = identifier.Callback(context.Background(), web.CreateRequest(request, session), nil)
		errResp := resp.(*web.ServerErrorResponse)
		assert.NotContains(t, errResp.Error.Error(), "state mismatch")

		now = time.Now
	})
}

type testOidcProvider struct {
	url string
}

func (p *testOidcProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch strings.Trim(r.URL.Path, "/") {
	case ".well-known/openid-configuration":
		_, _ = fmt.Fprintf(rw, `{
			"issuer": "%s",
			"token_endpoint": "%s/token",
			"jwks_uri": "%s/certs"
		}`, p.url, p.url, p.url)
	case "certs":
		_, _ = fmt.Fprint(rw, `{"keys":[{"kid":"3vfKZm_SqyuYCsD7isNlzs1EORs5guIF0XnisUqvjFQ","kty":"RSA","alg":"RS256","use":"sig","n":"o1ZomEdwneplkdgUJMUqHjaZRt3qCP7wQJq0XwK2J95LXYIwPaXy9b2IQruOfhhoy34hWnzsJkQpnugFj069qxZ3ni7Cb1Wau2x6xuhBoiko1iaB_IddG8tSi0FyMUhMtSpf8eBPQB9i5UX8Uymj36B4z2HfbDLLU8ld2ve4PNCFnRDwWCVzE_LwER0rIDeIpODptKVL8bEPsgLgGCB7WGIFIpVwmaS8UwWrCPtlMBJVhHOi1GwuAoPVhOTQhqyCrxh9c3PcZjO0o-yW2BBdSWSHs69zboa-DPGg4jUo7SHkGL4tC--Nvg46ZojNtmjlMJpquK3XJA6SC8l-W776Tw","e":"AQAB","x5c":["MIICmzCCAYMCBgF/MQ1x7DANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZtYXN0ZXIwHhcNMjIwMjI1MTMyMjE5WhcNMzIwMjI1MTMyMzU5WjARMQ8wDQYDVQQDDAZtYXN0ZXIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCjVmiYR3Cd6mWR2BQkxSoeNplG3eoI/vBAmrRfArYn3ktdgjA9pfL1vYhCu45+GGjLfiFafOwmRCme6AWPTr2rFneeLsJvVZq7bHrG6EGiKSjWJoH8h10by1KLQXIxSEy1Kl/x4E9AH2LlRfxTKaPfoHjPYd9sMstTyV3a97g80IWdEPBYJXMT8vARHSsgN4ik4Om0pUvxsQ+yAuAYIHtYYgUilXCZpLxTBasI+2UwElWEc6LUbC4Cg9WE5NCGrIKvGH1zc9xmM7Sj7JbYEF1JZIezr3Nuhr4M8aDiNSjtIeQYvi0L742+DjpmiM22aOUwmmq4rdckDpILyX5bvvpPAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAJIa6nWcYN6AHpYvpQeA62kjXzixXTb3sS5TCx1MVIrA1HK9oYkRqp8L6js0HZ9r4Bi6m7phuh9nssHQQo1HUWnkMvBXi7Su8OstUpMV3cef7E2eOiXl/XXoKOYzn00wuviajofGL6JopV9RpIGsZoU8mmjmpBpRcby/V9ILsCZeU6Q/mQw7xTG7eRZZOPtqgSdvOXxDWFrpycFk9ZaEBI8bVchfQ0B19VsmD/2P2ujGgRgxlZIx6R+gNDn6PiF+acbEaSqnl4WyrajgDQp0ZSsWhrSQ+AgrH3lGOe2KBMGrpaI93vEzSM/wMBgSwrGovnYaiiCD58uuAMnebrSPLnQ="],"x5t":"3pqpFR1KsffXGgRSwaEpW9HB0uQ","x5t#S256":"UXfYw2fnfmPyGfA4___BJynPIB9ZtHSvcE5A5jiKC14"}]}`)
	case "token":
		idtoken := `eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzdmZLWm1fU3F5dVlDc0Q3aXNObHpzMUVPUnM1Z3VJRjBYbmlzVXF2akZRIn0.eyJleHAiOjE2NDYzMTQ0NTQsImlhdCI6MTY0NjMxNDM5NCwiYXV0aF90aW1lIjoxNjQ2MzE0MTM0LCJqdGkiOiJhZGM0MDAzYy0wNTEzLTQ4ZjYtYjdiOS0wYTNjOGY2YmVlYWIiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoic2VjdXJpdHktYWRtaW4tY29uc29sZSIsInN1YiI6IjA0N2UwMDdhLTFhNTgtNDQ3Ni05ODc3LTFhMWIwN2U2OTlhMSIsInR5cCI6IklEIiwiYXpwIjoic2VjdXJpdHktYWRtaW4tY29uc29sZSIsIm5vbmNlIjoiNzI0ZjdhMTgtMmVkZi00MTBmLThhZmUtYWExYWNiNWI0YWI3Iiwic2Vzc2lvbl9zdGF0ZSI6IjQxN2MzZmExLTY2ZmUtNGQ2MC05ZDgyLTE0NTkxNmExMWNmNCIsImF0X2hhc2giOiJmUG4yWDZpakIxSFdVWWRNNGE3bjJnIiwiYWNyIjoiMSIsInNpZCI6IjQxN2MzZmExLTY2ZmUtNGQ2MC05ZDgyLTE0NTkxNmExMWNmNCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwicHJlZmVycmVkX3VzZXJuYW1lIjoibG9jYWxhZG1pbiJ9.AcnJ_RuTBWWBK_MvGTvvD2LBhLlIoOA_F7_TcAm9Cit2tRoUPUWAiclHE3pn2UsJ1YnEPJCraDFC1Ef0KBHTXZN9DmRrGCJxLrArQ89PuYrVSK6f7-C3TWwHKDTpHwTfONFBwuFmGxFbIGjzd8f3xX4jT0ISGdXCyP8Lv5yI0r7oVN_saqlIiC2X50etPicpFV8JrcVFhmBIDhGgATnl5FOWs5_xlMN1fWGWDyfZglkrWSyJPg73dt4JvEdnbZeI1vDLf_AXjjIVigzAAJfV2ZurcozGy1iaMO-UfghUsDIn8UUhQUtJk8EUSrCuYgIK8L5JImODp_IBLrcHAxQCow`
		accesstoken := fmt.Sprintf("header.%s", base64.RawURLEncoding.EncodeToString([]byte(`{"at_claim_1": "at-claim-1-value"}`)))
		rw.Header().Set("Content-type", "application/json")
		_, _ = fmt.Fprintf(rw, `{
			"access_token": "%s",
			"id_token": "%s",
			"legacy-token-response-claim": "legacy-token-response-claim-value"
		}`, accesstoken, idtoken)
	default:
		panic("unable to handle: " + r.URL.Path)
	}
}

type testAuthCodeOptioner struct{}

func (*testAuthCodeOptioner) Options(_ context.Context, _ string, req *web.Request) []oauth2.AuthCodeOption {
	return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("redirect_uri", "foobar123test")}
}

func TestOidcCallback(t *testing.T) {
	t.Run("Test Callback", func(t *testing.T) {
		t.Parallel()
		provider := &testOidcProvider{}
		testserver := httptest.NewServer(provider)
		defer testserver.Close()
		provider.url = testserver.URL

		identifier := new(openIDIdentifier)
		identifier.reverseRouter = new(mockRouter)
		identifier.eventRouter = &flamingo.DefaultEventRouter{}
		identifier.responder = new(web.Responder)
		identifier.verifierConfigurator = append(identifier.verifierConfigurator, func(c *oidc.Config) {
			c.SkipClientIDCheck = true
			c.SkipExpiryCheck = true
			c.SkipIssuerCheck = true
		})
		identifier.oidcConfig.Claims.AccessToken = map[string]string{
			"claim-1":      "at_claim_1",
			"legacy-claim": "legacy-token-response-claim",
		}

		var err error
		identifier.provider, err = oidc.NewProvider(context.Background(), testserver.URL)
		assert.NoError(t, err)
		identifier.oauth2Config = &oauth2.Config{
			Endpoint: identifier.provider.Endpoint(),
		}

		session := web.EmptySession()
		request := web.CreateRequest(nil, session)

		identifier.createSessionCode(request, "test-callback-state")

		request.Request().URL.RawQuery = "state=test-callback-state&code=test-callback-code"
		returnCalled := false
		identifier.Callback(context.Background(), request, func(request *web.Request) *url.URL {
			returnCalled = true
			return new(url.URL)
		})
		assert.True(t, returnCalled, "the return callback was not called")

		identity, err := identifier.Identify(context.Background(), request)
		assert.NoError(t, err)

		var testClaims struct {
			Claim1      string `json:"claim-1"`
			LegacyClaim string `json:"legacy-claim"`
		}
		assert.NoError(t, identity.(Identity).AccessTokenClaims(&testClaims))
		assert.Equal(t, "at-claim-1-value", testClaims.Claim1)
		assert.Equal(t, "legacy-token-response-claim-value", testClaims.LegacyClaim)
	})

	t.Run("Test optional callback error handler", func(t *testing.T) {
		t.Parallel()
		provider := &testOidcProvider{}
		testserver := httptest.NewServer(provider)
		defer testserver.Close()
		provider.url = testserver.URL

		identifier := new(openIDIdentifier)
		identifier.reverseRouter = new(mockRouter)
		identifier.eventRouter = &flamingo.DefaultEventRouter{}
		identifier.responder = new(web.Responder)
		identifier.verifierConfigurator = append(identifier.verifierConfigurator, func(c *oidc.Config) {
			c.SkipClientIDCheck = true
			c.SkipExpiryCheck = true
			c.SkipIssuerCheck = true
		})
		identifier.oidcConfig.Claims.AccessToken = map[string]string{
			"claim-1":      "at_claim_1",
			"legacy-claim": "legacy-token-response-claim",
		}

		var err error
		identifier.provider, err = oidc.NewProvider(context.Background(), testserver.URL)
		assert.NoError(t, err)
		identifier.oauth2Config = &oauth2.Config{
			Endpoint: identifier.provider.Endpoint(),
		}

		request := web.CreateRequest(nil, web.EmptySession())
		request.Request().URL.RawQuery = "error=login_required&error_description=The%20User%20is%20not%20logged%20in&state=foo-bar"

		mockCallback := &mockCallbackErrorHandler{}
		identifier.callbackErrorHandler = mockCallback

		result := identifier.Callback(context.Background(), request, func(request *web.Request) *url.URL { return new(url.URL) })
		assert.True(t, mockCallback.Called, "the error callback handler was not called")
		assert.Equal(t, mockCallback.SuppliedError, "login_required", "the error was not passed to the callback")
		assert.Equal(t, mockCallback.SuppliedErrorDescription, "The User is not logged in", "the error callback handler was not called")
		expectedURL, _ := url.Parse("https://example.com/callback-error-handler")
		assert.Equal(t, result, &web.URLRedirectResponse{URL: expectedURL}, "Result of callback handler was ignored")
	})

	t.Run("Test add auth code option to exchange token call", func(t *testing.T) {
		t.Parallel()

		testServer := httptest.NewUnstartedServer(nil)
		tokenCalled := false

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch strings.Trim(r.URL.Path, "/") {
			case "token":
				{
					_ = r.ParseForm()
					redirectURI := r.PostForm.Get("redirect_uri")

					assert.Equal(t, redirectURI, "foobar123test")
					tokenCalled = true

					w.Header().Set("Content-type", "application/json")
					return
				}
			case ".well-known/openid-configuration":
				_, _ = fmt.Fprintf(w, `{
				"issuer": "%s",
				"token_endpoint": "%s/token",
				"jwks_uri": "%s/certs"
			}`, testServer.URL, testServer.URL, testServer.URL)
			}
		})

		testServer.Config.Handler = handler
		testServer.Start()
		defer testServer.Close()

		identifier := new(openIDIdentifier)
		identifier.reverseRouter = new(mockRouter)
		identifier.responder = new(web.Responder)

		identifier.authCodeOptionerProvider = func() []AuthCodeOptioner {
			return []AuthCodeOptioner{new(testAuthCodeOptioner)}
		}

		var err error
		identifier.provider, err = oidc.NewProvider(context.Background(), testServer.URL)
		assert.NoError(t, err)
		identifier.oauth2Config = &oauth2.Config{
			Endpoint: oauth2.Endpoint{AuthURL: "", TokenURL: testServer.URL + "/token"},
		}

		session := web.EmptySession()
		request := web.CreateRequest(nil, session)

		identifier.createSessionCode(request, "test-callback-state")

		request.Request().URL.RawQuery = "state=test-callback-state&code=test-callback-code"
		identifier.Callback(context.Background(), request, func(request *web.Request) *url.URL {
			return new(url.URL)
		})

		assert.True(t, tokenCalled)
	})
}

func Test_openIDIdentifier_RefreshIdentity(t *testing.T) {
	t.Parallel()

	var identifier auth.RequestIdentifier = &openIDIdentifier{broker: "broker"}
	session := web.EmptySession()
	session.Store("core.auth.oidc.broker.sessiondata", sessionData{
		Token: &oauth2.Token{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(time.Minute * 5),
		},
	})
	ctx := web.ContextWithSession(context.Background(), session)

	req := web.CreateRequest(nil, session)
	ctx = web.ContextWithRequest(ctx, req)

	refresher, ok := identifier.(auth.WebIdentityRefresher)
	assert.True(t, ok)

	err := refresher.RefreshIdentity(ctx, req)
	assert.NoError(t, err)

	data, found := session.Load("core.auth.oidc.broker.sessiondata")
	assert.True(t, found)

	sessiondata, ok := data.(sessionData)
	assert.True(t, ok)
	assert.Empty(t, sessiondata.Token.AccessToken)
	assert.Equal(t, "refresh-token", sessiondata.Token.RefreshToken)
}
