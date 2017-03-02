package internalmock

import (
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/packages/internalmock/product"
)

// Register Services for internalmock package
func Register(c *di.Container) {
	c.Register(new(product.ProductService))
}
