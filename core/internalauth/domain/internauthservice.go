package domain

import (
	"context"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// InternalAuthService interface for internal oauth clients
// todo necessary?
type InternalAuthService interface {
	GetConfig(tokenEndpointPath string, clientID string, clientSecret string, grantType string) clientcredentials.Config
	GetOauthToken(ctx context.Context, config *clientcredentials.Config) (*oauth2.Token, error)
	GetClaimsFromToken(tokenString string) jwt.MapClaims
}
