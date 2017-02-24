package product

import (
	"encoding/json"
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

type ProductService struct{}

// Get returns a product struct
func (ps *ProductService) Get(id string) models.Product {
	var product models.Product

	p, _ := ioutil.ReadFile("frontend/src/mocks/product.json")
	json.Unmarshal(p, &product)

	product.Id = id
	product.Name = fmt.Sprintf("%s %s", nameprefix[rand.Intn(len(nameprefix))], namesuffix[rand.Intn(len(namesuffix))])

	return product
}

// GetByIDList returns a struct of Product Models identified by given Skus
func (ps *ProductService) GetByIDList(skus []string) []models.Product {
	var products = make([]models.Product, len(skus))

	for i, sku := range skus {
		products[i] = ps.Get(sku)
	}

	return products
}
