package application

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
	"github.com/stretchr/testify/assert"
	"testing"
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
	acc := &accounting.Accounting{
		Symbol:    "currency",
		Precision: 2,
	}
	assert.Equal(t, acc.FormatMoney(10000.5), service.FormatPrice(10000.5, "currency"))

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

	acc = &accounting.Accounting{
		Symbol:     "currency",
		Precision:  2,
		Decimal:    ".",
		Thousand:   ",",
		FormatZero: "%s 0.00",
		Format:     "%s %v",
	}
	assert.Equal(t, acc.FormatMoney(10000.5), service.FormatPrice(10000.5, "currency"))

	// get default and specific
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
			"currency1": config.Map{
				"decimal":    ",",
				"thousand":   ".",
				"formatZero": "%s 0.00",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
		},
	})

	acc = &accounting.Accounting{
		Symbol:     "currency1",
		Precision:  2,
		Decimal:    ",",
		Thousand:   ".",
		FormatZero: "%s 0.00",
		Format:     "%s %v",
	}
	assert.Equal(t, acc.FormatMoney(10000.5), service.FormatPrice(10000.5, "currency1"))
	acc = &accounting.Accounting{
		Symbol:     "currency2",
		Precision:  2,
		Decimal:    ".",
		Thousand:   ",",
		FormatZero: "%s 0.00",
		Format:     "%s %v",
	}
	assert.Equal(t, acc.FormatMoney(10000.5), service.FormatPrice(10000.5, "currency2"))
}
