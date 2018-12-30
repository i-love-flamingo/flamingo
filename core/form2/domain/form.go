package domain

type (
	Form struct {
		//Data  the form Data Struct (Forms DTO)
		Data interface{}
		//ValidationInfo for the form
		ValidationInfo ValidationInfo
		//IsSubmitted  flag if form was submitted and this is the result page
		IsSubmitted bool
		//IsValid flag if form was validated successfully
		IsValid bool
		//ValidationRules contains map with validation rules for all validatable fields
		ValidationRules map[string][]ValidationRule
	}
)
