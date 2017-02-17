package controller

import (
	"flamingo/core/core/app/web"
	"flamingo/core/core/app/web/responder"
	"flamingo/core/core/product/interfaces"
)

type ViewController struct {
	*responder.RenderAware `inject:""`

	interfaces.ProductService `inject:""`
}

func (p *ViewController) Get(c web.Context) web.Response {
	product := p.ProductService.Get(c.Param1("sku"))

	return p.Render(c, "pages/product/view", product)
}
