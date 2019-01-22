package templatefunctions

import (
	"context"

	"flamingo.me/flamingo/v3/core/locale/application"
	"github.com/leekchan/accounting"
)

// NumberFormatFunc for formatting numbers
type NumberFormatFunc struct {
	precision          float64
	decimal            string
	thousand           string
	translationService application.TranslationServiceInterface
}

// Inject dependencies
func (nff *NumberFormatFunc) Inject(
	serviceInterface application.TranslationServiceInterface,
	config *struct {
		Precision float64 `inject:"config:locale.numbers.precision"`
		Decimal   string  `inject:"config:locale.numbers.decimal"`
		Thousand  string  `inject:"config:locale.numbers.thousand"`
	},
) {
	nff.translationService = serviceInterface
	nff.precision = config.Precision
	nff.decimal = config.Decimal
	nff.thousand = config.Thousand
}

// Func as implementation of debug method
func (nff *NumberFormatFunc) Func(context.Context) interface{} {
	return func(value interface{}, params ...interface{}) string {

		precision := int(nff.precision)
		if len(params) > 0 {
			if precisionIntParam, ok := params[0].(int); ok {
				precision = precisionIntParam
			}
			// todo fix
			//if precisionNumberParam, ok := params[0].(pugjs.Number); ok {
			//	precision = int(precisionNumberParam)
			//}
		}

		return accounting.FormatNumber(value, precision, nff.thousand, nff.decimal)
	}
}
