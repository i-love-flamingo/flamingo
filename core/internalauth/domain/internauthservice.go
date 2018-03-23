package domain

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type (
	InternalAuthService interface {
		GetConfig(TokenEndpointPath string, ClientID string, ClientSecret string, GrantType string) clientcredentials.Config
		GetOauthToken(ctx context.Context, config *clientcredentials.Config) (*oauth2.Token, error)
		GetClaimsFromToken(tokenString string) jwt.MapClaims
	}
)
