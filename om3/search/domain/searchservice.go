package domain

import (
	"flamingo/framework/web"
)

// ProductService interface
type SearchService interface {
	Search(ctx web.Context, query string) (*SearchResult, error)
}
