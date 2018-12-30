package formData

import (
	"context"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	DefaultFormDataProviderImpl struct {}
)

var _ domain.DefaultFormDataProvider = &DefaultFormDataProviderImpl{}

func (p *DefaultFormDataProviderImpl) GetFormData(context.Context, *web.Request) (interface{}, error) {
	return map[string]string{}, nil
}
