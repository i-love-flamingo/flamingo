package internalmock

import (
	"flamingo/core/contrib/internalmock/brand"
	"flamingo/core/contrib/internalmock/product"
	"flamingo/core/core/app"
)

func Register(r *app.ServiceContainer) {
	r.Register(new(product.ProductService))
	r.Register(new(brand.BrandService))
}
