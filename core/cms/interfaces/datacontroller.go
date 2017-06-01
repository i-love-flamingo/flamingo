package interfaces

import (
	"flamingo/framework/web"
)

type (
	DataController struct{}
)

// Get Response for Product matching sku param
func (vc *DataController) Data(c web.Context) interface{} {
	return "Hello World Block."
}
