package domain

import "time"

// DateTimeFormatter has a couple of helpful methods to format date and times
type DateTimeFormatter struct {
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
