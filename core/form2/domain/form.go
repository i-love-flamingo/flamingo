package domain

type (
	Form struct {
		//Data  the form Data Struct (Forms DTO)
		Data interface{}
		//ValidationInfo for the form
		ValidationInfo ValidationInfo
		//submitted  flag if form was submitted and this is the result page
		submitted bool
		//validationRules contains map with validation rules for all validatable fields
		validationRules map[string][]ValidationRule
	}
)

func NewForm(submitted bool, validationRules map[string][]ValidationRule) Form {
	return Form{
		submitted:       submitted,
		validationRules: validationRules,
	}
}

func (f Form) IsValidAndSubmitted() bool {
	return f.IsValid() && f.IsSubmitted()
}

func (f Form) IsValid() bool {
	return f.ValidationInfo.IsValid()
}

func (f Form) IsSubmitted() bool {
	return f.submitted
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
	return f.ValidationInfo.GetErrorsForField(name)
}

//GetValidationRulesForField adds option to extract validation rules for desired field in templates
func (f Form) GetValidationRulesForField(name string) []ValidationRule {
	return f.validationRules[name]
}
