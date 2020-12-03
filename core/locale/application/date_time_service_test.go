package application

import (
	"flamingo.me/flamingo/v3/core/locale/domain"
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

func TestDateTimeService_GetDateTimeFormatterFromIsoString_InvalidTime(t *testing.T) {
	formatter, err := new(DateTimeService).GetDateTimeFormatterFromIsoString("invalid time")
	assert.Nil(t, formatter, "expected nil for invalid time")
	assert.Error(t, err, "error received")
}

func TestDateTimeService_GetDateTimeFormatterFromIsoString_InvalidLocation(t *testing.T) {
	dateTimeService := DateTimeService{}
	dateTimeService.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		Location: "invalid location",
	})
}

func TestDateTimeService_GetDateTimeFormatterFromIsoString_PlainFormatter(t *testing.T) {
	start := time.Now()
	now, err := time.Parse(time.RFC3339, start.Format(time.RFC3339))
	assert.NoError(t, err, "no error received")

	plain := &domain.DateTimeFormatter{}
	plain.SetDateTime(now, now.In(time.UTC))

	formatter, err := new(DateTimeService).GetDateTimeFormatterFromIsoString(now.Format(time.RFC3339))
	assert.Equal(t, plain, formatter, "got a formatter")
	assert.NoError(t, err, "no error received")
}

func TestDateTimeService_GetDateTimeFormatterFromIsoString_FormatterForConfiguredLocale(t *testing.T) {
	start := time.Now()
	now, err := time.Parse(time.RFC3339, start.Format(time.RFC3339))
	assert.NoError(t, err, "no error received")

	dateTimeService := DateTimeService{
		logger: flamingo.NullLogger{},
	}
	dateTimeService.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		DateFormat:     time.RFC822,
		TimeFormat:     time.Kitchen,
		DateTimeFormat: time.ANSIC,
		Location:       "America/New_York",
	})

	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err, "no error received")

	result := &domain.DateTimeFormatter{
		DateFormat:     time.RFC822,
		TimeFormat:     time.Kitchen,
		DateTimeFormat: time.ANSIC,
	}
	result.SetDateTime(now, now.In(loc))
	result.SetLogger(flamingo.NullLogger{})

	formatter, err := dateTimeService.GetDateTimeFormatterFromIsoString(now.Format(time.RFC3339))
	assert.Equal(t, result, formatter, "got a formatter")
	assert.NoError(t, err, "no error received")
}

func TestDateTimeService_GetTimeFormatter_InvalidLocation(t *testing.T) {
	now := time.Now()

	dateTimeService := DateTimeService{}
	dateTimeService.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		Location: "invalid location",
	})

	formatter, err := dateTimeService.GetDateTimeFormatter(now)
	assert.Nil(t, formatter, "got nil for formatter")
	assert.Error(t, err, "error received")
}

func TestDateTimeService_GetTimeFormatter_PlainFormatter(t *testing.T) {
	now := time.Now()

	plain := &domain.DateTimeFormatter{}
	plain.SetDateTime(now, now.In(time.UTC))

	formatter, err := new(DateTimeService).GetDateTimeFormatter(now)
	assert.Equal(t, plain, formatter, "got a formatter")
	assert.NoError(t, err, "no error received")
}

func TestDateTimeService_GetTimeFormatter_FormatterForConfiguredLocale(t *testing.T) {
	now := time.Now()

	dateTimeService := DateTimeService{
		logger: flamingo.NullLogger{},
	}
	dateTimeService.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		DateFormat:     time.RFC822,
		TimeFormat:     time.Kitchen,
		DateTimeFormat: time.ANSIC,
		Location:       "America/New_York",
	})

	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err, "no error received")

	expectedResult := &domain.DateTimeFormatter{
		DateFormat:     time.RFC822,
		TimeFormat:     time.Kitchen,
		DateTimeFormat: time.ANSIC,
	}
	expectedResult.SetDateTime(now, now.In(loc))
	expectedResult.SetLogger(flamingo.NullLogger{})

	formatter, err := dateTimeService.GetDateTimeFormatter(now)
	assert.Equal(t, expectedResult, formatter, "got a formatter")
	assert.NoError(t, err, "no error received")
}
