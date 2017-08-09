package product

import (
	"context"
	"flamingo/om3/searchperience/infrastructure"
	"net/http"
	"net/url"
)

type (
	// ProductApiClient is a specific SearchperienceClient
	ProductApiClient struct {
		SearchperienceClient infrastructure.SearchperienceClient `inject:""`
	}
)

// Get a Product
func (bc *ProductApiClient) Get(ctx context.Context, foreignID string) (*http.Response, error) {
	query := url.Values{}
	return bc.SearchperienceClient.Request(ctx, "product/"+foreignID, query)
}
