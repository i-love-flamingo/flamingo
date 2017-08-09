package searchperience

import (
	"flamingo/framework/dingo"
	searchdomain "flamingo/om3/search/domain"
	searchadapter "flamingo/om3/searchperience/infrastructure/search"

	productdomain "flamingo/core/product/domain"
	productadapter "flamingo/om3/searchperience/infrastructure/product"
)

// ensure types for the Ports and Adapters
var _ productdomain.ProductService = &productadapter.ProductService{}
var _ searchdomain.SearchService = &searchadapter.SearchService{}

type (
	// ProductClientModule for product client stuff
	ProductClientModule struct{}

	// SearchClientModule for searching
	SearchClientModule struct{}
)

// Configure DI
func (module *ProductClientModule) Configure(injector *dingo.Injector) {
	//injector.Bind(infrastructure.ProductClient{}).ToProvider(infrastructure.NewProductClient)
	injector.Bind((*productdomain.ProductService)(nil)).To(productadapter.ProductService{})
}

// Configure DI
func (module *SearchClientModule) Configure(injector *dingo.Injector) {
	//injector.Bind(infrastructure.SearchClient{}).ToProvider(infrastructure.NewSearchClient)
	injector.Bind((*searchdomain.SearchService)(nil)).To(searchadapter.SearchService{})
}
