package validators

import (
	"flamingo.me/flamingo/core/form2/application/provider/validators/mocks"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type (
	MinimumAgeValidatorTestSuite struct {
		suite.Suite

		validator *MinimumAgeValidator
	}
)

func TestMinimumAgeValidatorTestSuite(t *testing.T) {
	suite.Run(t, &MinimumAgeValidatorTestSuite{})
}

func (t *MinimumAgeValidatorTestSuite) SetupTest() {
	t.validator = &MinimumAgeValidator{}
	t.validator.Inject(&struct {
		DateFormat string `inject:"config:form.validator.dateFormat"`
	}{
		DateFormat: "2006-01-02",
	})
}

func (t *MinimumAgeValidatorTestSuite) TestValidatorName() {
	t.Equal("minimumage", t.validator.ValidatorName())
}

func (t *MinimumAgeValidatorTestSuite) TestValidateField() {
	now := time.Now()
	child := time.Date(now.Year()-7, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	justAdult := time.Date(now.Year()-18, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	almostAdult := justAdult.Add(24 * time.Hour)
	adult := time.Date(now.Year()-40, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

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
			Result: false,
		},
		{
			Date:  almostAdult.Format("2006-01-02"),
			Result: false,
		},
		{
			Date:  justAdult.Format("2006-01-02"),
			Result: true,
		},
		{
			Date:  adult.Format("2006-01-02"),
			Result: true,
		},
	}

	for _, testCase := range testCases {
		fieldLevel := &mocks.FieldLevel{}
		fieldLevel.On("Field").Return(reflect.ValueOf(testCase.Date)).Once()
		fieldLevel.On("Param").Return("18").Once()
		t.Equal(testCase.Result, t.validator.ValidateField(nil, fieldLevel))
		fieldLevel.AssertExpectations(t.T())
	}
}
