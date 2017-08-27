package controller

import (
	"context"
	"flamingo/core/cms/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	MockBlockService struct{}
)

func (m *MockBlockService) Get(ctx context.Context, identifier string) (*domain.Block, error) {
	return &domain.Block{
		Identifier: identifier,
	}, nil
}

func TestDataController_Data(t *testing.T) {
	dc := &DataController{
		BlockService: new(MockBlockService),
	}
	ctx := web.NewContext()

	ctx.LoadParams(router.P{"block": "test"})
	block, ok := dc.Data(ctx).(*domain.Block)

	assert.True(t, ok, "Result is not a block")
	assert.NotNil(t, block, "Block is nil")
	assert.Equal(t, "test", block.Identifier, "Wrong identifier %q, expected test", block.Identifier)
}
