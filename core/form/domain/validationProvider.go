package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"flamingo.me/flamingo/framework/config"
	"gopkg.in/go-playground/validator.v9"
)

type (
	FormValidator interface {
		ValidatorName() string
	}

	FormValidatorWithoutParam interface {
		Validate(interface{}) bool
	}

	FormValidatorWithParam interface {
		Validate(string, interface{}) bool
	}
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

func extractAgeParam(name string, param string) int {
	age, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("wrong format for %s parameter: %s", name, param))
	}

	return int(age)
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

func maximumAgeValidatorProvider(dateFormat string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateMaximumAge(fl.Field().String(), dateFormat, extractAgeParam("maximumage", fl.Param()))
	}
}

func minimumAgeValidatorProvider(dateFormat string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return validateMinimumAge(fl.Field().String(), dateFormat, extractAgeParam("minimumage", fl.Param()))
	}
}

func regexValidatorProvider(regexString string) validator.Func {
	regex := regexp.MustCompile(regexString)
	return func(fl validator.FieldLevel) bool {
		return validateRegex(fl.Field().String(), regex)
	}
}

func ValidatorProvider(formValidators []FormValidator, config *struct {
	DateFormat  string     `inject:"config:form.validator.dateFormat"`
	CustomRegex config.Map `inject:"config:form.validator.customRegex"`
}) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("dateformat", dateFormatValidatorProvider(config.DateFormat))
	validate.RegisterValidation("minimumage", minimumAgeValidatorProvider(config.DateFormat))
	validate.RegisterValidation("maximumage", maximumAgeValidatorProvider(config.DateFormat))

	for name, value := range config.CustomRegex {
		regex, ok := value.(string)
		if !ok {
			panic("wrong value passed as validation regex")
		}
		validate.RegisterValidation(name, regexValidatorProvider(regex))
	}

	for _, formValidator := range formValidators {
		attached := false
		if withoutParam, ok := formValidator.(FormValidatorWithoutParam); ok {
			attached = true
			validate.RegisterValidation(formValidator.ValidatorName(), func(fl validator.FieldLevel) bool {
				return withoutParam.Validate(fl.Field().Interface())
			})
		}
		if withParam, ok := formValidator.(FormValidatorWithParam); ok {
			attached = true
			validate.RegisterValidation(formValidator.ValidatorName(), func(fl validator.FieldLevel) bool {
				return withParam.Validate(fl.Param(), fl.Field().Interface())
			})
		}
		if !attached {
			panic("Validator must implement either FormValidatorWithoutParam or FormValidatorWithParam interface")
		}
	}

	return validate
}
