package domain

//ValidateDate  validates a string representing a date from a input field type date (see https://www.w3.org/TR/2011/WD-html-markup-20110405/input.date.html)
// Deprecated
func ValidateDate(date string) bool {
	return validateDateFormat(date, "2006-01-02")
}

//ValidateAge - validates a date and checks if the date is older than given age - used for birthday validations for example
// Deprecated
func ValidateAge(date string, age int) bool {
	return validateMinimumAge(date, "2006-01-02", age)
}
