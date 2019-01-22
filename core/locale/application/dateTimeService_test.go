package application

import (
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

func TestGetDateTimeFormatterFromIsoString(t *testing.T) {
	dateTimeService := new(DateTimeService)

	// check getting a date time formatter for a broken iso string
	f, e := dateTimeService.GetDateTimeFormatterFromIsoString("error")
	assert.Nil(t, f, "no formatter returned")
	assert.NotNil(t, e, "error returned")

	f, e = dateTimeService.GetDateTimeFormatterFromIsoString("2018-01-02T12:22:33Z")
	assert.NotNil(t, f, "formatter returned")
	assert.Nil(t, e, "no error returned")
}

func TestGetTimeFormatter(t *testing.T) {
	dateTimeService := DateTimeService{
		logger: flamingo.NullLogger{},
	}

	now := time.Now()

	// just get a plain formatter
	f, e := dateTimeService.GetDateTimeFormatter(now)
	assert.NotNil(t, f, "got a formatter")
	assert.Nil(t, e, "no error received")

	// get a formatter for a configured locale
	dateTimeService.location = "America/New_York"
	f, e = dateTimeService.GetDateTimeFormatter(now)
	assert.NotNil(t, f, "got a formatter")
	assert.Nil(t, e, "no error received")
}
