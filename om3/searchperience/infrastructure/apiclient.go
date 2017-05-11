package infrastructure

import (
	"context"
	"flamingo/framework/web"
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
	SearchClient  SearchperienceClient
)

func (ac *SearchperienceClient) request(ctx context.Context, path string, query url.Values) (*http.Response, error) {

	u, _ := url.Parse(ac.BaseURL)
	u.Path += path
	u.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}
	if ctx, ok := ctx.(web.Context); ok {
		defer ctx.Profile("searchperience", "GET "+u.String())()
		req.Header.Add("X-Request-ID", ctx.ID())
	}
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
	query.Set("type", "product")
	query.Set("foreignId", foreignID)
	return bc.common.request(ctx, "document", query)
}

// SearchClient provider
func NewSearchClient(ac *SearchperienceClient) *SearchClient {
	ac.common = ac
	return (*SearchClient)(ac)
}

// Get a search result
func (bc *SearchClient) Search(ctx context.Context, query url.Values) (*http.Response, error) {
	return bc.common.request(ctx, "search", query)
}
