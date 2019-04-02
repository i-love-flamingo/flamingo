package templatefunctions

import (
	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/config"
)

type (
	PriceFormatFunc struct {
		Config config.Map `inject:"config:accounting"`
	}
)

// Name alias for use in template
func (df PriceFormatFunc) Name() string {
	return "priceFormat"
}

// Func as implementation of debug method
func (f PriceFormatFunc) Func() interface{} {
	return func(value pugjs.Number, currency string) string {
		ac := accounting.Accounting{
			Symbol:    currency,
			Precision: 2,
		}
		decimal, ok := f.Config["decimal"].(string)
		if ok {
			ac.Decimal = decimal
		}
		thousand, ok := f.Config["thousand"].(string)
		if ok {
			ac.Thousand = thousand
		}
		formatZero, ok := f.Config["formatZero"].(string)
		if ok {
			ac.FormatZero = formatZero
		}
		format, ok := f.Config["format"].(string)
		if ok {
			ac.Format = format
		}
		return ac.FormatMoney(float64(value))
	}
}
