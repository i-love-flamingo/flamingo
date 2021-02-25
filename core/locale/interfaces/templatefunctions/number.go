package templatefunctions

import (
	"context"
	"math/big"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"github.com/leekchan/accounting"
)

// NumberFormatFunc for formatting numbers
type NumberFormatFunc struct {
	precision float64
	decimal   string
	thousand  string
	logger    flamingo.Logger
}

// Inject dependencies
func (nff *NumberFormatFunc) Inject(
	logger flamingo.Logger,
	config *struct {
		Precision float64 `inject:"config:core.locale.numbers.precision"`
		Decimal   string  `inject:"config:core.locale.numbers.decimal"`
		Thousand  string  `inject:"config:core.locale.numbers.thousand"`
	},
) {
	nff.precision = config.Precision
	nff.decimal = config.Decimal
	nff.thousand = config.Thousand
	nff.logger = logger
}

// Func returns the template function for formatting numbers
func (nff *NumberFormatFunc) Func(context.Context) interface{} {
	return func(value interface{}, params ...int) string {

		precision := int(nff.precision)
		if len(params) > 0 {
			precision = params[0]
		}

		defer func() {
			if err := recover(); err != nil {
				nff.logger.Error(err)
			}
		}()

		valueBigFloat, ok := value.(*big.Float)
		if ok {
			return accounting.FormatNumberBigFloat(valueBigFloat, precision, nff.thousand, nff.decimal)
		}

		return accounting.FormatNumber(value, precision, nff.thousand, nff.decimal)
	}
}
