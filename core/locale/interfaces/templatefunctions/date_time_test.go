package templatefunctions

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestDateTimeFormatFromIso_Func(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)

	service := &application.DateTimeService{}
	service.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		DateFormat:     "2006-01-02",
		TimeFormat:     "15:04:05Z07:00",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
		Location:       loc.String(),
	})

	tFuncProvider := &DateTimeFormatFromIso{}
	tFuncProvider.Inject(service, flamingo.NullLogger{})

	tFunc, ok := tFuncProvider.Func(context.Background()).(func(dateTimeString string) *domain.DateTimeFormatter)
	assert.True(t, ok)

	// error case
	formatter := tFunc("wrong")
	assert.Equal(t, &domain.DateTimeFormatter{}, formatter)

	// generated formatter case
	start := time.Now()
	now, err := time.Parse(time.RFC3339, start.Format(time.RFC3339))
	assert.NoError(t, err)

	formatter = tFunc(now.Format(time.RFC3339))
	expected := &domain.DateTimeFormatter{
		DateFormat:     "2006-01-02",
		TimeFormat:     "15:04:05Z07:00",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
	}
	assert.NoError(t, expected.SetLocation(loc.String()))
	expected.SetDateTime(now, now.In(loc))
	expected.SetLogger(flamingo.NullLogger{})
	assert.Equal(t, expected, formatter)
}

func TestDateTimeFormatFromTime_Func(t *testing.T) {
	start := time.Now()
	now, err := time.Parse(time.RFC3339, start.Format(time.RFC3339))
	assert.NoError(t, err)

	// error case
	service := &application.DateTimeService{}
	service.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		DateFormat:     "2006-01-02",
		TimeFormat:     "15:04:05Z07:00",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
		Location:       "wrong",
	})

	tFuncProvider := &DateTimeFormatFromTime{}
	tFuncProvider.Inject(service, flamingo.NullLogger{})

	tFunc, ok := tFuncProvider.Func(context.Background()).(func(dateTime time.Time) *domain.DateTimeFormatter)
	assert.True(t, ok)

	formatter := tFunc(now)
	assert.Equal(t, &domain.DateTimeFormatter{}, formatter)

	// generated formatter case
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)

	service = &application.DateTimeService{}
	service.Inject(flamingo.NullLogger{}, &struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	}{
		DateFormat:     "2006-01-02",
		TimeFormat:     "15:04:05Z07:00",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
		Location:       loc.String(),
	})

	tFuncProvider = &DateTimeFormatFromTime{}
	tFuncProvider.Inject(service, flamingo.NullLogger{})

	tFunc, ok = tFuncProvider.Func(context.Background()).(func(dateTime time.Time) *domain.DateTimeFormatter)
	assert.True(t, ok)

	formatter = tFunc(now)
	expected := &domain.DateTimeFormatter{
		DateFormat:     "2006-01-02",
		TimeFormat:     "15:04:05Z07:00",
		DateTimeFormat: "2006-01-02T15:04:05Z07:00",
	}
	assert.NoError(t, expected.SetLocation(loc.String()))
	expected.SetDateTime(now, now.In(loc))
	expected.SetLogger(flamingo.NullLogger{})
	assert.Equal(t, expected, formatter)
}
