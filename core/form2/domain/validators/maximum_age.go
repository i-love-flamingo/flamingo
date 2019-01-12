package validators

import (
	"context"
	"strconv"
	"strings"
	"time"

	"flamingo.me/flamingo/core/form2/domain"

	"gopkg.in/go-playground/validator.v9"
)

type (
	// MaximumAgeValidator defines maximum age validator which validates if passed date is after than desired years ago
	//
	// Data struct {
	//	 Date string `validate:"maximumage=150"`
	// }
	//
	MaximumAgeValidator struct {
		dateFormat string
	}
)

var _ domain.FieldValidator = &MaximumAgeValidator{}

// Inject is method used to set all dependencies as local variables
func (v *MaximumAgeValidator) Inject(cfg *struct {
	DateFormat string `inject:"config:form.validator.dateFormat"`
}) {
	v.dateFormat = cfg.DateFormat
}

// ValidatorName defines tag name of maximum age validator
func (v *MaximumAgeValidator) ValidatorName() string {
	return "maximumage"
}

// ValidateField validates string in date format for maximum age. Valid if string is empty or in wrong date format or in date range.
func (v *MaximumAgeValidator) ValidateField(_ context.Context, fl validator.FieldLevel) bool {
	param := fl.Param()
	years := 0

	if param != "" {
		value, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			panic(err.Error())
		}
		years = int(value)
	}

	converted, ok := fl.Field().Interface().(string)
	if !ok {
		return true
	}

	if len(strings.TrimSpace(converted)) == 0 {
		return true
	}

	date, err := time.Parse(v.dateFormat, converted)
	if err != nil {
		return true
	}

	now := time.Now()
	desired := time.Date(now.Year()-years, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date.After(desired) || date.Equal(desired)
}
