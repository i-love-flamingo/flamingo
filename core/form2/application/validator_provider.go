package application

import (
	"context"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// ValidatorProviderImpl as struct which implements interface ValidatorProvider
	ValidatorProviderImpl struct {
		validate *validator.Validate
	}
)

var _ domain.ValidatorProvider = &ValidatorProviderImpl{}

// Inject initialize instance of validator.Validate struct
func (p *ValidatorProviderImpl) Inject(fieldValidators []domain.FieldValidator, structValidators []domain.StructValidator) {
	validate := validator.New()
	p.attachFieldValidators(validate, fieldValidators)
	p.attachStructValidators(validate, structValidators)
	p.validate = validate
}

// Validate method which validates any struct and returns domain.ValidationInfo as a result of validation
func (p *ValidatorProviderImpl) Validate(ctx context.Context, req *web.Request, value interface{}) domain.ValidationInfo {
	reqCtx := web.Context_(ctx, req)
	validate := p.GetValidator()
	err := validate.StructCtx(reqCtx, value)

	return p.ErrorsToValidationInfo(err)
}

// GetValidator method which returns instance of validator.Validate struct with all injected field and struct validations
func (p *ValidatorProviderImpl) GetValidator() *validator.Validate {
	return p.validate
}

// ErrorsToValidationInfo method which transforms errors into domain.ValidationInfo
func (p *ValidatorProviderImpl) ErrorsToValidationInfo(err error) domain.ValidationInfo {
	validationInfo := domain.ValidationInfo{}

	if err == nil {
		return validationInfo
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			fieldName := p.getRelativeFieldNameFromValidationError(err)
			validationInfo.AddFieldError(fieldName, "formError."+fieldName+"."+err.Tag(), err.Field()+" "+err.Tag())
		}
	} else {
		validationInfo.AddGeneralError("formError.invalidValidation", err.Error())
	}

	return validationInfo
}

// attachFieldValidators method which attach all injected instances of FieldValidator interface into validator.Validate instance
func (p *ValidatorProviderImpl) attachFieldValidators(validate *validator.Validate, fieldValidators []domain.FieldValidator) {
	for _, fieldValidator := range fieldValidators {
		validate.RegisterValidationCtx(fieldValidator.ValidatorName(), fieldValidator.ValidateField)
	}
}

// attachFieldValidators method which attach all injected instances of StructValidator interface into validator.Validate instance
func (p *ValidatorProviderImpl) attachStructValidators(validate *validator.Validate, structValidators []domain.StructValidator) {
	for _, structValidator := range structValidators {
		validate.RegisterStructValidationCtx(structValidator.ValidateStruct, structValidator.StructType())
	}
}

// getRelativeFieldNameFromValidationError method which extracts relative field name depending on it's full namespace
func (p *ValidatorProviderImpl) getRelativeFieldNameFromValidationError(err validator.FieldError) string {
	var result []string

	namespace := err.Namespace()
	//first part of namespace is not required to have the relative path:
	fieldName := namespace[(strings.Index(namespace, ".") + 1):]
	for _, part := range strings.Split(fieldName, ".") {
		result = append(result, strings.ToLower(part[0:1])+part[1:])
	}

	return strings.Join(result, ".")
}
