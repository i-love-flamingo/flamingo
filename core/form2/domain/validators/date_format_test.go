package validators

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	"flamingo.me/flamingo/core/form2/domain/mocks"
)

type (
	DateFormatValidatorTestSuite struct {
		suite.Suite

		validator *DateFormatValidator
	}
)

func TestDateFormatValidatorTestSuite(t *testing.T) {
	suite.Run(t, &DateFormatValidatorTestSuite{})
}

func (t *DateFormatValidatorTestSuite) SetupTest() {
	t.validator = &DateFormatValidator{}
	t.validator.Inject(&struct {
		DateFormat string `inject:"config:form.validator.dateFormat"`
	}{
		DateFormat: "2006-01-02",
	})
}

func (t *DateFormatValidatorTestSuite) TestValidatorName() {
	t.Equal("dateformat", t.validator.ValidatorName())
}

func (t *DateFormatValidatorTestSuite) TestValidateField() {
	testCases := []struct {
		Date   string
		Result bool
	}{
		{
			Date:   "",
			Result: true,
		},
		{
			Date:   " ",
			Result: true,
		},
		{
			Date:   "wrong",
			Result: false,
		},
		{
			Date:   "2006-01-02",
			Result: true,
		},
		{
			Date:   "2006-30-02",
			Result: false,
		},
		{
			Date:   "2006.30.02",
			Result: false,
		},
		{
			Date:   "2006-30-13",
			Result: false,
		},
		{
			Date:   "2006-00-01",
			Result: false,
		},
		{
			Date:   "2006-01-00",
			Result: false,
		},
		{
			Date:   "1900-02-29",
			Result: false,
		},
		{
			Date:   "2006-01-31",
			Result: true,
		},
		{
			Date:   "2000-02-29",
			Result: true,
		},
	}

	for _, testCase := range testCases {
		fieldLevel := &mocks.FieldLevel{}
		fieldLevel.On("Field").Return(reflect.ValueOf(testCase.Date)).Once()
		t.Equal(testCase.Result, t.validator.ValidateField(nil, fieldLevel))
		fieldLevel.AssertExpectations(t.T())
	}
}
