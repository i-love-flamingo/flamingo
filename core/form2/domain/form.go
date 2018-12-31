package domain

type (
	Form struct {
		//Data  the form Data Struct (Forms DTO)
		Data interface{}
		//ValidationInfo for the form
		ValidationInfo ValidationInfo
		//IsSubmitted  flag if form was submitted and this is the result page
		IsSubmitted bool
		//ValidationRules contains map with validation rules for all validatable fields
		ValidationRules map[string][]ValidationRule
	}
)

func (f Form) IsValidAndSubmitted() bool {
	return f.ValidationInfo.IsValid() && f.IsSubmitted
}

func (f Form) HasErrorForField(name string) bool {
	return f.ValidationInfo.HasErrorsForField(name)
}

func (f Form) HasAnyFieldErrors() bool {
	return f.ValidationInfo.HasAnyFieldErrors()
}

func (f Form) HasGeneralErrors() bool {
	return f.ValidationInfo.HasGeneralErrors()
}

func (f Form) GetErrorsForField(name string) []Error {
	return f.ValidationInfo.GetFieldErrors(name)
}

//GetValidationRulesForField adds option to extract validation rules for desired field in templates
func (f Form) GetValidationRulesForField(name string) []ValidationRule {
	return f.ValidationRules[name]
}
