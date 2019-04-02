package domain

import "context"

// BlockService defines the cms block service
type BlockService interface {
	Get(context.Context, string) (*Block, error)
}
