package internalmock

import (
	"flamingo/core/packages/internalmock/brand"
	"flamingo/core/packages/internalmock/product"
	"flamingo/core/flamingo"
)

func Register(r *flamingo.ServiceContainer) {
	r.Register(new(product.ProductService))
	r.Register(new(brand.BrandService))
}
