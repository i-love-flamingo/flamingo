package templatefunctions

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
)

// PriceFormatLongFunc for formatting prices with additional currency code/label
type PriceFormatLongFunc struct {
	config       config.Map
	labelService *application.LabelService
	priceFormat  *PriceFormatFunc
}

// Inject dependencies
func (pff *PriceFormatLongFunc) Inject(
	labelService *application.LabelService,
	formatFunc *PriceFormatFunc,
	config *struct {
		Config config.Map `inject:"config:locale.accounting"`
	},
) {
	pff.config = config.Config
	pff.labelService = labelService
	pff.priceFormat = formatFunc
}

// Func formats the value, adds currency sign/symbol and add an additional currency code/label
// example output could be: $ 21,500.99 USD
func (pff *PriceFormatLongFunc) Func(ctx context.Context) interface{} {
	return func(value float64, currency string, currencyLabel string) string {
		priceFunc := pff.priceFormat.Func(ctx).(func(value float64, currency string) string)
		price := priceFunc(value, currency)
		currencyLabel = pff.labelService.NewLabel(currencyLabel).String()
		format, ok := pff.config["formatLong"].(string)
		if ok {
			return fmt.Sprintf(format, price, currencyLabel)
		}
		return price
	}
}
