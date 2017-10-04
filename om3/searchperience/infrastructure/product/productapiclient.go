package product

import (
	"context"
	"net/http"
	"net/url"

	"go.aoe.com/flamingo/om3/searchperience/infrastructure"
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
