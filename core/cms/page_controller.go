package cms

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/flamingo/web/responder"
)

// PageController is a demo cms page view controller
type PageController struct {
	*responder.RenderAware `inject:""`

	//pageservice interfaces.PageService
}

// Get renders specific page routes.
func (pc *PageController) Get(c web.Context) web.Response {
	return pc.Render(c, "pages/"+c.Param1("name"), nil)
}
