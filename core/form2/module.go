package form

import (
	"flamingo.me/flamingo/core/form2/application/provider"
	"flamingo.me/flamingo/core/form2/application/provider/validators"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
)

type (
	Module struct {
		CustomRegex config.Map `inject:"config:form.validator.customRegex"`
	}
)

func (m *Module) Configure(injector *dingo.Injector) {
	for name, value := range m.CustomRegex {
		regex, ok := value.(string)
		if !ok {
			panic("wrong value passed as validation regex")
		}
		regexValidator := validators.NewRegexValidator(name, regex)
		injector.BindMulti(new(provider.FieldValidator)).ToInstance(regexValidator)
	}
	injector.BindMulti(new(provider.FieldValidator)).To(validators.DateFormatValidator{})
	injector.BindMulti(new(provider.FieldValidator)).To(validators.MinimumAgeValidator{})
	injector.BindMulti(new(provider.FieldValidator)).To(validators.MaximumAgeValidator{})

	injector.Bind(new(provider.ValidatorProvider)).To(provider.ValidatorProviderImpl{})
}

// DefaultConfig method which is responsible for setting up default module configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"form.validator": config.Map{
			"dateFormat":  "2006-01-02",
			"customRegex": config.Map{},
		},
	}
}
