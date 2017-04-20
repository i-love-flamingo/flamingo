package domain

import "context"

// BrandService will be used to get brands
type (
	BrandService interface {
		Get(context context.Context, ID string) *Brand
	}
)
