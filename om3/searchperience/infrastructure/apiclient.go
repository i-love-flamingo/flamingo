package infrastructure

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

type (
	// SearchperienceClient talks to searchperience
	SearchperienceClient struct {
		BaseURL string `inject:"config:searchperience.frontend.baseurl"`
		common  *SearchperienceClient
	}
)

func (ac *SearchperienceClient) Request(ctx context.Context, path string, query url.Values) (*http.Response, error) {
	u, _ := url.Parse(ac.BaseURL)
	u.Path += path
	u.RawQuery = query.Encode()
	log.Printf("Searchperience Call to %v", u)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
