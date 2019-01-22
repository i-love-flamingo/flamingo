package templatefunctions

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
)

// PriceFormatLongFunc for formatting prices
type PriceFormatLongFunc struct {
	config             config.Map
	translationService application.TranslationServiceInterface
	priceFormat        *PriceFormatFunc
}

// Inject dependencies
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
func (pff *PriceFormatLongFunc) Func(ctx context.Context) interface{} {
	return func(value interface{}, currency string, currencyLabel string) string {
		priceFunc := pff.priceFormat.Func(ctx).(func(value interface{}, currency string) string)
		price := priceFunc(value, currency)
		currencyLabel = pff.translationService.Translate(currencyLabel, "", "", 1, nil)
		format, ok := pff.config["formatLong"].(string)
		if ok {
			return fmt.Sprintf(format, price, currencyLabel)
		}
		return price
	}
}
