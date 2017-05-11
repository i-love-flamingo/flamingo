package domain

import (
	"flamingo/framework/web"
	"net/url"
)

// SearchService interface
type SearchService interface {
	Search(ctx web.Context, query url.Values) (*SearchResult, error)
}
