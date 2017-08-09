package search

import (
	"context"
	"flamingo/om3/searchperience/infrastructure"
	"net/http"
	"net/url"
)

type (
	// ProductClient is a specific SearchperienceClient
	SearchClient struct {
		SearchperienceClient infrastructure.SearchperienceClient `inject:""`
	}
)

// Search gets a search result
func (bc *SearchClient) Search(ctx context.Context, query url.Values) (*http.Response, error) {
	return bc.SearchperienceClient.Request(ctx, "search", query)
}
