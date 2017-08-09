package brand

import (
	"context"
	"encoding/json"
	"flamingo/om3/brand/domain"
	"fmt"
	"io/ioutil"
)

// FakeBrandService is just mocking stuff
type FakeBrandService struct{}

// Get returns a brand struct
func (ps *FakeBrandService) Get(context context.Context, ID string) (*domain.Brand, error) {
	var brand domain.Brand

	b, _ := ioutil.ReadFile("frontend/src/mock/service.brand.mock.json")
	json.Unmarshal(b, &brand)
	brand.ID = ID
	fmt.Println("fake brand service called")
	fmt.Println(brand)
	return &brand, nil
}
