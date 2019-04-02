package domain

import "context"

// PageService for cms page retrieval
type PageService interface {
	Get(context.Context, string) (*Page, error)
}
