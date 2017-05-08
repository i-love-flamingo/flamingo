package searchperience

import (
	"flamingo/core/product/domain"
	"flamingo/framework/dingo"
	"flamingo/om3/searchperience/infrastructure"
)

// ensure types
var _ domain.ProductService = &infrastructure.ProductService{}

// ProductClientModule for product client stuff
type ProductClientModule struct{}

// Configure DI
func (module *ProductClientModule) Configure(injector *dingo.Injector) {
	injector.Bind(infrastructure.ProductClient{}).ToProvider(infrastructure.NewProductClient)
	injector.Bind((*domain.ProductService)(nil)).To(infrastructure.ProductService{})
}
