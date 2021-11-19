package oauth

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
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

func TestParallelStateRaceConditions(t *testing.T) {
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

	t.Run("test states", func(t *testing.T) {
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

	t.Run("test timeshift", func(t *testing.T) {
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
}
