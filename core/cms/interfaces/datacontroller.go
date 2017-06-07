package interfaces

import (
	"flamingo/core/cms/domain"
	"flamingo/framework/web"
)

type (
	// DataController for `get("cms.block", ...)` requests
	DataController struct {
		BlockService domain.BlockService `inject:""`
		DevMode bool `inject:"config:debug.mode"`
	}
)

// Data controller for blocks
func (vc *DataController) Data(c web.Context) interface{} {
	block, err := vc.BlockService.Get(c, c.MustParam1("block"))

	if err != nil && vc.DevMode {
		return domain.Block{
			Title: "[Block " + c.MustParam1("block") +  " not found]",
			Content: "[Block " + c.MustParam1("block") +  " not found]",
		}
	}

	return block
}
