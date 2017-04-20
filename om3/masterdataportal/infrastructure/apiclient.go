package infrastructure

import (
	"context"
	"flamingo/framework/web"
	"net/http"
)

type (
	// ApiClient requests masterdataportal api
	ApiClient struct {
		BaseURL string `inject:"config:masterdataportal.baseurl"`
		common  *ApiClient
	}

	// BrandsClient is a specific ApiClient
	BrandsClient ApiClient
)

func (ac *ApiClient) request(ctx context.Context, p string) *http.Response {
	req, err := http.NewRequest("GET", ac.BaseURL+p, nil)
	if err != nil {
		panic(err)
	}
	if ctx, ok := ctx.(web.Context); ok {
		defer ctx.Profile("masterdataportal", "GET "+ac.BaseURL+p)()
		req.Header.Add("X-Request-ID", ctx.ID())
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return r
}

// NewBrandsClient creates a BrandsClient from an ApiClient
func NewBrandsClient(ac *ApiClient) *BrandsClient {
	ac.common = ac
	return (*BrandsClient)(ac)
}

// Get a Brand
func (bc *BrandsClient) Get(ctx context.Context, ID string) *http.Response {
	return bc.common.request(ctx, "/brands/"+ID)
}
