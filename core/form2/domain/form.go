package domain

import (
	"net/url"
)

type (
	Form struct {
		//Data  the form Data Struct (Forms DTO)
		Data interface{}
		//ValidationInfo for the form
		ValidationInfo ValidationInfo
		//IsSubmitted  flag if form was submitted and this is the result page
		IsSubmitted bool
		//IsSuccess  - if IsValid && IsSubmitted && The planned Action was sucessfull (e.g. register user)
		IsSuccess bool
		//OriginalPostValues contain the original posted values
		OriginalPostValues url.Values
		//ValidationRules contains map with validation rules for all validatable fields
		ValidationRules map[string][]ValidationRule
	}
)
