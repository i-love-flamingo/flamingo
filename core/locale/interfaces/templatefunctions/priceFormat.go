package templatefunctions

import (
	"context"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
)

// PriceFormatFunc for formatting prices
type PriceFormatFunc struct {
	config       config.Map
	priceService *application.PriceService
}

// Inject dependencies
func (pff *PriceFormatFunc) Inject(priceService *application.PriceService) {
	pff.priceService = priceService
}

// Func formats the value and adds currency sign/symbol
// example output could be: $ 21,500.99
func (pff *PriceFormatFunc) Func(context.Context) interface{} {
	return func(value float64, currency string) string {
		return pff.priceService.FormatPrice(value, currency)
	}
}
