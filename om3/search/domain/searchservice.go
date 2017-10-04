package domain

import (
	"net/url"

	"go.aoe.com/flamingo/framework/web"
)

// SearchService interface
type SearchService interface {
	Search(ctx web.Context, query url.Values) (*SearchResult, error)
}
