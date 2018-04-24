package domain

import (
	"regexp"
	"time"
)

func ValidateDate(date string) bool {
	if matched, _ := regexp.MatchString(`\d{4}-\d{2}-\d{2}`, date); matched {
		return true
	}
	return false
}

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
