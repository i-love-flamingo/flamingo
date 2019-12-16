package application

import (
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
)

// PriceService for formatting prices
type PriceService struct {
	config       config.Map
	labelService *LabelService
}

// Inject dependencies
func (s *PriceService) Inject(labelService *LabelService, config *struct {
	Config config.Map `inject:"config:core.locale.accounting"`
}) {
	s.labelService = labelService
	s.config = config.Config
}

// GetConfigForCurrency get configuration for currency
func (s *PriceService) GetConfigForCurrency(currency string) config.Map {
	if configForCurrency, ok := s.config[currency]; ok {
		return configForCurrency.(config.Map)
	}

	if defaultConfig, ok := s.config["default"].(config.Map); ok {
		return defaultConfig
	}

	return s.config
}

// FormatPrice by price
func (s *PriceService) FormatPrice(value float64, currency string) string {
	currency = s.labelService.NewLabel(currency).String()

	configForCurrency := s.GetConfigForCurrency(currency)

	ac := accounting.Accounting{
		Symbol:    currency,
		Precision: 2,
	}
	decimal, ok := configForCurrency["decimal"].(string)
	if ok {
		ac.Decimal = decimal
	}
	thousand, ok := configForCurrency["thousand"].(string)
	if ok {
		ac.Thousand = thousand
	}
	formatZero, ok := configForCurrency["formatZero"].(string)
	if ok {
		ac.FormatZero = formatZero
	}
	format, ok := configForCurrency["format"].(string)
	if ok {
		ac.Format = format
	}

	return ac.FormatMoney(value)
}
