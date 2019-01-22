package locale

import (
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/template"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*application.TranslationServiceInterface)(nil)).To(application.TranslationService{})
	injector.Bind((*application.DateTimeServiceInterface)(nil)).To(application.DateTimeService{})

	template.BindFunc(injector, "__", new(templatefunctions.Label))
	template.BindFunc(injector, "priceFormat", new(templatefunctions.PriceFormatFunc))
	template.BindFunc(injector, "priceFormatLong", new(templatefunctions.PriceFormatLongFunc))
	template.BindFunc(injector, "numberFormat", new(templatefunctions.NumberFormatFunc))
	template.BindFunc(injector, "dateTimeFormatFromIso", new(templatefunctions.DateTimeFormatFromIso))
	template.BindFunc(injector, "dateTimeFormat", new(templatefunctions.DateTimeFormatFromTime))
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
