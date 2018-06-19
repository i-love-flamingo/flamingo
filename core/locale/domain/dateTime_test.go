package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type formatTest struct {
	timeStamp         string
	location          string
	expectedDate      string
	expectedLocalDate string
	expectedTime      string
	expectedLocalTime string
}

var formatTestData = []formatTest{
	// dst winter
	{"2018-01-01T01:00:00Z", "UTC", "01 Jan 2018", "01 Jan 2018", "01:00", "01:00"},
	{"2018-01-01T01:00:00Z", "America/New_York", "01 Jan 2018", "31 Dec 2017", "01:00", "20:00"},
	{"2018-01-01T01:00:00Z", "Europe/Berlin", "01 Jan 2018", "01 Jan 2018", "01:00", "02:00"},
	{"2018-01-01T01:00:00Z", "Europe/London", "01 Jan 2018", "01 Jan 2018", "01:00", "01:00"},

	// dst summer
	{"2018-06-01T01:00:00Z", "UTC", "01 Jun 2018", "01 Jun 2018", "01:00", "01:00"},
	{"2018-06-01T01:00:00Z", "America/New_York", "01 Jun 2018", "31 May 2018", "01:00", "21:00"},
	{"2018-06-01T01:00:00Z", "Europe/Berlin", "01 Jun 2018", "01 Jun 2018", "01:00", "03:00"},
	{"2018-06-01T01:00:00Z", "Europe/London", "01 Jun 2018", "01 Jun 2018", "01:00", "02:00"},
}

func TestFormat(t *testing.T) {
	for _, testData := range formatTestData {
		f := testGetFormatter(testData.timeStamp, testData.location)

		assert.Equal(
			t,
			testData.expectedDate,
			f.FormatDate(),
			fmt.Sprintf("Date of %v in %v", testData.timeStamp, testData.location),
		)
		assert.Equal(
			t,
			testData.expectedLocalDate,
			f.FormatToLocalDate(),
			fmt.Sprintf("Local date of %v in %v", testData.timeStamp, testData.location),
		)
		assert.Equal(
			t,
			testData.expectedTime,
			f.FormatTime(),
			fmt.Sprintf("Time of %v in %v", testData.timeStamp, testData.location),
		)
		assert.Equal(
			t,
			testData.expectedLocalTime,
			f.FormatToLocalTime(),
			fmt.Sprintf("Local time of %v in %v", testData.timeStamp, testData.location),
		)

	}
}

func testGetFormatter(timeString string, locationString string) *DateTimeFormatter {
	f := DateTimeFormatter{
		DateFormat:     "02 Jan 2006",
		TimeFormat:     "15:04",
		DateTimeFormat: "02 Jan 2006 15:04:05",
	}

	dateTime, e := time.Parse(time.RFC3339, timeString)
	if e != nil {
		// just panic here - that is enough for the moment
		panic(e)
	}
	location, e := time.LoadLocation(locationString)
	localTime := dateTime.In(location)

	f.SetDateTime(dateTime, localTime)

	return &f
}
