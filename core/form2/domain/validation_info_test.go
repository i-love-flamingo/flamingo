package domain

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type (
	ValidationInfoTestSuite struct {
		suite.Suite

		validationInfo ValidationInfo
	}
)

func TestValidationInfoTestSuite(t *testing.T) {
	suite.Run(t, &ValidationInfoTestSuite{})
}

func (t *ValidationInfoTestSuite) SetupTest() {
	t.validationInfo = ValidationInfo{}
}

func (t *ValidationInfoTestSuite) TestIsValid_Valid() {
	t.True(t.validationInfo.IsValid())
	t.False(t.validationInfo.HasAnyFieldErrors())
	t.False(t.validationInfo.HasErrorsForField("fieldName1"))
	t.False(t.validationInfo.HasErrorsForField("fieldName2"))
	t.False(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName2"))
	t.Equal(map[string][]Error{}, t.validationInfo.GetAllFieldErrors())
	t.Equal([]Error{}, t.validationInfo.GetGeneralErrors())
}

func (t *ValidationInfoTestSuite) TestIsValid_FieldError() {
	t.validationInfo.AddFieldError("fieldName1", "messageKey1", "defaultLabel1")

	t.False(t.validationInfo.IsValid())
	t.True(t.validationInfo.HasAnyFieldErrors())
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.False(t.validationInfo.HasErrorsForField("fieldName2"))
	t.False(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName2"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())
	t.Equal([]Error{}, t.validationInfo.GetGeneralErrors())
}

func (t *ValidationInfoTestSuite) TestIsValid_GeneralError() {
	t.validationInfo.AddGeneralError("messageKeyG", "defaultLabelG")

	t.False(t.validationInfo.IsValid())
	t.False(t.validationInfo.HasAnyFieldErrors())
	t.False(t.validationInfo.HasErrorsForField("fieldName1"))
	t.False(t.validationInfo.HasErrorsForField("fieldName2"))
	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName2"))
	t.Equal(map[string][]Error{}, t.validationInfo.GetAllFieldErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKeyG",
			DefaultLabel: "defaultLabelG",
		},
	}, t.validationInfo.GetGeneralErrors())
}

func (t *ValidationInfoTestSuite) TestAddGeneralError() {
	t.False(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{}, t.validationInfo.GetGeneralErrors())

	t.validationInfo.AddGeneralError("messageKey1", "defaultLabel1")

	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, t.validationInfo.GetGeneralErrors())

	t.validationInfo.AddGeneralError("messageKey2", "defaultLabel2")

	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetGeneralErrors())

	t.validationInfo.AddGeneralError("messageKey1", "defaultLabel1")

	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetGeneralErrors())
}

func (t *ValidationInfoTestSuite) TestAppendGeneralErrors() {
	t.False(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{}, t.validationInfo.GetGeneralErrors())

	t.validationInfo.AppendGeneralErrors([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	})

	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, t.validationInfo.GetGeneralErrors())

	t.validationInfo.AppendGeneralErrors([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	})

	t.True(t.validationInfo.HasGeneralErrors())
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetGeneralErrors())
}

func (t *ValidationInfoTestSuite) TestAddFieldError() {
	t.False(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{}, t.validationInfo.GetAllFieldErrors())

	t.validationInfo.AddFieldError("fieldName1", "messageKey1", "defaultLabel1")
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())

	t.validationInfo.AddFieldError("fieldName1", "messageKey2", "defaultLabel2")
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
			{
				MessageKey:   "messageKey2",
				DefaultLabel: "defaultLabel2",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())

	t.validationInfo.AddFieldError("fieldName1", "messageKey1", "defaultLabel1")
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
			{
				MessageKey:   "messageKey2",
				DefaultLabel: "defaultLabel2",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())
}

func (t *ValidationInfoTestSuite) TestAppendFieldErrors() {
	t.False(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{}, t.validationInfo.GetFieldErrors("fieldName1"))

	t.validationInfo.AppendFieldErrors(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
		},
	})
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())

	t.validationInfo.AppendFieldErrors(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
			{
				MessageKey:   "messageKey2",
				DefaultLabel: "defaultLabel2",
			},
		},
	})
	t.True(t.validationInfo.HasErrorsForField("fieldName1"))
	t.Equal([]Error{
		{
			MessageKey:   "messageKey1",
			DefaultLabel: "defaultLabel1",
		},
		{
			MessageKey:   "messageKey2",
			DefaultLabel: "defaultLabel2",
		},
	}, t.validationInfo.GetFieldErrors("fieldName1"))
	t.Equal(map[string][]Error{
		"fieldName1": {
			{
				MessageKey:   "messageKey1",
				DefaultLabel: "defaultLabel1",
			},
			{
				MessageKey:   "messageKey2",
				DefaultLabel: "defaultLabel2",
			},
		},
	}, t.validationInfo.GetAllFieldErrors())
}
