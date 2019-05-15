package templatefunctions

import (
	"context"
	"math/big"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
)

// PriceFormatFunc for formatting prices
type PriceFormatFunc struct {
	config       config.Map
	labelService *application.LabelService
}

// Inject dependencies
func (pff *PriceFormatFunc) Inject(labelService *application.LabelService, config *struct {
	Config config.Map `inject:"config:locale.accounting"`
}) {
	pff.labelService = labelService
	pff.config = config.Config
}

// Func formats the value and adds currency sign/symbol
// example output could be: $ 21,500.99
// (supported value types : int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, *big.Rat, *big.Float)
func (pff *PriceFormatFunc) Func(context.Context) interface{} {
	return func(value interface{}, currency string) string {
		currency = pff.labelService.NewLabel(currency).String()
		ac := accounting.Accounting{
			Symbol:    currency,
			Precision: 2,
		}
		decimal, ok := pff.config["decimal"].(string)
		if ok {
			ac.Decimal = decimal
		}
		thousand, ok := pff.config["thousand"].(string)
		if ok {
			ac.Thousand = thousand
		}
		formatZero, ok := pff.config["formatZero"].(string)
		if ok {
			ac.FormatZero = formatZero
		}
		format, ok := pff.config["format"].(string)
		if ok {
			ac.Format = format
		}

		valueBigFloat, ok := value.(*big.Float)
		if ok {
			return ac.FormatMoneyBigFloat(valueBigFloat)
		}

		return ac.FormatMoney(value)
	}
}
