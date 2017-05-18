package searchperience

import (
	"flamingo/core/product/domain"
	"flamingo/framework/dingo"
	searchdomain "flamingo/om3/search/domain"
	"flamingo/om3/searchperience/infrastructure"
)

// ensure types
var _ domain.ProductService = &infrastructure.ProductService{}

type (
	// ProductClientModule for product client stuff
	ProductClientModule struct{}

	// SearchClientModule for searching
	SearchClientModule struct{}
)

// Configure DI
func (module *ProductClientModule) Configure(injector *dingo.Injector) {
	injector.Bind(infrastructure.ProductClient{}).ToProvider(infrastructure.NewProductClient)
	injector.Bind((*domain.ProductService)(nil)).To(infrastructure.ProductService{})
}

// Configure DI
func (module *SearchClientModule) Configure(injector *dingo.Injector) {
	injector.Bind(infrastructure.SearchClient{}).ToProvider(infrastructure.NewSearchClient)
	injector.Bind((*searchdomain.SearchService)(nil)).To(infrastructure.SearchService{})
}
