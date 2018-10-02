package domain

import (
	"regexp"
	"time"

	"go.aoe.com/flamingo/framework/config"
	"gopkg.in/go-playground/validator.v9"
)

func toDate(value string, dateFormat string) *time.Time {
	if value == "" {
		return nil
	}

	date, err := time.Parse(dateFormat, value)
	if err != nil {
		return nil
	}

	return &date
}

func yearsFromNow(years int) time.Time {
	date := time.Now()

	return time.Date(
		date.Year()+years,
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location(),
	)
}

func validateDateFormat(value string, dateFormat string) bool {
	_, err := time.Parse(dateFormat, value)
	return err == nil
}

func validateMinimumAge(value string, dateFormat string, minimumAge int) bool {
	date := toDate(value, dateFormat)
	if date == nil {
		return true
	}

	required := yearsFromNow(-minimumAge)

	return date.Add(-24 * time.Hour).Before(required)
}

func validateMaximumAge(value string, dateFormat string, maximumAge int) bool {
	date := toDate(value, dateFormat)
	if date == nil {
		return true
	}

	required := yearsFromNow(-maximumAge)

	return date.After(required)
}

func validateRegex(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

func dateFormatValidatorProvider(dateFormat string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateDateFormat(fl.Field().String(), dateFormat)
	}
}

func maximumAgeValidatorProvider(dateFormat string, maximumAge int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateMaximumAge(fl.Field().String(), dateFormat, maximumAge)
	}
}

func minimumAgeValidatorProvider(dateFormat string, minimumAge int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateMinimumAge(fl.Field().String(), dateFormat, minimumAge)
	}
}

func regexValidatorProvider(regexString string) validator.Func {
	regex := regexp.MustCompile(regexString)
	return func(fl validator.FieldLevel) bool {
		return validateRegex(fl.Field().String(), regex)
	}
}

func ValidatorProvider(config *struct {
	DateFormat  string     `inject:"config:form.validator.dateFormat"`
	MinimumAge  float64    `inject:"config:form.validator.minimumAge"`
	MaximumAge  float64    `inject:"config:form.validator.maximumAge"`
	CustomRegex config.Map `inject:"config:form.validator.customRegex"`
}) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("dateformat", dateFormatValidatorProvider(config.DateFormat))
	validate.RegisterValidation("minimumage", minimumAgeValidatorProvider(config.DateFormat, int(config.MinimumAge)))
	validate.RegisterValidation("maximumage", maximumAgeValidatorProvider(config.DateFormat, int(config.MaximumAge)))
	validate.RegisterValidation("minimumnow", minimumAgeValidatorProvider(config.DateFormat, 0))
	validate.RegisterValidation("maximumnow", maximumAgeValidatorProvider(config.DateFormat, 0))

	for name, value := range config.CustomRegex {
		regex, ok := value.(string)
		if !ok {
			panic("wrong value passed as validation regex")
		}
		validate.RegisterValidation(name, regexValidatorProvider(regex))
	}

	return validate
}
