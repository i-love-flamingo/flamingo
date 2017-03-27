package brand

import (
	"encoding/json"
	"io/ioutil"
	"flamingo/core/flamingo/web"
	"flamingo/akl/src/brand/models"
)

// BrandService is just mocking stuff
type BrandService struct{}

// Get returns a brand struct
func (ps *BrandService) Get(context web.Context, id string) (models.Brand) {
	var brand models.Brand

	b, _ := ioutil.ReadFile("frontend/src/mocks/brand.json")
	json.Unmarshal(b, &brand)
	brand.ID = id

	return brand
}
