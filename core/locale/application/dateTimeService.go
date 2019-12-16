package application

import (
	"time"

	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"
)

type (
	// DateTimeServiceInterface to define a service to obtain formatted date time data
	DateTimeServiceInterface interface {
		GetDateTimeFormatterFromIsoString(dateTimeString string) (*domain.DateTimeFormatter, error)
		GetDateTimeFormatter(dateTime time.Time) (*domain.DateTimeFormatter, error)
	}

	// DateTimeService is a basic support service for date/time parsing
	DateTimeService struct {
		dateFormat     string
		timeFormat     string
		dateTimeFormat string
		location       string
		logger         flamingo.Logger
	}
)

// check interface implementation
var _ DateTimeServiceInterface = (*DateTimeService)(nil)

// Inject dependencies
func (dts *DateTimeService) Inject(
	logger flamingo.Logger,
	config *struct {
		DateFormat     string `inject:"config:core.locale.date.dateFormat"`
		TimeFormat     string `inject:"config:core.locale.date.timeFormat"`
		DateTimeFormat string `inject:"config:core.locale.date.dateTimeFormat"`
		Location       string `inject:"config:core.locale.date.location"`
	},
) {
	dts.logger = logger
	dts.dateFormat = config.DateFormat
	dts.timeFormat = config.TimeFormat
	dts.dateTimeFormat = config.DateTimeFormat
	dts.location = config.Location
}

// GetDateTimeFormatterFromIsoString Need string in format ISO: "2017-11-25T06:30:00Z"
func (dts *DateTimeService) GetDateTimeFormatterFromIsoString(dateTimeString string) (*domain.DateTimeFormatter, error) {
	timeResult, err := time.Parse(time.RFC3339, dateTimeString) //"2006-01-02T15:04:05Z"
	if err != nil {
		return nil, errors.Errorf("could not parse date in defined format: %v / Error: %v", dateTimeString, err)
	}

	return dts.GetDateTimeFormatter(timeResult)
}

// GetDateTimeFormatter from time
func (dts *DateTimeService) GetDateTimeFormatter(timeValue time.Time) (*domain.DateTimeFormatter, error) {
	loc, err := dts.loadLocation()
	if err != nil {
		if dts.logger != nil {
			dts.logger.Warn("dateTime Parsing error - could not load location - use UTC as fallback", dts.location)
		}
		loc = time.UTC
	}

	dateTime := domain.DateTimeFormatter{
		DateFormat:     dts.dateFormat,
		TimeFormat:     dts.timeFormat,
		DateTimeFormat: dts.dateTimeFormat,
	}

	if dts.logger != nil {
		dateTime.SetLogger(dts.logger)
	}

	dateTime.SetDateTime(timeValue, timeValue.In(loc))

	return &dateTime, nil
}

func (dts *DateTimeService) loadLocation() (*time.Location, error) {
	if dts.location == "" {
		return nil, errors.Errorf("No location configured")
	}

	// try to load the configured location
	return time.LoadLocation(dts.location)
}
