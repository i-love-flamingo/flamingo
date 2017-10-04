package fakeservices

import (
	"go.aoe.com/flamingo/om3/fakeservices/brand"
	"go.aoe.com/flamingo/om3/fakeservices/category"
	"go.aoe.com/flamingo/om3/fakeservices/cmsblock"
	"go.aoe.com/flamingo/om3/fakeservices/cmspage"
	"go.aoe.com/flamingo/om3/fakeservices/product"
	"go.aoe.com/flamingo/om3/fakeservices/search"

	categorydomain "go.aoe.com/flamingo/core/category/domain"
	cmsdomain "go.aoe.com/flamingo/core/cms/domain"
	productdomain "go.aoe.com/flamingo/core/product/domain"
	searchdomain "go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	branddomain "go.aoe.com/flamingo/om3/brand/domain"
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

	if v, ok := module.Config["category"].(bool); v && ok {
		injector.Override((*categorydomain.CategoryService)(nil), "").To(category.FakeCategoryService{})
	}
}

func (module *Module) DefaultConfig() config.Map {
	return config.Map{
		"fakeservices": config.Map{},
	}
}
