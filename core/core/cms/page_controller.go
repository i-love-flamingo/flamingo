package cms

import (
	"flamingo/core/core/app/web"
	"flamingo/core/core/app/web/responder"
	"flamingo/core/core/cms/interfaces"
)

type PageController struct {
	responder.RenderAware

	pageservice interfaces.PageService
}

func (i *PageController) Copy() interface{} {
	var ic PageController = *i
	return &ic
}

func (pc *PageController) Get(c web.Context) web.Response {
	return pc.Render(c, "pages/"+c.Param1("name"))
}
