package searchperience

import (
	categorydomain "go.aoe.com/flamingo/core/category/domain"
	productdomain "go.aoe.com/flamingo/core/product/domain"
	searchdomain "go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/framework/dingo"
	categoryadapter "go.aoe.com/flamingo/om3/searchperience/infrastructure/category"
	productadapter "go.aoe.com/flamingo/om3/searchperience/infrastructure/product"
	searchadapter "go.aoe.com/flamingo/om3/searchperience/infrastructure/search"
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
