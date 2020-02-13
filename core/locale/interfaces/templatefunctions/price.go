package templatefunctions

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
)

// PriceFormatFunc for formatting prices
type PriceFormatFunc struct {
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

// PriceFormatLongFunc for formatting prices with additional currency code/label
type PriceFormatLongFunc struct {
	config       config.Map
	labelService *application.LabelService
	priceFormat  *PriceFormatFunc
	priceService *application.PriceService
}

// Inject dependencies
func (pff *PriceFormatLongFunc) Inject(
	labelService *application.LabelService,
	formatFunc *PriceFormatFunc,
	priceService *application.PriceService,
	config *struct {
		Config config.Map `inject:"config:core.locale.accounting"`
	},
) {
	pff.config = config.Config
	pff.labelService = labelService
	pff.priceFormat = formatFunc
	pff.priceService = priceService
}

// Func formats the value, adds currency sign/symbol and add an additional currency code/label
// example output could be: $ 21,500.99 USD
func (pff *PriceFormatLongFunc) Func(ctx context.Context) interface{} {
	return func(value float64, currency string, currencyLabel string) string {
		priceFunc := pff.priceFormat.Func(ctx).(func(value float64, currency string) string)
		price := priceFunc(value, currency)
		currencyLabel = pff.labelService.NewLabel(currencyLabel).String()

		// get config for currency or default config
		formatConfig := pff.priceService.GetConfigForCurrency(currency)

		if formatConfig.FormatLong != "" {
			return fmt.Sprintf(formatConfig.FormatLong, price, currencyLabel)
		}

		return price
	}
}
