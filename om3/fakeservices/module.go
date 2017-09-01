package fakeservices

import (
	"flamingo/om3/fakeservices/brand"
	"flamingo/om3/fakeservices/cmsblock"
	"flamingo/om3/fakeservices/cmspage"
	"flamingo/om3/fakeservices/product"
	"flamingo/om3/fakeservices/search"

	cmsdomain "flamingo/core/cms/domain"
	productdomain "flamingo/core/product/domain"
	"flamingo/framework/config"
	"flamingo/framework/dingo"
	branddomain "flamingo/om3/brand/domain"
	searchdomain "flamingo/om3/search/domain"
)

// Module for AKL internalmock configuration
type Module struct {
	Config config.Map `inject:"config:fakeservices"`
}

// Configure DI
func (module *Module) Configure(injector *dingo.Injector) {
	if v, ok := module.Config["brand"].(bool); v && ok {
		injector.Override((*branddomain.BrandService)(nil), "").To(brand.FakeBrandService{})
	}

	if v, ok := module.Config["product"].(bool); v && ok {
		injector.Override((*productdomain.ProductService)(nil), "").To(product.FakeProductService{})
	}

	if v, ok := module.Config["search"].(bool); v && ok {
		injector.Override((*searchdomain.SearchService)(nil), "").To(search.FakeSearchService{})
	}

	if v, ok := module.Config["cmspage"].(bool); v && ok {
		injector.Override((*cmsdomain.PageService)(nil), "").To(cmspage.FakePageService{})
	}

	if v, ok := module.Config["cmsblock"].(bool); v && ok {
		injector.Override((*cmsdomain.BlockService)(nil), "").To(cmsblock.FakeBlockService{})
	}
}
