package cms

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/flamingo/web/responder"
)

type PageController struct {
	*responder.RenderAware `inject:""`

	//pageservice interfaces.PageService
}

func (pc *PageController) Get(c web.Context) web.Response {
	return pc.Render(c, "pages/"+c.Param1("name"), nil)
}
