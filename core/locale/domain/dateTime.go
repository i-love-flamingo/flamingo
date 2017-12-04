package domain

import (
	"time"

	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	DateTime struct {
		DateFormat     string
		TimeFormat     string
		DateTimeFormat string
		Location       string
		Logger         flamingo.Logger
		dateTime       time.Time
	}
)

//SetDateTime Setter for private member
func (dt *DateTime) SetDateTime(time time.Time) {
	dt.dateTime = time
}

func (dts *DateTime) Format(format string) string {
	return dts.dateTime.Format(format)
}

func (dts *DateTime) FormatLocale(format string) string {
	return dts.getLocalTime().Format(format)
}

func (dts *DateTime) FormatDate() string {
	dts.Logger.Errorf("ääääääääää %v", dts.dateTime)
	return dts.dateTime.Format(dts.DateFormat)
}

func (dts *DateTime) FormatTime(dateTimeString string) string {
	return dts.dateTime.Format(dts.TimeFormat)
}

func (dts *DateTime) FormatDateTime(dateTimeString string) string {
	return dts.dateTime.Format(dts.DateTimeFormat)
}

func (dts *DateTime) FormatToLocalDate() string {
	return dts.getLocalTime().Format(dts.DateFormat)
}

func (dts *DateTime) FormatToLocalTime(dateTimeString string) string {
	return dts.getLocalTime().Format(dts.TimeFormat)
}

func (dts *DateTime) FormatToLocalDateTime(dateTimeString string) string {
	return dts.getLocalTime().Format(dts.DateTimeFormat)
}

func (dts *DateTime) getLocalTime() time.Time {
	loc, e := time.LoadLocation(dts.Location)
	if e != nil {
		if dts.Logger != nil {
			dts.Logger.Errorf("dateTime Parsing error - could not load location %v", dts.Location)
		}
		return dts.dateTime
	}
	return dts.dateTime.In(loc)
}
