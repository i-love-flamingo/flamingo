package validators

import (
	"flamingo.me/flamingo/core/form2/application/provider/validators/mocks"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type (
	MaximumAgeValidatorTestSuite struct {
		suite.Suite

		validator *MaximumAgeValidator
	}
)

func TestMaximumAgeValidatorTestSuite(t *testing.T) {
	suite.Run(t, &MaximumAgeValidatorTestSuite{})
}

func (t *MaximumAgeValidatorTestSuite) SetupTest() {
	t.validator = &MaximumAgeValidator{}
	t.validator.Inject(&struct {
		DateFormat string `inject:"config:form.validator.dateFormat"`
	}{
		DateFormat: "2006-01-02",
	})
}

func (t *MaximumAgeValidatorTestSuite) TestValidatorName() {
	t.Equal("maximumage", t.validator.ValidatorName())
}

func (t *MaximumAgeValidatorTestSuite) TestValidateField() {
	now := time.Now()
	child := time.Date(now.Year()-7, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	adult := time.Date(now.Year()-40, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	walkingDead := time.Date(now.Year()-150, now.Month(), now.Day() - 1, 0, 0, 0, 0, now.Location())

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
			Result: true,
		},
		{
			Date:  child.Format("2006-01-02"),
			Result: true,
		},
		{
			Date:  adult.Format("2006-01-02"),
			Result: true,
		},
		{
			Date:  walkingDead.Add(24 * time.Hour).Format("2006-01-02"),
			Result: true,
		},
		{
			Date:  walkingDead.Format("2006-01-02"),
			Result: false,
		},
	}

	for _, testCase := range testCases {
		fieldLevel := &mocks.FieldLevel{}
		fieldLevel.On("Field").Return(reflect.ValueOf(testCase.Date)).Once()
		fieldLevel.On("Param").Return("150").Once()
		t.Equal(testCase.Result, t.validator.ValidateField(nil, fieldLevel))
		fieldLevel.AssertExpectations(t.T())
	}
}
