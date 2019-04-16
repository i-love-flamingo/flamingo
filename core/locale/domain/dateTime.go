package domain

import (
	"fmt"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

// DateTimeFormatter has a couple of helpful methods to format date and times
type DateTimeFormatter struct {
	logger         flamingo.Logger
	DateFormat     string
	TimeFormat     string
	DateTimeFormat string
	localDateTime  time.Time
	dateTime       time.Time
}

// SetDateTime setter for private member
func (dtf *DateTimeFormatter) SetDateTime(time time.Time, localtime time.Time) {
	dtf.dateTime = time
	dtf.localDateTime = localtime
}

// SetLocation sets location for local time
func (dtf *DateTimeFormatter) SetLocation(loc string) *DateTimeFormatter {
	location, err := time.LoadLocation(loc)

	if err != nil && dtf.logger != nil {
		dtf.logger.WithField(flamingo.LogKeyMethod, "SetLocation").Error(
			fmt.Sprintf("%s: failed to load time zone location: %v", loc, err),
		)
	} else {
		dtf.localDateTime = dtf.dateTime.In(location)
	}

	return dtf
}

// SetLogger sets the logger instance
func (dtf *DateTimeFormatter) SetLogger(logger flamingo.Logger) {
	dtf.logger = logger
	dtf.logger.WithField(flamingo.LogKeyCategory, "DateTimeFormatter")
}

// Format datetime
func (dtf *DateTimeFormatter) Format(format string) string {
	return dtf.dateTime.Format(format)
}

// FormatLocale formats the local time
func (dtf *DateTimeFormatter) FormatLocale(format string) string {
	return dtf.localDateTime.Format(format)
}

// FormatDate formats the date
func (dtf *DateTimeFormatter) FormatDate() string {
	return dtf.dateTime.Format(dtf.DateFormat)
}

// FormatTime formats the time
func (dtf *DateTimeFormatter) FormatTime() string {
	return dtf.dateTime.Format(dtf.TimeFormat)
}

// FormatDateTime formats both date and time
func (dtf *DateTimeFormatter) FormatDateTime() string {
	return dtf.dateTime.Format(dtf.DateTimeFormat)
}

// FormatToLocalDate formats for local date
func (dtf *DateTimeFormatter) FormatToLocalDate() string {
	return dtf.localDateTime.Format(dtf.DateFormat)
}

// FormatToLocalTime formats the local time
func (dtf *DateTimeFormatter) FormatToLocalTime() string {
	return dtf.localDateTime.Format(dtf.TimeFormat)
}

// FormatToLocalDateTime formats both locale date and time
func (dtf *DateTimeFormatter) FormatToLocalDateTime() string {
	return dtf.localDateTime.Format(dtf.DateTimeFormat)
}
