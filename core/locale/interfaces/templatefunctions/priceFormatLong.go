package templatefunctions

import (
	"fmt"

	"flamingo.me/flamingo/core/locale/application"
	"flamingo.me/flamingo/framework/config"
)

type (
	// PriceFormatFunc for formatting prices
	PriceFormatLongFunc struct {
		config             config.Map
		translationService application.TranslationServiceInterface
		priceFormat        *PriceFormatFunc
	}
)

func (pff *PriceFormatLongFunc) Inject(
	serviceInterface application.TranslationServiceInterface,
	formatFunc *PriceFormatFunc,
	config *struct {
		Config config.Map `inject:"config:locale.accounting"`
	},
) {
	pff.config = config.Config
	pff.translationService = serviceInterface
	pff.priceFormat = formatFunc
}

// Func as implementation of debug method
func (pff *PriceFormatLongFunc) Func() interface{} {
	return func(value interface{}, currency string, currencyLabel string) string {
		priceFunc := pff.priceFormat.Func().(func(value interface{}, currency string) string)
		price := priceFunc(value, currency)
		currencyLabel = pff.translationService.Translate(currencyLabel, "", "", 1, nil)
		format, ok := pff.config["formatLong"].(string)
		if ok {
			return fmt.Sprintf(format, price, currencyLabel)
		}
		return price
	}
}
