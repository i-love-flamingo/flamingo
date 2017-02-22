package internalmock

import (
	"flamingo/core/flamingo/service_container"
	"flamingo/core/packages/internalmock/brand"
	"flamingo/core/packages/internalmock/product"
)

func Register(r *service_container.ServiceContainer) {
	r.Register(new(product.ProductService))
	r.Register(new(brand.BrandService))
}
