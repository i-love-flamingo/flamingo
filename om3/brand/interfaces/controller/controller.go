package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"flamingo/om3/brand/domain"
)

type (
	// ViewController demonstrates a brand view controller
	ViewController struct {
		*responder.RenderAware `inject:""`
		domain.BrandService    `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Brand *domain.Brand
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	return vc.Render(c, "pages/brand/view", ViewData{Brand: vc.BrandService.Get(c, c.Param1("uid"))})
}
