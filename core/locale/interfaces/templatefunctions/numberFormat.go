package templatefunctions

import (
	"flamingo.me/flamingo/core/locale/application"
	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"github.com/leekchan/accounting"
)

type (
	// NumberFormatFunc for formatting numbers
	NumberFormatFunc struct {
		precision          float64
		decimal            string
		thousand           string
		translationService application.TranslationServiceInterface
	}
)

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

// Name alias for use in template
func (nff NumberFormatFunc) Name() string {
	return "numberFormat"
}

// Func as implementation of debug method
func (nff NumberFormatFunc) Func() interface{} {
	return func(value interface{}, params ...interface{}) string {

		precision := int(nff.precision)
		if len(params) > 0 {
			if precisionIntParam, ok := params[0].(int); ok {
				precision = precisionIntParam
			}
			if precisionNumberParam, ok := params[0].(pugjs.Number); ok {
				precision = int(precisionNumberParam)
			}
		}

		return accounting.FormatNumber(value, precision, nff.thousand, nff.decimal)
	}
}
