package templatefunctions

import (
	"time"

	"flamingo.me/flamingo/core/locale/application"
	"flamingo.me/flamingo/core/locale/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// DateTime template helper function
	DateTimeFormatFromIso struct {
		dateTimeService *application.DateTimeService
		logger          flamingo.Logger
	}

	// DateTime template helper function
	DateTimeFormatFromTime struct {
		dateTimeService *application.DateTimeService
		logger          flamingo.Logger
	}
)

func (tf *DateTimeFormatFromIso) Inject(service *application.DateTimeService, logger flamingo.Logger) {
	tf.dateTimeService = service
	tf.logger = logger
}

// Func template function factory
func (tf *DateTimeFormatFromIso) Func() interface{} {
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

func (tf *DateTimeFormatFromTime) Inject(service *application.DateTimeService, logger flamingo.Logger) {
	tf.dateTimeService = service
	tf.logger = logger
}

// Func template function factory
func (tf *DateTimeFormatFromTime) Func() interface{} {
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
