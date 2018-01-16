package templatefunctions

import (
	"time"

	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/locale/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	// DateTime template helper function
	DateTimeFormatFromIso struct {
		DateTimeService application.DateTimeService `inject:""`
		Logger          flamingo.Logger             `inject:""`
	}

	// DateTime template helper function
	DateTimeFormatFromTime struct {
		DateTimeService application.DateTimeService `inject:""`
		Logger          flamingo.Logger             `inject:""`
	}
)

// Name alias for use in template
func (tf DateTimeFormatFromIso) Name() string {
	return "dateTimeFormatFromIso"
}

// Func template function factory
func (tf DateTimeFormatFromIso) Func() interface{} {
	// Usage
	// dateTimeFormatFromIso(dateTimeString).formatDate()
	return func(dateTimeString string) *domain.DateTimeFormatter {
		dateTimeFormatter, e := tf.DateTimeService.GetDateTimeFormatterFromIsoString(dateTimeString)
		if e != nil {
			tf.Logger.Errorf("Error Parsing dateTime %v / %v", dateTimeString, e)
			return &domain.DateTimeFormatter{}
		}
		return dateTimeFormatter
	}
}

// Name alias for use in template
func (tf DateTimeFormatFromTime) Name() string {
	return "dateTimeFormat"
}

// Func template function factory
func (tf DateTimeFormatFromTime) Func() interface{} {
	// Usage
	// dateTimeFormat(dateTime).formatDate()
	return func(dateTime time.Time) *domain.DateTimeFormatter {
		dateTimeFormatter, e := tf.DateTimeService.GetDateTimeFormatter(dateTime)
		if e != nil {
			tf.Logger.Errorf("Error getting formatter dateTime %v", e)
			return &domain.DateTimeFormatter{}
		}
		return dateTimeFormatter
	}
}
