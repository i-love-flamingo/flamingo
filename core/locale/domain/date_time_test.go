package domain

import (
	"flamingo.me/flamingo/v3/framework/flamingo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTimeFormatter_SetDateTime(t *testing.T) {
	formatter := &DateTimeFormatter{}
	assert.Equal(t, &DateTimeFormatter{}, formatter)

	now := getUTCNow()
	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)

	formatter.SetDateTime(now, now.In(loc))
	assert.Equal(t, &DateTimeFormatter{
		dateTime:      now,
		localDateTime: now.In(loc),
	}, formatter)
}

func TestDateTimeFormatter_SetLocation(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		logger:   flamingo.NullLogger{},
		dateTime: now,
	}

	assert.Errorf(t, formatter.SetLocation("wrong"), ErrInvalidLocation)

	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)
	assert.NoError(t, formatter.SetLocation(loc.String()))
	assert.Equal(t, &DateTimeFormatter{
		logger:        flamingo.NullLogger{},
		dateTime:      now,
		localDateTime: now.In(loc),
	}, formatter)
}

func TestDateTimeFormatter_SetLogger(t *testing.T) {
	formatter := &DateTimeFormatter{}
	assert.Nil(t, formatter.logger)

	formatter.SetLogger(flamingo.NullLogger{})
	assert.Equal(t, flamingo.NullLogger{}, formatter.logger)
}

func TestDateTimeFormatter_Format(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime: now,
	}
	assert.Equal(t, now.Format(time.RFC1123), formatter.Format(time.RFC1123))
}

func TestDateTimeFormatter_FormatLocale(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime: now,
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)
	assert.NoError(t, formatter.SetLocation(loc.String()))

	assert.Equal(t, now.In(loc).Format(time.RFC1123), formatter.FormatLocale(time.RFC1123))
}

func TestDateTimeFormatter_FormatDate(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:   now,
		DateFormat: "02 Jan 06",
	}
	assert.Equal(t, now.Format("02 Jan 06"), formatter.FormatDate())
}

func TestDateTimeFormatter_FormatTime(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:   now,
		TimeFormat: "15:04:05",
	}
	assert.Equal(t, now.Format("15:04:05"), formatter.FormatTime())
}

func TestDateTimeFormatter_FormatDateTime(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:       now,
		DateTimeFormat: time.RFC3339Nano,
	}
	assert.Equal(t, now.Format(time.RFC3339Nano), formatter.FormatDateTime())
}

func TestDateTimeFormatter_FormatToLocalDate(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:   now,
		DateFormat: "02 Jan 06",
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)
	assert.NoError(t, formatter.SetLocation(loc.String()))

	assert.Equal(t, now.In(loc).Format("02 Jan 06"), formatter.FormatToLocalDate())
}

func TestDateTimeFormatter_FormatToLocalTime(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:   now,
		TimeFormat: "15:04:05",
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)
	assert.NoError(t, formatter.SetLocation(loc.String()))

	assert.Equal(t, now.In(loc).Format("15:04:05"), formatter.FormatToLocalTime())
}

func TestDateTimeFormatter_FormatToLocalDateTime(t *testing.T) {
	now := getUTCNow()

	formatter := &DateTimeFormatter{
		dateTime:       now,
		DateTimeFormat: time.RFC3339Nano,
	}

	loc, err := time.LoadLocation("Europe/Berlin")
	assert.NoError(t, err)
	assert.NoError(t, formatter.SetLocation(loc.String()))

	assert.Equal(t, now.In(loc).Format(time.RFC3339Nano), formatter.FormatToLocalDateTime())
}

func getUTCNow() time.Time {
	start := time.Now()
	return time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
}
