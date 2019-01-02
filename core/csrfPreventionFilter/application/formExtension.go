package application

import (
	"context"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	CrsfTokenFormExtension struct {
		service Service
	}
)

func (f *CrsfTokenFormExtension) Inject(service Service) {
	f.service = service
}

func (f *CrsfTokenFormExtension) Name() string {
	return "formExtension.csrfToken"
}

func (f *CrsfTokenFormExtension) Validate(_ context.Context, req *web.Request, _ domain.ValidatorProvider, _ interface{}) (*domain.ValidationInfo, error) {
	validationInfo := domain.ValidationInfo{}

	if !f.service.IsValid(req) {
		validationInfo.AddGeneralError("formError.crsfToken.invalid", "Invalid crsf token.")
	}

	return &validationInfo, nil
}
