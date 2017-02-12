package mock

import (
	"fmt"
	"flamingo/core/backend"
	"math/rand"
)

// ProductService ps
type ProductService struct {
	//profiler *profiler.Profile
}

// NewProductService nps
func NewProductService(dsn string) backend.ProductServicer {
	return ProductService{}
}

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

// WithProfiler make it profileable
/*func (ps ProductService) WithProfiler(p *profiler.Profile) backend.ProductServicer {
	nps := ps
	nps.profiler = p
	return nps
}*/

// Get a mock product
func (ps ProductService) Get(sku string) backend.Producter {
	/*if ps.profiler != nil {
		ps.profiler.Start("API.mock", "Get Product "+sku)
		defer ps.profiler.End()
	}*/

	return Product{
		sku:         sku,
		name:        fmt.Sprintf("%s %s", nameprefix[rand.Intn(len(nameprefix))], namesuffix[rand.Intn(len(namesuffix))]),
		description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
}

// GetBySkuList gbsl
func (ps ProductService) GetBySkuList(skus []string) []backend.Producter {
	var products = make([]backend.Producter, len(skus))

	/*if ps.profiler != nil {
		ps.profiler.Start("API.mock", "Get Product List for "+strings.Join(skus, ", "))
		defer ps.profiler.End()
	}*/

	for i, sku := range skus {
		products[i] = ps.Get(sku)
	}

	return products
}

// Product p
type Product struct {
	sku         string
	name        string
	description string
}

// Sku s
func (p Product) Sku() string {
	return p.sku
}

// Name n
func (p Product) Name() string {
	return p.name
}

// Description d
func (p Product) Description() string {
	return p.description
}

// Price p
func (p Product) Price() float64 {
	return rand.Float64() * 100
}
