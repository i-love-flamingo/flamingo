package domain

type (
	// ValidationInfo - represents the complete Validation Informations of your form. It can contain GeneralErrors and form field related errors.
	ValidationInfo struct {
		// fieldErrors list of errors per form field.
		fieldErrors map[string][]Error
		// generalErrors list of general form errors, that are not related to any field
		generalErrors []Error
	}

	// ValidationRule - contains single validation rule for field. Name is mandatory (required|email|max|len|...), Value is optional and adds additional info (like "128" for "max=128" rule)
	ValidationRule struct {
		// Name validator tag name
		Name string
		// Value additional parameter provided as condition for validation tag
		Value string
	}

	// Error - representation of an Error Message - intented usage is to display errors in the view to the end user
	Error struct {
		// MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		// DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
	}
)

// IsValid method which defines if validation info is related to valid data or not
func (vi *ValidationInfo) IsValid() bool {
	return !vi.HasGeneralErrors() && !vi.HasAnyFieldErrors()
}

// HasGeneralErrors method which defines if there is any general validations error
func (vi *ValidationInfo) HasGeneralErrors() bool {
	return len(vi.generalErrors) > 0
}

// HasAnyFieldErrors method which defines if there is any field validations error for any field
func (vi *ValidationInfo) HasAnyFieldErrors() bool {
	if vi.fieldErrors == nil {
		return false
	}

	for fieldName := range vi.fieldErrors {
		if vi.HasErrorsForField(fieldName) {
			return true
		}
	}
	return false
}

// HasErrorsForField method which defines if there is any field validations error for specific field
func (vi *ValidationInfo) HasErrorsForField(fieldName string) bool {
	return vi.fieldErrors != nil && len(vi.fieldErrors[fieldName]) > 0
}

// AppendGeneralErrors method which appends all provided validation errors to general errors, without duplicating existing ones
func (vi *ValidationInfo) AppendGeneralErrors(errs []Error) {
	for _, err := range errs {
		vi.AddGeneralError(err.MessageKey, err.DefaultLabel)
	}
}

// AddError method which adds a general error with the passed MessageKey and DefaultLabel
func (vi *ValidationInfo) AddGeneralError(messageKey string, defaultLabel string) {
	keys := vi.getExistingMessageKeys(vi.generalErrors)

	if keys[messageKey] {
		return
	}

	err := Error{
		MessageKey:   messageKey,
		DefaultLabel: defaultLabel,
	}

	vi.generalErrors = append(vi.generalErrors, err)
}

// GetGeneralErrors method which returns list of all general validation errors
func (vi *ValidationInfo) GetGeneralErrors() []Error {
	return vi.generalErrors
}

// AppendFieldErrors method which appends all provided validation errors to field errors, without duplicating existing ones
func (vi *ValidationInfo) AppendFieldErrors(fieldErrors map[string][]Error) {
	for fieldName, errs := range fieldErrors {
		for _, err := range errs {
			vi.AddFieldError(fieldName, err.MessageKey, err.DefaultLabel)
		}
	}
}

// AddFieldError method which adds a field error with the passed field name, message key and default label
func (vi *ValidationInfo) AddFieldError(fieldName string, messageKey string, defaultLabel string) {
	if vi.fieldErrors == nil {
		vi.fieldErrors = map[string][]Error{}
	}

	keys := vi.getExistingMessageKeys(vi.fieldErrors[fieldName])

	if keys[messageKey] {
		return
	}

	err := Error{
		MessageKey:   messageKey,
		DefaultLabel: defaultLabel,
	}

	vi.fieldErrors[fieldName] = append(vi.fieldErrors[fieldName], err)
}

// GetGeneralErrors method which returns list of all field validation errors for all fields
func (vi *ValidationInfo) GetErrorsForAllFields() map[string][]Error {
	return vi.fieldErrors
}

// GetFieldErrors method which returns list of all general validation errors for specific field
func (vi *ValidationInfo) GetErrorsForField(fieldName string) []Error {
	return vi.fieldErrors[fieldName]
}

// getExistingMessageKeys method which returns all message keys used in specific list of validation errors
func (vi *ValidationInfo) getExistingMessageKeys(errs []Error) map[string]bool {
	keys := make(map[string]bool, len(errs))
	for _, err := range errs {
		keys[err.MessageKey] = true
	}
	return keys
}
