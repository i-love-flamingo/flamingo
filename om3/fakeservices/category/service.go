package category

import (
	"context"
	"strconv"

	"go.aoe.com/flamingo/core/category/domain"
	productdomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/om3/fakeservices/product"
)

type (
	// FakeBlockService for CMS Blocks
	FakeCategoryService struct{}

	fakeCategory struct {
		code       string
		name       string
		categories []domain.Category
		active     bool
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

// Active indicator
func (f *fakeCategory) Active() bool {
	return f.active
}

// Get returns a category struct
func (cs *FakeCategoryService) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	r := &fakeCategory{
		name: "Test",
		code: "test",
		categories: []domain.Category{
			&fakeCategory{
				name: "Sub 1",
				code: "sub1",
				categories: []domain.Category{
					&fakeCategory{
						name: "Sub 1 / 1",
						code: "sub11",
					},
					&fakeCategory{
						name: "Sub 1 / 2",
						code: "sub12",
					},
					&fakeCategory{
						name: "Sub 1 / 3",
						code: "sub13",
					},
				},
			},
			&fakeCategory{
				name: "Sub 2",
				code: "sub2",
				categories: []domain.Category{
					&fakeCategory{
						name: "Sub 2 / 1",
						code: "sub21",
					},
					&fakeCategory{
						name: "Sub 2 / 2",
						code: "sub22",
					},
					&fakeCategory{
						name: "Sub 2 / 3",
						code: "sub23",
					},
				},
			},
			&fakeCategory{
				name: "Sub 3",
				code: "sub3",
				categories: []domain.Category{
					&fakeCategory{
						name: "Sub 3 / 1",
						code: "sub31",
					},
					&fakeCategory{
						name: "Sub 3 / 2",
						code: "sub32",
					},
					&fakeCategory{
						name: "Sub 3 / 3",
						code: "sub33",
					},
				},
			},
		},
	}
	markActive(r, categoryCode)
	return r, nil
}

func markActive(sc *fakeCategory, categoryCode string) (marked bool) {
	for _, sub := range sc.categories {
		if markActive(sub.(*fakeCategory), categoryCode) {
			sc.active = true
			return true
		}
	}
	if sc.code == categoryCode {
		sc.active = true
		return true
	}
	return
}

func (cs *FakeCategoryService) GetProducts(ctx context.Context, categoryCode string) ([]productdomain.BasicProduct, error) {
	products := make([]productdomain.BasicProduct, 30)

	for i := 1; i <= 30; i++ {
		products[i] = product.FakeSimple("product-" + strconv.Itoa(i))
	}

	return products, nil
}
