package product

import (
	"html/template"
	"math/rand"
)

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

func (p Product) Url() template.URL {
	return "/test"
}
