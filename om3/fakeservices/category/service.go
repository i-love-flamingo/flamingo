package category

import (
	"context"
	"flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	"flamingo/om3/fakeservices/product"
)

type (
	// FakeBlockService for CMS Blocks
	FakeCategoryService struct{}

	fakeCategory struct {
		code       string
		name       string
		categories []domain.Category
	}
)

// Code return the category code
func (f *fakeCategory) Code() string {
	return f.code
}

// Name returns the category name
func (f *fakeCategory) Name() string {
	return f.name
}

// Categories returns a list of child categories
func (f *fakeCategory) Categories() []domain.Category {
	return f.categories
}

// Get returns a category struct
func (cs *FakeCategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	return &fakeCategory{
		name: "Test",
		code: "test",
		categories: []domain.Category{
			&fakeCategory{
				name: "Test2",
				code: "test2",
			},
			&fakeCategory{
				name: "Test3",
				code: "test3",
			},
			&fakeCategory{
				name: "Test4",
				code: "test4",
			},
		},
	}, nil
}

func (cs *FakeCategoryService) GetProducts(ctx context.Context, categoryCode string) ([]productdomain.BasicProduct, error) {
	return []productdomain.BasicProduct{
		product.FakeSimple("product-1"),
		product.FakeSimple("product-2"),
		product.FakeSimple("product-3"),
		product.FakeSimple("product-4"),
		product.FakeSimple("product-5"),
	}, nil
}
