package domain

import (
	"regexp"
	"time"
)

//ValidateDate  validates a string representing a date from a input field type date (see https://www.w3.org/TR/2011/WD-html-markup-20110405/input.date.html)
func ValidateDate(date string) bool {
	if matched, _ := regexp.MatchString(`\d{4}-(0[1-9]|1[012])-(0[1-9]|1[0-9]|2[0-9]|3[0-1])`, date); matched {
		return true
	}
	return false
}

//ValidateAge - validates a date and checks if the date is older than given age - used for birthday validations for example
func ValidateAge(date string, age int) bool {
	if !ValidateDate(date) {
		return false
	}
	timev, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	now := time.Now()
	years := now.Year() - timev.Year()
	if now.Month() < timev.Month() || (now.Month() == timev.Month() && now.Day() < timev.Day()) {
		years--
	}
	if years < age {
		return false
	}
	return true
}
