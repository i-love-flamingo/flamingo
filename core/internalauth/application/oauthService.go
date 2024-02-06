package application

import (
	"context"
	"net/url"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OauthService for internal direct token grant
type OauthService struct {
	baseURL string
}

// Inject configuration
func (os *OauthService) Inject(config *struct {
	BaseURL string `inject:"config:core.internalauth.baseurl"`
}) {
	os.baseURL = config.BaseURL
}

// GetConfig returns an oauth config object
func (os *OauthService) GetConfig(TokenEndpointPath string, ClientID string, ClientSecret string, _ string) clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:       ClientID,
		ClientSecret:   ClientSecret,
		TokenURL:       strings.TrimRight(os.baseURL, "/") + "/" + strings.TrimLeft(TokenEndpointPath, "/"),
		EndpointParams: url.Values{},
	}
}

// GetOauthToken wraps the oauth2 call to retrieve a token
func (os *OauthService) GetOauthToken(ctx context.Context, config *clientcredentials.Config) (*oauth2.Token, error) {
	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetClaimsFromToken is a "fix" for the oauth2 libs inability to decode extra data from the token
func (os *OauthService) GetClaimsFromToken(tokenString string) jwt.MapClaims {
	claims := jwt.MapClaims{}

	jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(nil), nil
	})

	return claims
}
