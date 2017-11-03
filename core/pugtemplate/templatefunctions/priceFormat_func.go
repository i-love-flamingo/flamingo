package templatefunctions

import (
	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/pugtemplate/pugjs"
	"go.aoe.com/flamingo/framework/config"
)

type (
	// PriceFormatFunc for formatting prices
	PriceFormatFunc struct {
		Config config.Map `inject:"config:accounting"`
	}
)

// Name alias for use in template
func (pff PriceFormatFunc) Name() string {
	return "priceFormat"
}

// Func as implementation of debug method
func (pff PriceFormatFunc) Func() interface{} {
	return func(value pugjs.Number, currency string) string {
		ac := accounting.Accounting{
			Symbol:    currency,
			Precision: 2,
		}
		decimal, ok := pff.Config["decimal"].(string)
		if ok {
			ac.Decimal = decimal
		}
		thousand, ok := pff.Config["thousand"].(string)
		if ok {
			ac.Thousand = thousand
		}
		formatZero, ok := pff.Config["formatZero"].(string)
		if ok {
			ac.FormatZero = formatZero
		}
		format, ok := pff.Config["format"].(string)
		if ok {
			ac.Format = format
		}
		return ac.FormatMoney(float64(value))
	}
}
