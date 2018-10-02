package form

import (
	"flamingo.me/flamingo/core/form/domain"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"gopkg.in/go-playground/validator.v9"
)

type (
	Module struct{}
)

func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(validator.Validate)).ToProvider(domain.ValidatorProvider)
}

// DefaultConfig method which is responsible for setting up default module configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"form.validator": config.Map{
			"dateFormat": "2006-01-02",
			"minimumAge": 18.0,
			"maximumAge": 150.0,
		},
	}
}
