package controller

import (
	"go.aoe.com/flamingo/core/cms/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// DataController for `get("cms.block", ...)` requests
	DataController struct {
		BlockService domain.BlockService `inject:""`
	}
)

// Data controller for blocks
func (vc *DataController) Data(c web.Context) interface{} {
	block, _ := vc.BlockService.Get(c, c.MustParam1("block"))
	return block
}
