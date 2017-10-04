package cmsblock

//go:generate go-bindata -pkg cmsblock -prefix mocks/ mocks/

import (
	"context"
	"encoding/json"

	"go.aoe.com/flamingo/core/cms/domain"
)

// FakeBlockService for CMS Blocks
type FakeBlockService struct{}

// Get returns a Block struct
func (ps *FakeBlockService) Get(ctx context.Context, name string) (*domain.Block, error) {
	var block domain.Block

	b, _ := Asset("service.cms.block.mock.json")
	json.Unmarshal(b, &block)
	block.Identifier = name

	return &block, nil
}
