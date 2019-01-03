package formData

import (
	"context"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	DefaultFormDataValidatorImpl struct{}
)

var _ domain.DefaultFormDataValidator = &DefaultFormDataValidatorImpl{}

func (p *DefaultFormDataValidatorImpl) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	validationInfo := validatorProvider.Validate(ctx, req, formData)
	return &validationInfo, nil
}
