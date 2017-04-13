package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"flamingo/om3/brand/interfaces"
	"flamingo/om3/brand/models"
)

type (
	// ViewController demonstrates a brand view controller
	ViewController struct {
		*responder.ErrorAware     `inject:""`
		*responder.RenderAware    `inject:""`
		interfaces.BrandService   `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Brand models.Brand
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	brand := vc.BrandService.Get(c, c.Param1("uid"))
	return vc.Render(c, "pages/brand/view", ViewData{Brand: brand})
}
