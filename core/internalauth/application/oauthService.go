package application

import (
	"context"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type (
	OauthService struct {
		BaseUrl string `inject:"config:internalauth.baseurl"`
	}
)

// GetConfig returns an oauth config object
func (os *OauthService) GetConfig(TokenEndpointPath string, ClientID string, ClientSecret string, GrantType string) clientcredentials.Config {
	return clientcredentials.Config{
		ClientID: ClientID,
		ClientSecret: ClientSecret,
		TokenURL: os.BaseUrl + TokenEndpointPath,
		EndpointParams: url.Values{},
	}
}

// GetOauthToken wraps the oauth2 call to retrieve a token
func (os *OauthService) GetOauthToken(ctx context.Context, config *clientcredentials.Config) (*oauth2.Token, error) {
	token, err:= config.Token(ctx)
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
