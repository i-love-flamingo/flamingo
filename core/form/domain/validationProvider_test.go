package domain

import (
	"testing"
	"time"

	"flamingo.me/flamingo/framework/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"
)

type (
	ValidationProviderTestSuite struct {
		suite.Suite

		validate *validator.Validate
	}
)

func TestRunValidationProviderTestSuite(t *testing.T) {
	suite.Run(t, &ValidationProviderTestSuite{})
}

func (t *ValidationProviderTestSuite) SetupTest() {
	cfg := struct {
		DateFormat  string     `inject:"config:form.validator.dateFormat"`
		CustomRegex config.Map `inject:"config:form.validator.customRegex"`
	}{
		DateFormat: "2006-01-02",
		CustomRegex: config.Map{
			"onlynumber": "^[0-9]{1}$",
			"justthis":   "^justthis$",
		},
	}
	t.validate = ValidatorProvider([]FieldValidator{}, &cfg)
}

func (t *ValidationProviderTestSuite) TestDateFormat() {
	formData := struct {
		Date string `validate:"dateformat"`
	}{}

	testCases := []struct {
		Value  string
		Result bool
	}{
		{
			Value:  "",
			Result: false,
		},
		{
			Value:  "2006-30-02",
			Result: false,
		},
		{
			Value:  "2006.30.02",
			Result: false,
		},
		{
			Value:  "2006-30-13",
			Result: false,
		},
		{
			Value:  "2006-00-01",
			Result: false,
		},
		{
			Value:  "2006-01-00",
			Result: false,
		},
		{
			Value:  "1900-02-29",
			Result: false,
		},
		{
			Value:  "2006-01-31",
			Result: true,
		},
		{
			Value:  "2000-02-29",
			Result: true,
		},
	}

	for _, testCase := range testCases {
		formData.Date = testCase.Value
		err := t.validate.Struct(formData)
		if testCase.Result {
			t.NoError(err)
		} else {
			t.Error(err)
			fieldErrors, ok := err.(validator.ValidationErrors)
			t.True(ok)
			t.Len(fieldErrors, 1)
			t.Equal("dateformat", fieldErrors[0].Tag())
			t.Equal("Date", fieldErrors[0].Field())
		}
	}
}

func (t *ValidationProviderTestSuite) TestMinimumAge() {
	now := time.Now()
	justAdult := time.Date(now.Year()-18, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	almostAdult := justAdult.Add(24 * time.Hour)

	formData := struct {
		Date string `validate:"minimumage=18"`
	}{}

	testCases := []struct {
		Value  string
		Result bool
	}{
		{
			Value:  "",
			Result: true,
		},
		{
			Value:  almostAdult.Format("2006-01-02"),
			Result: false,
		},
		{
			Value:  justAdult.Format("2006-01-02"),
			Result: true,
		},
	}

	for _, testCase := range testCases {
		formData.Date = testCase.Value
		err := t.validate.Struct(formData)
		if testCase.Result {
			t.NoError(err)
		} else {
			t.Error(err)
			fieldErrors, ok := err.(validator.ValidationErrors)
			t.True(ok)
			t.Len(fieldErrors, 1)
			t.Equal("minimumage", fieldErrors[0].Tag())
			t.Equal("Date", fieldErrors[0].Field())
		}
	}
}

func (t *ValidationProviderTestSuite) TestMaximumAge() {
	now := time.Now()
	child := time.Date(now.Year()-7, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	adult := time.Date(now.Year()-40, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	walkingDead := time.Date(now.Year()-200, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	formData := struct {
		Date string `validate:"maximumage=150"`
	}{}

	testCases := []struct {
		Value  string
		Result bool
	}{
		{
			Value:  "",
			Result: true,
		},
		{
			Value:  child.Format("2006-01-02"),
			Result: true,
		},
		{
			Value:  adult.Format("2006-01-02"),
			Result: true,
		},
		{
			Value:  walkingDead.Format("2006-01-02"),
			Result: false,
		},
	}

	for _, testCase := range testCases {
		formData.Date = testCase.Value
		err := t.validate.Struct(formData)
		if testCase.Result {
			t.NoError(err)
		} else {
			t.Error(err)
			fieldErrors, ok := err.(validator.ValidationErrors)
			t.True(ok)
			t.Len(fieldErrors, 1)
			t.Equal("maximumage", fieldErrors[0].Tag())
			t.Equal("Date", fieldErrors[0].Field())
		}
	}
}

func (t *ValidationProviderTestSuite) TestRegexOne() {
	formData := struct {
		Number string `validate:"onlynumber"`
	}{}

	testCases := []struct {
		Value  string
		Result bool
	}{
		{
			Value:  "",
			Result: false,
		},
		{
			Value:  "a",
			Result: false,
		},
		{
			Value:  "A",
			Result: false,
		},
		{
			Value:  "00",
			Result: false,
		},
		{
			Value:  "1",
			Result: true,
		},
	}

	for _, testCase := range testCases {
		formData.Number = testCase.Value
		err := t.validate.Struct(formData)
		if testCase.Result {
			t.NoError(err)
		} else {
			t.Error(err)
			fieldErrors, ok := err.(validator.ValidationErrors)
			t.True(ok)
			t.Len(fieldErrors, 1)
			t.Equal("onlynumber", fieldErrors[0].Tag())
			t.Equal("Number", fieldErrors[0].Field())
		}
	}
}

func (t *ValidationProviderTestSuite) TestRegexTwo() {
	formData := struct {
		JustThis string `validate:"justthis"`
	}{}

	testCases := []struct {
		Value  string
		Result bool
	}{
		{
			Value:  "",
			Result: false,
		},
		{
			Value:  "justthis",
			Result: true,
		},
	}

	for _, testCase := range testCases {
		formData.JustThis = testCase.Value
		err := t.validate.Struct(formData)
		if testCase.Result {
			t.NoError(err)
		} else {
			t.Error(err)
			fieldErrors, ok := err.(validator.ValidationErrors)
			t.True(ok)
			t.Len(fieldErrors, 1)
			t.Equal("justthis", fieldErrors[0].Tag())
			t.Equal("JustThis", fieldErrors[0].Field())
		}
	}
}
