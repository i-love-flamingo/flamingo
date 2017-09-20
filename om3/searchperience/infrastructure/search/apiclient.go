package search

import (
	"context"
	"flamingo/om3/searchperience/infrastructure"
	"net/http"
	"net/url"
)

type (
	// ProductClient is a specific SearchperienceClient
	Client struct {
		SearchperienceClient infrastructure.SearchperienceClient `inject:""`
	}
)

// Search gets a search result
func (bc *Client) Search(ctx context.Context, query url.Values) (*http.Response, error) {
	return bc.SearchperienceClient.Request(ctx, "search", query)
}

// Category product listing request
func (bc *Client) Category(ctx context.Context, category string, query url.Values) (*http.Response, error) {
	return bc.SearchperienceClient.Request(ctx, "product/category/"+category, query)
}
