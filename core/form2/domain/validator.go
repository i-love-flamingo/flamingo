package domain

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"gopkg.in/go-playground/validator.v9"
)

type (
	// ValidatorProvider as interface for defining main validator provider
	ValidatorProvider interface {
		// Validate method which validates any struct and returns domain.ValidationInfo as a result of validation
		Validate(ctx context.Context, req *web.Request, value interface{}) ValidationInfo
		// GetValidator method which returns instance of validator.Validate struct with all injected field and struct validations
		GetValidator() *validator.Validate
		// ErrorsToValidationInfo method which transforms errors into domain.ValidationInfo
		ErrorsToValidationInfo(err error) ValidationInfo
	}

	// FieldValidator as interface for defining custom field validation
	FieldValidator interface {
		// ValidatorName defines validator name used in fields' tags inside structs
		ValidatorName() string
		// ValidateField defines validation method called when field is validated
		ValidateField(ctx context.Context, fl validator.FieldLevel) bool
	}

	// StructValidator as interface for defining custom struct validation
	StructValidator interface {
		// StructType defines struct type which should be validated
		StructType() interface{}
		// ValidateStruct defines validation method called when struct is validated
		ValidateStruct(ctx context.Context, sl validator.StructLevel)
	}
)
