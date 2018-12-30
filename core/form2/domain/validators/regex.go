package validators

import (
	"context"
	"flamingo.me/flamingo/core/form2/domain"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

type (
	/* RegexValidator defines regex validator. Name and regex pattern is passed to fresh instance during creation of validator.

	validator := NewRegexValidator("postcode", "^[0-9]{5}$")

	...

	Data struct {
		PostCode string `validate:"postcode"`
	}

	 */
	RegexValidator struct {
		name  string
		regex *regexp.Regexp
	}
)

var _ domain.FieldValidator = &RegexValidator{}

// NewRegexValidator creates new instance of RegexValidator by defining it's tag name and regex pattern
func NewRegexValidator(name string, regex string) *RegexValidator {
	return &RegexValidator{
		name:  name,
		regex: regexp.MustCompile(regex),
	}
}

// ValidatorName defines tag name of regex validator
func (v *RegexValidator) ValidatorName() string {
	return v.name
}

// ValidateField validates string if match right regex. Valid if string is empty or match defined regex pattern.
func (v *RegexValidator) ValidateField(_ context.Context, fl validator.FieldLevel) bool {
	converted, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	if len(converted) == 0 {
		return true
	}

	return v.regex.MatchString(converted)
}
