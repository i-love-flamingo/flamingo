package cms

import (
	"flamingo/core/web"
	"flamingo/core/web/responder"
)

type PageController struct {
	responder.RenderTemplate
}

func (pc *PageController) Get(c web.Context) web.Response {
	return pc.RenderResponse(c, "pages/"+c.Param1("name"))
}
