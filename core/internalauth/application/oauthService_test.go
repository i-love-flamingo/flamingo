package application_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/core/internalauth/application"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/clientcredentials"
)

func TestOauthService(t *testing.T) {
	service := new(application.OauthService)
	service.Inject(&struct {
		BaseURL string "inject:\"config:core.internalauth.baseurl\""
	}{
		BaseURL: "http://example.com/",
	})

	t.Run("GetConfig", func(t *testing.T) {
		cfg := service.GetConfig("/path", "client-id", "client-secret", "grant-type")
		assert.Equal(t, "client-id", cfg.ClientID)
		assert.Equal(t, "client-secret", cfg.ClientSecret)
		assert.Equal(t, url.Values{}, cfg.EndpointParams)
		assert.Equal(t, "http://example.com/path", cfg.TokenURL)
	})

	t.Run("GetOauthToken", func(t *testing.T) {
		_, err := service.GetOauthToken(context.Background(), &clientcredentials.Config{})
		assert.Error(t, err)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			fmt.Fprintf(w, `{"access_token": "test"}`)
		}))
		defer server.Close()
		token, err := service.GetOauthToken(context.Background(), &clientcredentials.Config{
			TokenURL: server.URL,
		})
		assert.NoError(t, err)
		assert.Equal(t, "test", token.AccessToken)
	})

	t.Run("GetClaimsFromToken", func(t *testing.T) {
		claims := service.GetClaimsFromToken("token")
		assert.Empty(t, claims)

		claims = service.GetClaimsFromToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZXN0IjoidmFsdWUifQ.signature")
		assert.NotEmpty(t, claims)
		assert.Equal(t, "value", claims["test"])
	})
}
