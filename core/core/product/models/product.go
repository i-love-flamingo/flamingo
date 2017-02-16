package models

type (
	ProductAttribute struct {
		ID    string `structs:"id" json:"id"`
		Name  string `structs:"name" json:"name"`
		Value string `structs:"value" json:"value"`
	}

	Product struct {
		ID          string  `structs:"id" json:"id"`
		Name        string  `structs:"name" json:"name"`
		Description string  `structs:"description" json:"description"`
		Price       float64 `structs:"price" json:"price"`
	}
)
