package internalmock

import (
	"flamingo/core/app"
	"flamingo/core/packages/internalmock/brand"
	"flamingo/core/packages/internalmock/product"
)

func Register(r *app.ServiceContainer) {
	r.Register(new(product.ProductService))
	r.Register(new(brand.BrandService))
}
