package infrastructure

import (
	"context"
	"flamingo/framework/web"
	"net/http"
)

type (
	// SearchperienceClient talks to searchperience
	SearchperienceClient struct {
		BaseURL string `inject:"config:searchperience.frontend.baseurl"`
		common  *SearchperienceClient
	}

	// ProductClient is a specific SearchperienceClient
	ProductClient SearchperienceClient
	SearchClient SearchperienceClient
)

func (ac *SearchperienceClient) request(ctx context.Context, p string) (*http.Response, error) {
	req, err := http.NewRequest("GET", ac.BaseURL+p, nil)
	if err != nil {
		panic(err)
	}
	if ctx, ok := ctx.(web.Context); ok {
		defer ctx.Profile("searchperience", "GET "+ac.BaseURL+p)()
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
	return bc.common.request(ctx, "document?type=product&foreignId="+foreignID)
}

// SearchClient provider
func NewSearchClient(ac *SearchperienceClient) *SearchClient {
	ac.common = ac
	return (*SearchClient)(ac)
}

// Get a search result
func (bc *SearchClient) Search(ctx context.Context, query string) (*http.Response, error) {
	return bc.common.request(ctx, "search?q="+query)
}
