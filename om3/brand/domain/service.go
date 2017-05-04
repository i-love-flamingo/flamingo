package domain

import "context"

type (
	// BrandService will be used to get brands
	BrandService interface {
		Get(context context.Context, ID string) (brand *Brand, err error)
	}

	BrandNotFound struct {
		Name string
	}
)

func (b BrandNotFound) Error() string {
	return "Brand " + b.Name + " Not Found"
}
