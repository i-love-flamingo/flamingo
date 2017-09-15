package category

import (
	"context"
	"flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	"flamingo/om3/fakeservices/product"
)

// FakeBlockService for CMS Blocks
type FakeCategoryService struct{}

// Get returns a category struct
func (cs *FakeCategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	return domain.Category{
		Categories: []*domain.Category{
			{
				Name: "Test2",
				Code: "test2",
			},
			{
				Name: "Test3",
				Code: "test3",
			},
			{
				Name: "Test4",
				Code: "test4",
			},
		},
		Name: "Test",
		Code: "test",
	}, nil
}

func (cs *FakeCategoryService) GetProducts(ctx context.Context, categoryCode string) ([]productdomain.BasicProduct, error) {
	ps := new(product.FakeProductService)
	p, _ := ps.Get(ctx, "fake_simple")

	return []productdomain.BasicProduct{
		p,
		p,
		p,
		p,
		p,
	}, nil
}
