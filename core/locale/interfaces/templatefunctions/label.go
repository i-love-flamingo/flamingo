package templatefunctions

import (
	"context"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo/v3/core/locale/application"
)

// Label is exported as a template function
type (
	Label struct {
		labelService *application.LabelService
		logger flamingo.Logger
	}
)

// Inject dependencies
func (tf *Label) Inject(labelService *application.LabelService,logger flamingo.Logger) {
	tf.labelService = labelService
	tf.logger = logger.WithField("module","locale").WithField("category","templatefunctions.label")
}

// Func template function factory
// todo fix
func (tf *Label) Func(context.Context) interface{} {

	// Usage:  __("key")
	// __("key","default")
	// __("key","Hello Mr {{.userName}}",{UserName: "Max"})
	// Force other than configured locale: __("switch_to_german","",{},"de-DE")
	return func(key string, params ...interface{}) *domain.Label {

		if len(params) > 0 {
			tf.logger.Warn("Depricated unsupported paramaters given! Use the Setters provided by the returned Label")

		}
		return tf.labelService.NewLabel(key)
	}
}
