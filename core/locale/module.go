package locale

import (
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/locale/interfaces/templatefunctions"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/template"
)

type (
	// Module registers our profiler
	Module struct{}
)

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*application.TranslationServiceInterface)(nil)).To(application.TranslationService{})

	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.Label{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.PriceFormatFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.PriceFormatLongFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.NumberFormatFunc{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.DateTimeFormatFromIso{})
	injector.BindMulti((*template.Function)(nil)).To(templatefunctions.DateTimeFormatFromTime{})
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
