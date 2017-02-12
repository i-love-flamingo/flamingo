package product

import (
	"flamingo/core/backend"
	"flamingo/core/web"
	"flamingo/core/web/responder"
)

type ViewController struct {
	responder.RenderTemplate

	productservice backend.ProductServicer
}

func NewViewController(ps backend.ProductServicer) *ViewController {
	return &ViewController{
		productservice: ps,
	}
}

func (p *ViewController) Get(c web.Context) web.Response {
	//products := p.productservice.Get(c.Param1("sku"))

	return p.RenderResponse(c, "product")
}
