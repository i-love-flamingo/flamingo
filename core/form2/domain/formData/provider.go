package formdata

import (
	"context"

	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// DefaultFormDataProviderImpl represents implementation of default domain.FormDataProvider.
	DefaultFormDataProviderImpl struct {}
)

var _ domain.DefaultFormDataProvider = &DefaultFormDataProviderImpl{}

// GetFormData performs default form data providing, by passing simple form data as instance of map[string]string.
func (p *DefaultFormDataProviderImpl) GetFormData(context.Context, *web.Request) (interface{}, error) {
	return map[string]string{}, nil
}
