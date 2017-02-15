package internalmock

import (
	"flamingo/core/contrib/internalmock/product"
	"flamingo/core/core/app"
)

func Register(r *app.Registrator) {
	r.Object(new(product.ProductService))
}
