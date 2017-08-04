package infrastructure

import (
	"context"
	"net/http"
	"net/url"
)

type (
	// SearchperienceClient talks to searchperience
	SearchperienceClient struct {
		BaseURL string `inject:"config:searchperience.frontend.baseurl"`
		common  *SearchperienceClient
	}

	// ProductClient is a specific SearchperienceClient
	ProductClient SearchperienceClient

	// SearchClient is a specific SearchperienceClient
	SearchClient SearchperienceClient
)

func (ac *SearchperienceClient) request(ctx context.Context, path string, query url.Values) (*http.Response, error) {

	u, _ := url.Parse(ac.BaseURL)
	u.Path += path
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)
	return http.DefaultClient.Do(req)
}

// NewProductClient provider
func NewProductClient(ac *SearchperienceClient) *ProductClient {
	ac.common = ac
	return (*ProductClient)(ac)
}

// Get a Product
func (bc *ProductClient) Get(ctx context.Context, foreignID string) (*http.Response, error) {
	query := url.Values{}
	return bc.common.request(ctx, "product/"+foreignID, query)
}

// NewSearchClient provider
func NewSearchClient(ac *SearchperienceClient) *SearchClient {
	ac.common = ac
	return (*SearchClient)(ac)
}

// Search gets a search result
func (bc *SearchClient) Search(ctx context.Context, query url.Values) (*http.Response, error) {
	return bc.common.request(ctx, "search", query)
}
