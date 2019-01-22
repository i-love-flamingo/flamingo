package form

import (
	"flamingo.me/flamingo/v3/core/form/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
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
		},
	}
}
