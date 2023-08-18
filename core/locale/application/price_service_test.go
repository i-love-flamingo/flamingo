package application

import (
	"testing"

	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
)

func TestPriceService_GetConfigForCurrency(t *testing.T) {
	// get empty
	assert.Equal(t, PriceFormatConfig{}, new(PriceService).GetConfigForCurrency("currency"))

	// get default
	service := &PriceService{}
	service.Inject(nil, nil, &struct {
		Config config.Map `inject:"config:core.locale.accounting"`
	}{
		Config: config.Map{
			"default": config.Map{
				"decimal":    ".",
				"thousand":   ",",
				"formatZero": "%s 0.00",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
		},
	})

	assert.Equal(t, PriceFormatConfig{
		Decimal:    ".",
		Thousand:   ",",
		FormatZero: "%s 0.00",
		Format:     "%s %v",
		FormatLong: "%v %v",
	}, service.GetConfigForCurrency("currency"))

	// get default and specific
	service = &PriceService{}
	service.Inject(nil, nil, &struct {
		Config config.Map `inject:"config:core.locale.accounting"`
	}{
		Config: config.Map{
			"default": config.Map{
				"decimal":    ".",
				"thousand":   ",",
				"formatZero": "%s 0.00",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
			"currency1": config.Map{
				"decimal":    ". 1",
				"thousand":   ", 1",
				"formatZero": "%s 0.00 1",
				"format":     "%s %v 1",
				"formatLong": "%v %v 1",
			},
		},
	})

	assert.Equal(t, PriceFormatConfig{
		Decimal:    ". 1",
		Thousand:   ", 1",
		FormatZero: "%s 0.00 1",
		Format:     "%s %v 1",
		FormatLong: "%v %v 1",
	}, service.GetConfigForCurrency("currency1"))
	assert.Equal(t, PriceFormatConfig{
		Decimal:    ".",
		Thousand:   ",",
		FormatZero: "%s 0.00",
		Format:     "%s %v",
		FormatLong: "%v %v",
	}, service.GetConfigForCurrency("currency2"))
}

func TestPriceService_FormatPrice(t *testing.T) {
	labelService := &LabelService{}
	labelService.Inject(func() *domain.Label {
		label := &domain.Label{}
		label.Inject(&fake.TranslationService{})
		return label
	}, nil, nil, nil)

	// get empty
	service := &PriceService{}
	service.Inject(labelService, nil, nil)
	assert.Equal(t, "currency10,000.50", service.FormatPrice(10000.5, "currency"))
	assert.Equal(t, "currency0.00", service.FormatPrice(0, "currency"))

	// get default
	service = &PriceService{}
	service.Inject(labelService, nil, &struct {
		Config config.Map `inject:"config:core.locale.accounting"`
	}{
		Config: config.Map{
			"default": config.Map{
				"decimal":    ".",
				"thousand":   ",",
				"formatZero": "%s 0.00",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
		},
	})
	assert.Equal(t, "currency 10,000.50", service.FormatPrice(10000.5, "currency"))
	assert.Equal(t, "currency 0.00", service.FormatPrice(0, "currency"))

	// get default and specific
	service = &PriceService{}
	service.Inject(labelService, nil, &struct {
		Config config.Map `inject:"config:core.locale.accounting"`
	}{
		Config: config.Map{
			"default": config.Map{
				"decimal":    ".",
				"thousand":   ",",
				"formatZero": "0.00 %s",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
			"currency1": config.Map{
				"decimal":    ",",
				"thousand":   ".",
				"formatZero": "%s 0,0",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
		},
	})

	assert.Equal(t, "currency1 10.000,50", service.FormatPrice(10000.5, "currency1"))
	assert.Equal(t, "currency1 0,0", service.FormatPrice(0, "currency1"))

	// unknown currency should take default:
	assert.Equal(t, "currency2 10,000.50", service.FormatPrice(10000.5, "currency2"))
	assert.Equal(t, "0.00 currency2", service.FormatPrice(0, "currency2"))
}
