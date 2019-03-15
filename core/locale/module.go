package locale

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.TranslationService)).To(infrastructure.TranslationService{})
	injector.Bind(new(application.DateTimeServiceInterface)).To(application.DateTimeService{})

	flamingo.BindTemplateFunc(injector, "__", new(templatefunctions.Label))
	flamingo.BindTemplateFunc(injector, "priceFormat", new(templatefunctions.PriceFormatFunc))
	flamingo.BindTemplateFunc(injector, "priceFormatLong", new(templatefunctions.PriceFormatLongFunc))
	flamingo.BindTemplateFunc(injector, "numberFormat", new(templatefunctions.NumberFormatFunc))
	flamingo.BindTemplateFunc(injector, "dateTimeFormatFromIso", new(templatefunctions.DateTimeFormatFromIso))
	flamingo.BindTemplateFunc(injector, "dateTimeFormat", new(templatefunctions.DateTimeFormatFromTime))
}

// DefaultConfig for this module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"locale": config.Map{
			"locale": "en-US",
			"accounting": config.Map{
				"decimal":    ",",
				"thousand":   ".",
				"formatZero": "%s -,-",
				"format":     "%s %v",
				"formatLong": "%v %v",
			},
			"numbers": config.Map{
				"decimal":   ",",
				"thousand":  ".",
				"precision": float64(2),
			},
			"date": config.Map{
				"dateFormat":     "02 Jan 2006",
				"timeFormat":     "15:04:05",
				"dateTimeFormat": "02 Jan 2006 15:04:05",
				"location":       "Europe/London",
			},
		},
	}
}
