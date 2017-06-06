package infrastructure

import (
	"context"
	"net/http"
)

type (
	// APIClient requests masterdataportal api
	APIClient struct {
		BaseURL string `inject:"config:masterdataportal.baseurl"`
		common  *APIClient
	}

	// BrandsClient is a specific APIClient
	BrandsClient APIClient
)

func (ac *APIClient) request(ctx context.Context, p string) (*http.Response, error) {
	req, err := http.NewRequest("GET", ac.BaseURL+p, nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)
	return http.DefaultClient.Do(req)
}

// NewBrandsClient creates a BrandsClient from an APIClient
func NewBrandsClient(ac *APIClient) *BrandsClient {
	ac.common = ac
	return (*BrandsClient)(ac)
}

// Get a Brand
func (bc *BrandsClient) Get(ctx context.Context, ID string) (*http.Response, error) {
	return bc.common.request(ctx, "/brands/"+ID)
}
