package search

import (
	"context"
	"net/http"
	"net/url"

	"go.aoe.com/flamingo/om3/searchperience/infrastructure"
)

type (
	// ProductClient is a specific SearchperienceClient
	Client struct {
		SearchperienceClient infrastructure.SearchperienceClient `inject:""`
		Locale               string                              `inject:"config:locale"`
		Channel              string                              `inject:"config:searchperience.frontend.channel"`
	}
)

// Search gets a search result
func (bc *Client) Search(ctx context.Context, query url.Values) (*http.Response, error) {
	return bc.SearchperienceClient.Request(ctx, "search", query)
}

// Category product listing request
func (bc *Client) Category(ctx context.Context, category string, query url.Values) (*http.Response, error) {
	return bc.SearchperienceClient.Request(ctx, "product/category/"+category+"-"+bc.Locale+"-"+bc.Channel, query)
}
