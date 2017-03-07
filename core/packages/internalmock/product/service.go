package product

import (
	"encoding/json"
	"flamingo/core/flamingo/web"
	"flamingo/core/product/models"
	"fmt"
	"io/ioutil"
	"math/rand"
)

var nameprefix = [...]string{
	"Cool",
	"Fancy",
	"Modern",
	"Soft",
	"Stylish",
	"Hip",
	"Green",
	"Red",
	"Super Fancy",
	"Sparkling",
	"Glittery",
}

var namesuffix = [...]string{
	"Bag",
	"Hat",
	"Shirt",
	"Top",
	"Jeans",
	"Pants",
	"BackPack",
}

// ProductService is just mocking stuff
type ProductService struct{}

// Get returns a product struct
func (ps *ProductService) Get(context web.Context, id string) (models.Product, models.AppError) {
	defer context.Profile("service", "get product "+id)()
	var product models.Product

	e := models.AppError{}

	p, _ := ioutil.ReadFile("frontend/src/mocks/product.json")
	json.Unmarshal(p, &product)

	product.Id = id
	product.Name = fmt.Sprintf("%s %s", nameprefix[rand.Intn(len(nameprefix))], namesuffix[rand.Intn(len(namesuffix))])

	return product, e
}

// GetByIDList returns a struct of Product Models identified by given Skus
func (ps *ProductService) GetByIDList(context web.Context, skus []string) []models.Product {
	defer context.Profile("service", "get product list")()
	var products = make([]models.Product, len(skus))

	for i, sku := range skus {
		products[i], _ = ps.Get(context, sku)
	}

	return products
}
