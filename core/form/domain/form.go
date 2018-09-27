package domain

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/framework/web"
)

type (
	//Form represents a Form - its intended usage is to pass to your view
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

	//ValidationInfo - represents the complete Validation Informations of your form. It can contain GeneralErrors and form field related errors.
	ValidationInfo struct {
		//FieldErrors list of errors per form field.
		FieldErrors map[string][]Error
		//GeneralErrors list of general form errors, that are not related to any field
		GeneralErrors []Error
		//IsValid  flag if data was valid
		IsValid bool
	}

	//ValidationRule - contains single validation rule for field. Name is mandatory (required|email|max|len|...), Value is optional and adds additional info (like "128" for "max=128" rule)
	ValidationRule struct {
		Name  string
		Value string
	}

	//Error - representation of an Error Message - intented usage is to display errors in the view to the end user
	Error struct {
		//Tag - contains the validation tag that failed. if the
		Tag string
		//MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		//DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
	}

	// FormService interface that need to be implemented, in case you want to use "application.ProcessFormRequest"
	FormService interface {
		//ParseFormData is responsible of mapping the passed formValues to your FormData Struct (Forms DTO)
		ParseFormData(ctx context.Context, r *web.Request, formValues url.Values) (interface{}, error)
	}

	// ValidateFormData interface that need to be implemented, in case you want to use "application.ProcessFormRequest" with validation that doesn't include context
	ValidateFormData interface {
		//ValidateFormData is responsible to run validations on the Data, the returned error type can be a slice of errors. each error is converted to a validation Error
		ValidateFormData(data interface{}) (ValidationInfo, error)
	}

	// ValidateFormDataWithContext interface that need to be implemented, in case you want to use "application.ProcessFormRequest" with validation that include context
	ValidateFormDataWithContext interface {
		//ValidateFormDataWithContext is responsible to run validations on the Data, with provided context, the returned error type can be a slice of errors. each error is converted to a validation Error
		ValidateFormDataWithContext(ctx context.Context, data interface{}) (ValidationInfo, error)
	}

	// GetDefaultFormData interface
	GetDefaultFormData interface {
		//GetDefaultFormData
		GetDefaultFormData(parsedData interface{}) interface{}
	}

	// GetDefaultFormDataWithContext interface
	GetDefaultFormDataWithContext interface {
		//GetDefaultFormDataWithContext
		GetDefaultFormDataWithContext(ctx context.Context, parsedData interface{}) interface{}
	}
)

//AddGeneralUnknownError - adds a general unknown error to the validation infos
func (vi *ValidationInfo) AddGeneralUnknownError(err error) {
	vi.GeneralErrors = append(vi.GeneralErrors, Error{MessageKey: "unknown_error", DefaultLabel: "An error occured!"})
	vi.IsValid = false
}

//AddError adds a general error with the passed MessageKey and DefaultLabel
func (vi *ValidationInfo) AddError(messageKey string, defaultLabel string) {
	vi.GeneralErrors = append(vi.GeneralErrors, Error{MessageKey: messageKey, DefaultLabel: defaultLabel})
	vi.IsValid = false
}

//AddFieldError -adds a Error with the passed messageKey and defaultLabel - for a (form)field with the given name
func (vi *ValidationInfo) AddFieldError(fieldName string, messageKey string, defaultLabel string) {
	if vi.FieldErrors == nil {
		vi.FieldErrors = make(map[string][]Error)
	}
	if _, ok := vi.FieldErrors[fieldName]; !ok {
		vi.FieldErrors[fieldName] = make([]Error, 0)
	}
	vi.FieldErrors[fieldName] = append(vi.FieldErrors[fieldName], Error{MessageKey: messageKey, DefaultLabel: defaultLabel})
	vi.IsValid = false
}

func (f Form) IsValidAndSubmitted() bool {
	return f.ValidationInfo.IsValid && f.IsSubmitted
}

func (f Form) HasErrorForField(name string) bool {
	if _, ok := f.ValidationInfo.FieldErrors[name]; ok {
		return true
	}
	return false
}

func (f Form) HasAnyFieldErrors() bool {
	return len(f.ValidationInfo.FieldErrors) > 0
}

func (f Form) HasGeneralErrors() bool {
	return len(f.ValidationInfo.GeneralErrors) > 0
}

func (f Form) GetErrorsForField(name string) []Error {
	if v, ok := f.ValidationInfo.FieldErrors[name]; ok {
		return v
	}
	return nil
}

func (f Form) GetOriginalPostValue1(key string) string {
	if f.OriginalPostValues == nil {
		return ""
	}
	values := f.OriginalPostValues[key]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

//GetValidationRulesForField adds option to extract validation rules for desired field in templates
func (f Form) GetValidationRulesForField(name string) []ValidationRule {
	return f.ValidationRules[name]
}
