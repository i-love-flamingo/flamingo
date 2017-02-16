package product

import (
	"flamingo/core/core/app"
	"flamingo/core/core/product/controller"
)

func Register(r *app.Registrator) {
	r.Handle("product.view", new(controller.ViewController))
	r.Route("/product/{sku}", "product.view")
}
