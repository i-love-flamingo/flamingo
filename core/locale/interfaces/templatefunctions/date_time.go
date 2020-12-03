package templatefunctions

import (
	"context"
	"time"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// DateTimeFormatFromIso template helper function
type DateTimeFormatFromIso struct {
	dateTimeService *application.DateTimeService
	logger          flamingo.Logger
}

// DateTimeFormatFromTime template helper function
type DateTimeFormatFromTime struct {
	dateTimeService *application.DateTimeService
	logger          flamingo.Logger
}

// Inject dependencies
func (tf *DateTimeFormatFromIso) Inject(service *application.DateTimeService, logger flamingo.Logger) {
	tf.dateTimeService = service
	tf.logger = logger
}

// Func template function factory
func (tf *DateTimeFormatFromIso) Func(context.Context) interface{} {
	// Usage
	// dateTimeFormatFromIso(dateTimeString).formatDate()
	return func(dateTimeString string) *domain.DateTimeFormatter {
		dateTimeFormatter, e := tf.dateTimeService.GetDateTimeFormatterFromIsoString(dateTimeString)
		if e != nil {
			tf.logger.Error("Error Parsing dateTime %v / %v", dateTimeString, e)
			return &domain.DateTimeFormatter{}
		}
		return dateTimeFormatter
	}
}

// Inject dependencies
func (tf *DateTimeFormatFromTime) Inject(service *application.DateTimeService, logger flamingo.Logger) {
	tf.dateTimeService = service
	tf.logger = logger
}

// Func template function factory
func (tf *DateTimeFormatFromTime) Func(context.Context) interface{} {
	// Usage
	// dateTimeFormat(dateTime).formatDate()
	return func(dateTime time.Time) *domain.DateTimeFormatter {
		dateTimeFormatter, e := tf.dateTimeService.GetDateTimeFormatter(dateTime)
		if e != nil {
			tf.logger.Error("Error getting formatter dateTime %v", e)
			return &domain.DateTimeFormatter{}
		}
		return dateTimeFormatter
	}
}
