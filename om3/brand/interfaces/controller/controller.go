package controller

import (
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"flamingo/om3/brand/domain"

	"github.com/pkg/errors"
)

type (
	// ViewController demonstrates a brand view controller
	ViewController struct {
		responder.RenderAware `inject:""`
		responder.ErrorAware  `inject:""`
		domain.BrandService   `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Brand *domain.Brand
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	brand, err := vc.BrandService.Get(c, c.MustParam1("uid"))

	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.BrandNotFound:
			return vc.ErrorNotFound(c, err)

		default:
			return vc.Error(c, err)
		}
	}

	return vc.Render(c, "brand/view", ViewData{Brand: brand})
}
