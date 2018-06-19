package application

import (
	"time"

	"github.com/pkg/errors"
	"flamingo.me/flamingo/core/locale/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	DateTimeServiceInterface interface {
		GetDateTimeFormatterFromIsoString(dateTimeString string) (*domain.DateTimeFormatter, error)
		GetDateTimeFormatter(dateTime time.Time) (*domain.DateTimeFormatter, error)
	}

	// DateTimeService is a basic support service for date/time parsing
	DateTimeService struct {
		DateFormat     string          `inject:"config:locale.date.dateFormat"`
		TimeFormat     string          `inject:"config:locale.date.timeFormat"`
		DateTimeFormat string          `inject:"config:locale.date.dateTimeFormat"`
		Location       string          `inject:"config:locale.date.location"`
		Logger         flamingo.Logger `inject:""`
	}
)

// check interface implementation
var _ DateTimeServiceInterface = (*DateTimeService)(nil)

//GetDateTimeFromString Need string in format ISO: "2017-11-25T06:30:00Z"
func (dts *DateTimeService) GetDateTimeFormatterFromIsoString(dateTimeString string) (*domain.DateTimeFormatter, error) {
	timeResult, err := time.Parse(time.RFC3339, dateTimeString) //"2006-01-02T15:04:05Z"
	if err != nil {
		return nil, errors.Errorf("could not parse date in defined format: %v / Error: %v", dateTimeString, err)
	}

	return dts.GetDateTimeFormatter(timeResult)
}

//GetDateTimeFormatter from time
func (dts *DateTimeService) GetDateTimeFormatter(timeValue time.Time) (*domain.DateTimeFormatter, error) {
	loc, err := dts.loadLocation()
	if err != nil {
		if dts.Logger != nil {
			dts.Logger.Warn("dateTime Parsing error - could not load location - use UTC as fallback", dts.Location)
		}
		loc = time.UTC
	}

	dateTime := domain.DateTimeFormatter{
		DateFormat:     dts.DateFormat,
		TimeFormat:     dts.TimeFormat,
		DateTimeFormat: dts.DateTimeFormat,
	}
	dateTime.SetDateTime(timeValue, timeValue.In(loc))

	return &dateTime, nil
}

func (dts *DateTimeService) loadLocation() (*time.Location, error) {
	if dts.Location == "" {
		return nil, errors.Errorf("No location configured")
	}

	// try to load the configured location
	return time.LoadLocation(dts.Location)
}
