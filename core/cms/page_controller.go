package cms

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

// PageController is a demo cms page view controller
type PageController struct {
	*responder.RenderAware `inject:""`

	//pageservice interfaces.PageService
}

// Get renders specific page routes.
func (pc *PageController) Get(c web.Context) web.Response {
	return pc.Render(c, "pages/"+c.MustParam1("name"), nil)
}
