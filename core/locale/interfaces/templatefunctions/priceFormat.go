package templatefunctions

import (
	"strconv"

	"flamingo.me/flamingo/core/locale/application"
	"flamingo.me/flamingo/core/pugtemplate/pugjs"
	"flamingo.me/flamingo/framework/config"
	"github.com/leekchan/accounting"
)

type (
	// PriceFormatFunc for formatting prices
	PriceFormatFunc struct {
		config             config.Map
		translationService application.TranslationServiceInterface
	}
)

func (pff *PriceFormatFunc) Inject(serviceInterface application.TranslationServiceInterface, config *struct {
	Config config.Map `inject:"config:locale.accounting"`
}) {
	pff.translationService = serviceInterface
	pff.config = config.Config
}

// Func as implementation of debug method
func (pff *PriceFormatFunc) Func() interface{} {
	return func(value interface{}, currency string) string {
		currency = pff.translationService.Translate(currency, currency, "", 1, nil)
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
		if valueNumber, ok := value.(pugjs.Number); ok {
			return ac.FormatMoney(float64(valueNumber))
		} else if valueString, ok := value.(pugjs.String); ok {
			float, err := strconv.ParseFloat(string(valueString), 64)
			if err != nil {
				float = 0.0
			}
			return ac.FormatMoney(float)
		} else {
			return ac.FormatMoney(0)
		}
	}
}
