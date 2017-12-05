package domain

import (
	"time"
)

type (
	// DateTimeFormatter has a couple of helpful methods to format date and times
	DateTimeFormatter struct {
		DateFormat     string
		TimeFormat     string
		DateTimeFormat string
		localDateTime  time.Time
		dateTime       time.Time
	}
)

//SetDateTime Setter for private member
func (dt *DateTimeFormatter) SetDateTime(time time.Time, localtime time.Time) {
	dt.dateTime = time
	dt.localDateTime = localtime
}

// Format datetime
func (dts *DateTimeFormatter) Format(format string) string {
	return dts.dateTime.Format(format)
}

// FormatLocale formats the local time
func (dts *DateTimeFormatter) FormatLocale(format string) string {
	return dts.localDateTime.Format(format)
}

// FormatDate formats the date
func (dts *DateTimeFormatter) FormatDate() string {
	return dts.dateTime.Format(dts.DateFormat)
}

// FormatTime formats the time
func (dts *DateTimeFormatter) FormatTime() string {
	return dts.dateTime.Format(dts.TimeFormat)
}

// FormatDateTime formats both date and time
func (dts *DateTimeFormatter) FormatDateTime() string {
	return dts.dateTime.Format(dts.DateTimeFormat)
}

// FormatToLocalDate formats for local date
func (dts *DateTimeFormatter) FormatToLocalDate() string {
	return dts.localDateTime.Format(dts.DateFormat)
}

// FormatToLocalTime formats the local time
func (dts *DateTimeFormatter) FormatToLocalTime() string {
	return dts.localDateTime.Format(dts.TimeFormat)
}

// FormatToLocalDateTime formats both locale date and time
func (dts *DateTimeFormatter) FormatToLocalDateTime() string {
	return dts.localDateTime.Format(dts.DateTimeFormat)
}
