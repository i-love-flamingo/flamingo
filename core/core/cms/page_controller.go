package cms

import (
	"flamingo/core/web"
	"flamingo/core/web/responder"
)

type PageController struct {
	responder.RenderAware
}

func (pc *PageController) Get(c web.Context) web.Response {
	return pc.Render(c, "pages/"+c.Param1("name"))
}
