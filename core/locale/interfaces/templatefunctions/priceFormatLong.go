package templatefunctions

import (
	"fmt"
	"flamingo.me/flamingo/core/locale/application"
	"flamingo.me/flamingo/framework/config"
)

type (
	// PriceFormatFunc for formatting prices
	PriceFormatLongFunc struct {
		Config             config.Map                              `inject:"config:locale.accounting"`
		TranslationService application.TranslationServiceInterface `inject:""`
		PriceFormat        *PriceFormatFunc                        `inject:""`
	}
)

// Name alias for use in template
func (pff PriceFormatLongFunc) Name() string {
	return "priceFormatLong"
}

// Func as implementation of debug method
func (pff PriceFormatLongFunc) Func() interface{} {
	return func(value interface{}, currency string, currencyLabel string) string {
		priceFunc := pff.PriceFormat.Func().(func(value interface{}, currency string) string)
		price := priceFunc(value, currency)
		currencyLabel = pff.TranslationService.Translate(currencyLabel, "", "", 1, nil)
		format, ok := pff.Config["formatLong"].(string)
		if ok {
			return fmt.Sprintf(format, price, currencyLabel)
		}
		return price
	}
}
