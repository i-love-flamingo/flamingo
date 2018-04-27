package domain

import (
	"net/url"

	"go.aoe.com/flamingo/framework/web"
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
	}

	ValidationInfo struct {
		FieldErrors   map[string][]Error
		GeneralErrors []Error
		//IsValid  flag if data was valid
		IsValid bool
	}

	Error struct {
		Tag          string
		MessageKey   string
		DefaultLabel string
	}

	// FormService interface
	FormService interface {
		//ParseFormData is responsible of mapping the passed formValues to your FormData Struct
		ParseFormData(ctx web.Context, formValues url.Values) (interface{}, error)
		//ValidateFormData is responsible to run validations on the Data, the returned error type can be a slice of errors. each error is converted to a validation Error
		ValidateFormData(data interface{}) (ValidationInfo, error)
	}
	// GetDefaultFormData interface
	GetDefaultFormData interface {
		//GetDefaultFormData
		GetDefaultFormData(parsedData interface{}) interface{}
	}
)

func (vi *ValidationInfo) AddGeneralUnknownError(err error) {
	vi.GeneralErrors = append(vi.GeneralErrors, Error{MessageKey: "unknown_error", DefaultLabel: "An error occured!"})
	vi.IsValid = false
}

func (vi *ValidationInfo) AddError(key string, defaultLabel string) {
	vi.GeneralErrors = append(vi.GeneralErrors, Error{MessageKey: key, DefaultLabel: defaultLabel})
	vi.IsValid = false
}

func (vi *ValidationInfo) AddFieldError(fieldName string, key string, defaultLabel string) {
	if vi.FieldErrors == nil {
		vi.FieldErrors = make(map[string][]Error)
	}
	if _, ok := vi.FieldErrors[fieldName]; !ok {
		vi.FieldErrors[fieldName] = make([]Error, 0)
	}
	vi.FieldErrors[fieldName] = append(vi.FieldErrors[fieldName], Error{MessageKey: key, DefaultLabel: defaultLabel})
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
