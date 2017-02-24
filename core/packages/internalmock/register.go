package internalmock

import (
	"flamingo/core/flamingo/service_container"
	"flamingo/core/packages/internalmock/brand"
	"flamingo/core/packages/internalmock/product"
)

// Register Services for internalmock package
func Register(r *service_container.ServiceContainer) {
	r.Register(new(product.ProductService))
	r.Register(new(brand.BrandService))
}
