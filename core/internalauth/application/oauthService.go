package application

import (
	"context"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type (
	OauthService struct {
		BaseUrl string `inject:"config:internalauth.baseurl"`
	}
)

func (os *OauthService) GetConfig(TokenEndpointPath string, ClientID string, ClientSecret string, GrantType string) clientcredentials.Config {
	return clientcredentials.Config{
		ClientID: ClientID,
		ClientSecret: ClientSecret,
		TokenURL: os.BaseUrl + TokenEndpointPath,
		EndpointParams: url.Values{},
	}
}

func (os *OauthService) GetOauthToken(ctx context.Context, config *clientcredentials.Config) (*oauth2.Token, error) {
	token, err:= config.Token(ctx)
	if err != nil {
		return nil, err
	}

	return token, nil
}
