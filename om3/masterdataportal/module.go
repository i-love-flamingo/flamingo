package masterdataportal

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/om3/brand/domain"
	"go.aoe.com/flamingo/om3/masterdataportal/infrastructure"
)

// check types at compile time ;)
var _ domain.BrandService = &infrastructure.BrandService{}

// Module is our MasterDataPortal Module
type Module struct{}

// Configure DI
func (module *Module) Configure(injector *dingo.Injector) {
	injector.Bind(infrastructure.BrandsClient{}).ToProvider(infrastructure.NewBrandsClient)
	injector.Bind((*domain.BrandService)(nil)).To(infrastructure.BrandService{})
}
