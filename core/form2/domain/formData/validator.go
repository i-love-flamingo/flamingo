package formdata

import (
	"context"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// DefaultFormDataValidatorImpl represents implementation of default domain.FormDataValidator.
	DefaultFormDataValidatorImpl struct{}
)

var _ domain.DefaultFormDataValidator = &DefaultFormDataValidatorImpl{}

// Validate performs default form data validation, by using go-playground validator package and storing results into domain.ValidationInfo instance.
func (p *DefaultFormDataValidatorImpl) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	validationInfo := validatorProvider.Validate(ctx, req, formData)
	return &validationInfo, nil
}
