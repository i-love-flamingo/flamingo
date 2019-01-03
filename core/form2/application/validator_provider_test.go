package application

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/validator.v9"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/core/form2/domain/mocks"
	"flamingo.me/flamingo/framework/web"
)

type (
	ValidatorProviderTestSuite struct {
		suite.Suite

		provider *ValidatorProviderImpl

		firstFieldValidator  *mocks.FieldValidator
		secondFieldValidator *mocks.FieldValidator

		structValidator *mocks.StructValidator
	}

	validatorProviderTestData struct {
		First  string `validate:"firstfield"`
		Second string `validate:"secondfield"`
	}
)

func TestValidatorProviderTestSuite(t *testing.T) {
	suite.Run(t, &ValidatorProviderTestSuite{})
}

func (t *ValidatorProviderTestSuite) SetupTest() {
	t.firstFieldValidator = &mocks.FieldValidator{}
	t.secondFieldValidator = &mocks.FieldValidator{}
	t.structValidator = &mocks.StructValidator{}
	t.provider = &ValidatorProviderImpl{}
	t.provider.Inject([]domain.FieldValidator{
		t.firstFieldValidator,
		t.secondFieldValidator,
	}, []domain.StructValidator{
		t.structValidator,
	})
}

func (t *ValidatorProviderTestSuite) TearDownTest() {
	t.firstFieldValidator.AssertExpectations(t.T())
	t.firstFieldValidator = nil
	t.secondFieldValidator.AssertExpectations(t.T())
	t.secondFieldValidator = nil
	t.structValidator.AssertExpectations(t.T())
	t.structValidator = nil
	t.provider = nil
}

func (t *ValidatorProviderTestSuite) TestGetValidator() {
	t.firstFieldValidator.On("ValidatorName").Return("firstfield").Once()
	t.secondFieldValidator.On("ValidatorName").Return("secondfield").Once()
	t.structValidator.On("StructType").Return(validatorProviderTestData{}).Once()

	t.IsType(&validator.Validate{}, t.provider.GetValidator())
}

func (t *ValidatorProviderTestSuite) TestErrorsToValidationInfo_Empty() {
	validationInfo := t.provider.ErrorsToValidationInfo(nil)
	t.True(validationInfo.IsValid())
}

func (t *ValidatorProviderTestSuite) TestErrorsToValidationInfo_FieldError() {
	err := &mocks.FieldError{}
	err.On("Namespace").Return("formData.fieldName1").Once()
	err.On("Tag").Return("firstfield").Twice()
	err.On("Field").Return("FieldName1").Once()

	validationInfo := t.provider.ErrorsToValidationInfo(validator.ValidationErrors{
		err,
	})
	t.False(validationInfo.IsValid())
	t.Equal(map[string][]domain.Error{
		"fieldName1": {
			{
				MessageKey:   "formError.fieldName1.firstfield",
				DefaultLabel: "FieldName1 firstfield",
			},
		},
	}, validationInfo.GetErrorsForAllFields())

	err.AssertExpectations(t.T())
}

func (t *ValidatorProviderTestSuite) TestErrorsToValidationInfo_GeneralError() {
	err := errors.New("error")
	validationInfo := t.provider.ErrorsToValidationInfo(err)
	t.False(validationInfo.IsValid())
	t.Equal([]domain.Error{
		{
			MessageKey:   "formError.invalidValidation",
			DefaultLabel: "error",
		},
	}, validationInfo.GetGeneralErrors())
}

func (t *ValidatorProviderTestSuite) TestGetRelativeFieldNameFromValidationError() {
	testCases := []struct {
		Namespace string
		Result    string
	}{
		{
			Namespace: "formData.fieldName1",
			Result:    "fieldName1",
		},
		{
			Namespace: "fieldName1",
			Result:    "fieldName1",
		},
		{
			Namespace: "formData.subData.fieldName1",
			Result:    "subData.fieldName1",
		},
	}

	for _, testCase := range testCases {
		err := &mocks.FieldError{}
		err.On("Namespace").Return(testCase.Namespace).Once()
		t.Equal(testCase.Result, t.provider.getRelativeFieldNameFromValidationError(err))
		err.AssertExpectations(t.T())
	}
}

func (t *ValidatorProviderTestSuite) TestValidate() {
	ctx := context.Background()
	request := &web.Request{}
	reqCtx := web.Context_(ctx, request)

	t.firstFieldValidator.On("ValidatorName").Return("firstfield").Once()
	t.secondFieldValidator.On("ValidatorName").Return("secondfield").Once()

	t.structValidator.On("StructType").Return(validatorProviderTestData{}).Once()

	t.firstFieldValidator.On("ValidateField", reqCtx, mock.Anything).Return(false).Once()
	t.secondFieldValidator.On("ValidateField", reqCtx, mock.Anything).Return(true).Once()

	t.structValidator.On("ValidateStruct", reqCtx, mock.Anything).Return().Once()

	validationInfo := t.provider.Validate(ctx, request, validatorProviderTestData{
		First:  "first",
		Second: "second",
	})
	t.False(validationInfo.IsValid())
	t.Equal(map[string][]domain.Error{
		"first": {
			{
				MessageKey:   "formError.first.firstfield",
				DefaultLabel: "First firstfield",
			},
		},
	}, validationInfo.GetErrorsForAllFields())
}
