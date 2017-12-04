package domain

import (
	"time"
)

type (
	DateTimeFormatter struct {
		DateFormat     string
		TimeFormat     string
		DateTimeFormat string

		localDateTime time.Time
		dateTime      time.Time
	}
)

//SetDateTime Setter for private member
func (dt *DateTimeFormatter) SetDateTime(time time.Time, localtime time.Time) {
	dt.dateTime = time
	dt.localDateTime = localtime
}

func (dts *DateTimeFormatter) Format(format string) string {
	return dts.dateTime.Format(format)
}

func (dts *DateTimeFormatter) FormatLocale(format string) string {
	return dts.localDateTime.Format(format)
}

func (dts *DateTimeFormatter) FormatDate() string {
	return dts.dateTime.Format(dts.DateFormat)
}

func (dts *DateTimeFormatter) FormatTime() string {
	return dts.dateTime.Format(dts.TimeFormat)
}

func (dts *DateTimeFormatter) FormatDateTime() string {
	return dts.dateTime.Format(dts.DateTimeFormat)
}

func (dts *DateTimeFormatter) FormatToLocalDate() string {
	return dts.localDateTime.Format(dts.DateFormat)
}

func (dts *DateTimeFormatter) FormatToLocalTime() string {
	return dts.localDateTime.Format(dts.TimeFormat)
}

func (dts *DateTimeFormatter) FormatToLocalDateTime() string {
	return dts.localDateTime.Format(dts.DateTimeFormat)
}
