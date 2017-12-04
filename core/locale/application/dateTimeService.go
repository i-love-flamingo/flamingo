package application

import (
	"time"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/locale/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	DateTimeService struct {
		DateFormat     string          `inject:"config:locale.dateTime.dateFormat"`
		TimeFormat     string          `inject:"config:locale.dateTime.timeFormat"`
		DateTimeFormat string          `inject:"config:locale.dateTime.dateTimeFormat"`
		Location       string          `inject:"config:locale.dateTime.location"`
		Logger         flamingo.Logger `inject:""`
	}
)

//GetDateTimeFromString Need string in format ISO: "2017-11-25T06:30:00Z"
func (dts *DateTimeService) GetDateTimeFromIsoString(dateTimeString string) (*domain.DateTime, error) {
	//"scheduledDateTime": "2017-11-25T06:30:00Z",
	timeResult, e := time.Parse(time.RFC3339, dateTimeString) //"2006-01-02T15:04:05Z"
	if e != nil {
		return nil, errors.Errorf("could not parse date in defined format: %v / Error: %v", dateTimeString, e)
	}
	dateTime := domain.DateTime{
		DateFormat:     dts.DateFormat,
		TimeFormat:     dts.TimeFormat,
		DateTimeFormat: dts.DateTimeFormat,
		Location:       dts.Location,
		Logger:         dts.Logger,
	}
	dateTime.SetDateTime(timeResult)
	return &dateTime, nil
}
