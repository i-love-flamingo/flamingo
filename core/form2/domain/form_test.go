package domain

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type (
	FormTestSuite struct {
		suite.Suite
	}
)

func TestFormTestSuite(t *testing.T) {
	suite.Run(t, &FormTestSuite{})
}

func (t *FormTestSuite) TestNewForm() {
	form := NewForm(false, map[string][]ValidationRule{
		"fieldName1": {
			{
				Name:  "gte",
				Value: "10",
			},
		},
	})

	t.False(form.IsSubmitted())
	t.Equal([]ValidationRule{
		{
			Name:  "gte",
			Value: "10",
		},
	}, form.GetValidationRulesForField("fieldName1"))

	form = NewForm(true, map[string][]ValidationRule{
		"fieldName1": {
			{
				Name:  "gte",
				Value: "10",
			},
			{
				Name: "required",
			},
		},
	})

	t.True(form.IsSubmitted())
	t.Equal([]ValidationRule{
		{
			Name:  "gte",
			Value: "10",
		},
		{
			Name: "required",
		},
	}, form.GetValidationRulesForField("fieldName1"))
}

func (t *FormTestSuite) TestIsValidAndSubmitted() {
	form := NewForm(false, map[string][]ValidationRule{})
	t.True(form.IsValid())
	t.False(form.IsSubmitted())
	t.False(form.IsValidAndSubmitted())

	form = NewForm(true, map[string][]ValidationRule{})
	t.True(form.IsValid())
	t.True(form.IsSubmitted())
	t.True(form.IsValidAndSubmitted())

	validationInfo := ValidationInfo{}
	validationInfo.AddGeneralError("messageKey1", "defaultLabel1")
	form.ValidationInfo = validationInfo
	t.False(form.IsValid())
	t.True(form.IsSubmitted())
	t.False(form.IsValidAndSubmitted())
}

func (t *FormTestSuite) TestErrors() {
	form := NewForm(false, map[string][]ValidationRule{})
	t.False(form.HasAnyFieldErrors())
	t.False(form.HasGeneralErrors())
	t.False(form.HasErrorForField("fieldName1"))

	validationInfo := ValidationInfo{}
	validationInfo.AddGeneralError("messageKey1", "defaultLabel1")
	validationInfo.AddFieldError("fieldName1", "messageKey1", "defaultLabel1")
	form.ValidationInfo = validationInfo
	t.True(form.HasAnyFieldErrors())
	t.True(form.HasGeneralErrors())
	t.True(form.HasErrorForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, form.GetGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, form.GetErrorsForField("fieldName1"))
}
