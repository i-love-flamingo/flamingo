package templatefunctions

import (
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/locale/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	DateTime struct {
		DateTimeService application.DateTimeService `inject:""`
		Logger          flamingo.Logger             `inject:""`
	}
)

// Name alias for use in template
func (tf DateTime) Name() string {
	return "dateTimeFormat"
}

func (tf DateTime) Func() interface{} {
	// Usage
	// dateTimeFormat(dateTimeString).formatDate()
	return func(dateTimeString string) domain.DateTime {

		dateTime, e := tf.DateTimeService.GetDateTimeFromIsoString(dateTimeString)
		if e != nil {
			tf.Logger.Errorf("Error Parsing dateTime %v / %v", dateTimeString, e)
			return domain.DateTime{}
		}
		tf.Logger.Errorf(" Parsing dateTime %v / %#v", dateTimeString, *dateTime)
		return *dateTime
	}
}
