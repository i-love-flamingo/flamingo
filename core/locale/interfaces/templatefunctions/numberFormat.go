package templatefunctions

import (
	"context"

	"github.com/leekchan/accounting"
)

// NumberFormatFunc for formatting numbers
type NumberFormatFunc struct {
	precision          float64
	decimal            string
	thousand           string
}

// Inject dependencies
func (nff *NumberFormatFunc) Inject(
	config *struct {
		Precision float64 `inject:"config:locale.numbers.precision"`
		Decimal   string  `inject:"config:locale.numbers.decimal"`
		Thousand  string  `inject:"config:locale.numbers.thousand"`
	},
) {
	nff.precision = config.Precision
	nff.decimal = config.Decimal
	nff.thousand = config.Thousand
}

// Func as implementation of debug method
func (nff *NumberFormatFunc) Func(context.Context) interface{} {
	return func(value interface{}, params ...int) string {

		precision := int(nff.precision)
		if len(params) > 0 {
			precision = params[0]
		}

		return accounting.FormatNumber(value, precision, nff.thousand, nff.decimal)
	}
}
