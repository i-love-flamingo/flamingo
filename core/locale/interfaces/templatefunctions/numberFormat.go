package templatefunctions

import (
	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
)

type (
	// NumberFormatFunc for formatting numbers
	NumberFormatFunc struct {
		Precision          float64                                 `inject:"config:locale.numbers.precision"`
		Decimal            string                                  `inject:"config:locale.numbers.decimal"`
		Thousand           string                                  `inject:"config:locale.numbers.thousand"`
		TranslationService application.TranslationServiceInterface `inject:""`
	}
)

// Name alias for use in template
func (nff NumberFormatFunc) Name() string {
	return "numberFormat"
}

// Func as implementation of debug method
func (nff NumberFormatFunc) Func() interface{} {
	return func(value interface{}, params ...interface{}) string {

		precision := int(nff.Precision)
		if len(params) > 0 {
			if precisionIntParam, ok := params[0].(int); ok {
				precision = precisionIntParam
			}
			if precisionNumberParam, ok := params[0].(pugjs.Number); ok {
				precision = int(precisionNumberParam)
			}
		}

		return accounting.FormatNumber(value, precision, nff.Thousand, nff.Decimal)
	}
}
