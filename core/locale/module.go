package locale

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure"
	"flamingo.me/flamingo/v3/core/locale/interfaces/controllers"
	"flamingo.me/flamingo/v3/core/locale/interfaces/templatefunctions"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module registers our profiler
type Module struct {
	EnableTranslationAPI bool `inject:"config:core.locale.enableTranslationApi,optional"`
}

type routes struct {
	translationController *controllers.TranslationController
}

// Configure the product URL
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.TranslationService)).In(dingo.ChildSingleton).To(infrastructure.TranslationService{})
	injector.Bind(new(application.DateTimeServiceInterface)).To(application.DateTimeService{})

	if m.EnableTranslationAPI {
		web.BindRoutes(injector, new(routes))
	}

	flamingo.BindTemplateFunc(injector, "__", new(templatefunctions.LabelFunc))
	flamingo.BindTemplateFunc(injector, "priceFormat", new(templatefunctions.PriceFormatFunc))
	flamingo.BindTemplateFunc(injector, "priceFormatLong", new(templatefunctions.PriceFormatLongFunc))
	flamingo.BindTemplateFunc(injector, "numberFormat", new(templatefunctions.NumberFormatFunc))
	flamingo.BindTemplateFunc(injector, "dateTimeFormatFromIso", new(templatefunctions.DateTimeFormatFromIso))
	flamingo.BindTemplateFunc(injector, "dateTimeFormat", new(templatefunctions.DateTimeFormatFromTime))
}

func (r *routes) Inject(
	tc *controllers.TranslationController,
) {
	r.translationController = tc
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("api.translations", r.translationController.GetAllTranslations)
	registry.MustRoute("/api/translations", "api.translations")
}

// CueConfig for this module
func (m *Module) CueConfig() string {
	return `
core: locale: {
	locale: string | *"en-US"
	translationFile: string | *""
	translationFiles: [...string] | *[]
	accounting: {
		default: {
			decimal: string | *"."
			thousand: string | *","
			formatZero: string | *"%s 0.00"
			format: string | *"%s %v"
			formatLong: string | *"%v %v"
		}
	}
	numbers: {
		decimal: string | *"."
		thousand: string | *","
		precision: float | int | *2
	}
	date: {
		dateFormat: string | *"02 Jan 2006"
		timeFormat: string | *"15:04:05"
		dateTimeFormat: string | *"02 Jan 2006 15:04:05"
		location: string | *"Europe/London"
	}
}
`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	alias := make(map[string]string)
	for _, v := range []string{
		"locale.locale",
		"locale.accounting.default.decimal",
		"locale.accounting.default.thousand",
		"locale.accounting.default.formatZero",
		"locale.accounting.default.format",
		"locale.accounting.default.formatLong",
		"locale.numbers.decimal",
		"locale.numbers.thousand",
		"locale.numbers.precision",
		"locale.date.dateFormat",
		"locale.date.timeFormat",
		"locale.date.dateTimeFormat",
		"locale.date.location",
	} {
		alias[v] = "core." + v
	}
	return alias
}
