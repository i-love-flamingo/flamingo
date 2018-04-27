package application

import (
	"errors"
	"log"
	"net/url"

	"strings"

	"go.aoe.com/flamingo/core/form/domain"
	"go.aoe.com/flamingo/framework/web"
	"gopkg.in/go-playground/validator.v9"
)

//ProcessFormRequest: Parses and Validates a Request to a Form - with the Help of the passed FormService
func ProcessFormRequest(ctx web.Context, service domain.FormService) (domain.Form, error) {
	form := domain.Form{}

	urlValues, err := getPostValues(ctx)
	if err != nil {
		form.ValidationInfo.AddGeneralUnknownError(err)
		return form, err
	}
	form.OriginalPostValues = urlValues

	form.Data, err = parseFormData(urlValues, service, ctx)
	if err != nil {
		form.ValidationInfo.AddGeneralUnknownError(err)
		return form, err
	}

	//Run Validation only if form was submitted
	if urlValues.Get("novalidate") != "true" && ctx.Request().Method == "POST" {
		form.IsSubmitted = true
		form.ValidationInfo, err = service.ValidateFormData(form.Data)
		if err != nil {
			form.ValidationInfo = ValidationErrorsToValidationInfo(err)
		}
	} else {
		if defaultFormDataService, ok := service.(domain.GetDefaultFormData); ok {
			form.Data = defaultFormDataService.GetDefaultFormData(form.Data)
			log.Printf("############ %v", form.Data)
		}
	}

	return form, nil
}

//SimpleProcessFormRequest: Parses Post Values and returns a simple map - can be used if you dont need/want to implement a domain.FormService
func SimpleProcessFormRequest(ctx web.Context) (domain.Form, error) {
	var err error
	var urlValues url.Values
	form := domain.Form{}

	if ctx.Request().Method != "POST" {
		form.IsSubmitted = false
		form.ValidationInfo.IsValid = true
		return form, nil
	}

	form.IsSubmitted = true

	urlValues, err = getPostValues(ctx)
	if err != nil {
		form.ValidationInfo.AddGeneralUnknownError(err)
		return form, err
	}
	form.ValidationInfo.IsValid = true
	dataMap := make(map[string]string)
	for k, v := range urlValues {
		dataMap[k] = strings.Join(v, " ")
	}
	form.ValidationInfo.IsValid = true
	form.Data = dataMap

	return form, nil
}

func ValidationErrorsToValidationInfo(err error) domain.ValidationInfo {
	var validationInfo domain.ValidationInfo

	validationInfo.IsValid = true
	validationInfo.FieldErrors = make(map[string][]domain.Error)

	if err == nil {
		return validationInfo
	}

	if err1, ok := err.(*validator.InvalidValidationError); ok {
		validationInfo.IsValid = false
		validationInfo.AddGeneralUnknownError(err1)
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			err.Tag()
			log.Printf("Error Form %#v", err)
			var errorValue domain.Error
			validationInfo.IsValid = false
			fieldName := getRelativeFieldNameFromValidationError(err)
			errorValue = domain.Error{
				Tag:          err.Tag(),
				MessageKey:   "formerror_" + fieldName + "_" + err.Tag(),
				DefaultLabel: err.Field() + " wrong",
			}
			validationInfo.FieldErrors[fieldName] = append(validationInfo.FieldErrors[fieldName], errorValue)
		}
	}

	return validationInfo
}

func getRelativeFieldNameFromValidationError(err validator.FieldError) string {
	var result []string
	fieldName := err.Namespace()
	//first part of namespace is not required to have the relative path:
	fieldName = fieldName[(strings.Index(fieldName, ".") + 1):]
	for _, part := range strings.Split(fieldName, ".") {
		result = append(result, strings.ToLower(part[0:1])+part[1:])
	}
	return strings.Join(result, ".")
}

func getPostValues(ctx web.Context) (url.Values, error) {
	err := ctx.Request().ParseForm()
	if err != nil {
		log.Printf("form.application: Parse Form Error %v", err)
		return ctx.Request().Form, errors.New("unkown_error")
	}
	return ctx.Request().Form, nil
}

func parseFormData(values url.Values, service domain.FormService, ctx web.Context) (interface{}, error) {
	formData, err := service.ParseFormData(ctx, values)
	if err != nil {
		log.Printf("form.application: ParseForm Error %v", err)
		return formData, errors.New("unkown_error")
	}
	return formData, nil
}
