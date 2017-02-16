package interfaces

import "html/template"

// ProductService interface
type ProductService interface {
	Get(string) Product
	GetBySkuList([]string) []Product
}

// Product interface
type Product interface {
	Sku() string
	Name() string
	Description() string
	Price() float64
	Url() template.URL
}
