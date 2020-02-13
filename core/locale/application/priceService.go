package application

import (
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/leekchan/accounting"
)

// PriceService for formatting prices
type PriceService struct {
	configs      map[string]PriceFormatConfig
	labelService *LabelService
}

// PriceFormatConfig represents price formatting configuration which is possible to specify.
type PriceFormatConfig struct {
	Decimal    string `json:"decimal"`
	Thousand   string `json:"thousand"`
	FormatZero string `json:"formatZero"`
	FormatLong string `json:"formatLong"`
	Format     string `json:"format"`
}

// Inject dependencies
func (s *PriceService) Inject(labelService *LabelService, logger flamingo.Logger, config *struct {
	Config config.Map `inject:"config:core.locale.accounting"`
}) {
	s.labelService = labelService

	if config == nil {
		return
	}

	err := config.Config.MapInto(&s.configs)
	if err != nil {
		logger.WithField("category", "PriceService").Error(err)
	}
}

// GetConfigForCurrency get configuration for currency
func (s *PriceService) GetConfigForCurrency(currency string) PriceFormatConfig {
	if configForCurrency, ok := s.configs[currency]; ok {
		return configForCurrency
	}

	return s.configs["default"]
}

// FormatPrice by price
func (s *PriceService) FormatPrice(value float64, currency string) string {
	currency = s.labelService.NewLabel(currency).String()

	configForCurrency := s.GetConfigForCurrency(currency)

	ac := accounting.Accounting{
		Symbol:    currency,
		Precision: 2,
	}

	if configForCurrency.Decimal != "" {
		ac.Decimal = configForCurrency.Decimal
	}

	if configForCurrency.Thousand != "" {
		ac.Thousand = configForCurrency.Thousand
	}

	if configForCurrency.FormatZero != "" {
		ac.FormatZero = configForCurrency.FormatZero
	}

	if configForCurrency.Format != "" {
		ac.Format = configForCurrency.Format
	}

	return ac.FormatMoney(value)
}
