package domain

import (
	"context"

	"gopkg.in/go-playground/validator.v9"
)

type (
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
