package searchperience

import (
	categorydomain "flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	searchdomain "flamingo/core/search/domain"
	"flamingo/framework/dingo"
	categoryadapter "flamingo/om3/searchperience/infrastructure/category"
	productadapter "flamingo/om3/searchperience/infrastructure/product"
	searchadapter "flamingo/om3/searchperience/infrastructure/search"
)

// ensure types for the Ports and Adapters
var _ productdomain.ProductService = &productadapter.ProductService{}
var _ searchdomain.SearchService = &searchadapter.Service{}
var _ categorydomain.CategoryService = &categoryadapter.Service{}

type (
	// ProductClientModule for product client stuff
	ProductClientModule struct{}

	// SearchClientModule for searching
	SearchClientModule struct{}

	// CategoryClientModule for searching
	CategoryClientModule struct{}
)

// Configure DI
func (module *ProductClientModule) Configure(injector *dingo.Injector) {
	injector.Bind((*productdomain.ProductService)(nil)).To(productadapter.ProductService{})
}

// Configure DI
func (module *SearchClientModule) Configure(injector *dingo.Injector) {
	injector.Bind((*searchdomain.SearchService)(nil)).To(searchadapter.Service{})
}

// Configure DI
func (module *CategoryClientModule) Configure(injector *dingo.Injector) {
	injector.Bind((*categorydomain.CategoryService)(nil)).To(categoryadapter.Service{})
}
