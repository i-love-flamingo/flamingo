package validators

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"flamingo.me/flamingo/core/form2/domain/mocks"
)

type (
	RegexValidatorTestSuite struct {
		suite.Suite

		validator *RegexValidator
	}
)

func TestRegexValidatorTestSuite(t *testing.T) {
	suite.Run(t, &RegexValidatorTestSuite{})
}

func (t *RegexValidatorTestSuite) SetupTest() {
	t.validator = NewRegexValidator("onlynumber", "^[0-9]{1}$")
}

func (t *RegexValidatorTestSuite) TestValidatorName() {
	t.Equal("onlynumber", t.validator.ValidatorName())
}

func (t *RegexValidatorTestSuite) TestValidateField() {
	testCases := []struct {
		Value   string
		Result bool
	}{
		{
			Value:  "",
			Result: true,
		},
		{
			Value:  " ",
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
		fieldLevel := &mocks.FieldLevel{}
		fieldLevel.On("Field").Return(reflect.ValueOf(testCase.Value)).Once()
		t.Equal(testCase.Result, t.validator.ValidateField(nil, fieldLevel))
		fieldLevel.AssertExpectations(t.T())
	}
}
