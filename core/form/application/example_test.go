package application_test

import (
	"fmt"

	"context"
	"net/http"

	"strings"

	"flamingo.me/flamingo/core/form/application"
	"flamingo.me/flamingo/framework/web"
	"gopkg.in/go-playground/validator.v9"
)

func ExampleValidationErrorsToValidationInfo() {
	type (
		CustomerEditFormData struct {
			FirstName string `validate:"required"`
			LastName  string `validate:"required"`
			Title     string ``
		}
	)
	formData := CustomerEditFormData{}
	//validate - result from package validator "gopkg.in/go-playground/validator.v9"
	validate := validator.New()
	result := application.ValidationErrorsToValidationInfo(validate.Struct(formData))

	fmt.Printf("%v\n", result.IsValid)
	fmt.Printf(result.FieldErrors["firstName"][0].MessageKey)

	// Output: false
	// formerror_firstName_required
}

func ExampleSimpleProcessFormRequest() {

	httpRequest, _ := http.NewRequest("POST", "?test=demo", strings.NewReader(""))

	flamingoWebRequest := web.RequestFromRequest(httpRequest, nil)
	form, _ := application.SimpleProcessFormRequest(context.Background(), flamingoWebRequest)

	fmt.Printf("%v\n", form.IsSubmitted)
	fmt.Print(form.Data.(map[string]string)["test"])

	// Output: true
	// demo
}
