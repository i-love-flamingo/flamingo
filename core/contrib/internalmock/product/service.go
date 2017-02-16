package product

import (
	"flamingo/core/core/product/models"
	"fmt"
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

func (ps *ProductService) Get(id string) models.Product {
	return models.Product{
		ID:          id,
		Name:        fmt.Sprintf("%s %s", nameprefix[rand.Intn(len(nameprefix))], namesuffix[rand.Intn(len(namesuffix))]),
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
}

// GetBySkuList gbsl
func (ps *ProductService) GetByIDList(skus []string) []models.Product {
	var products = make([]models.Product, len(skus))

	for i, sku := range skus {
		products[i] = ps.Get(sku)
	}

	return products
}
