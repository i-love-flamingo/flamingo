package validators

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type (
	/* MinimumAgeValidator defines minimum age validator which validates if passed date is before than desired years ago

	Data struct {
		Date string `validate:"minimumage=18"`
	}

	 */
	MinimumAgeValidator struct {
		dateFormat string
	}
)

func (v *MinimumAgeValidator) Inject(cfg *struct {
	DateFormat string `inject:"config:form.validator.dateFormat"`
}) {
	v.dateFormat = cfg.DateFormat
}

// ValidatorName defines tag name of minimum age validator
func (v *MinimumAgeValidator) ValidatorName() string {
	return "minimumage"
}

// ValidateField validates string in date format for minimum age. Valid if string is empty or in wrong date format or in date range.
func (v *MinimumAgeValidator) ValidateField(_ context.Context, fl validator.FieldLevel) bool {
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
	desired := time.Date(now.Year()-years, now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

	return date.Before(desired)
}
