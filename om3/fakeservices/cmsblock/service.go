package cmsblock

import (
	"context"
	"encoding/json"
	"flamingo/core/cms/domain"
	"io/ioutil"
)

// FakeBlockService for CMS Blocks
type FakeBlockService struct{}

// Get returns a Block struct
func (ps *FakeBlockService) Get(ctx context.Context, name string) (*domain.Block, error) {
	var block domain.Block

	b, _ := ioutil.ReadFile("src/fakeservices/cmsblock/service.cms.block.mock.json")
	json.Unmarshal(b, &block)
	block.Identifier = name

	return &block, nil
}
