package interfaces

import (
	"flamingo/framework/web"
	"flamingo/om3/brand/models"
)

// BrandService will be used to get brands
type BrandService interface {
	Get(web.Context, string) (models.Brand)
}
