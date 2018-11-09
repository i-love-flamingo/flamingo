# Form Package

This package provides helper to use forms in your interfaces.

## Usage

### Define a Datatransfer object (DTO) for your form:

Add your Data Representation of your form to your package ("/interfaces/controller/form")

  Example *(Example contains already annotations for the suggested libs - see below)*
```go
  package form
  
  type (
    AddressFormData struct {
      RegionCode   string `form:"regionCode" validate:"required" conform:"name"`
      CountryCode  string `form:"countryCode" validate:"required" conform:"name"`
      Company      string `form:"company" validate:"required" conform:"name"`
      Street       string `form:"street" validate:"required" conform:"name"`
      StreetNr     string `form:"streetNr" validate:"required" conform:"name"`
      AddressLine1 string `form:"addressLine1" validate:"required" conform:"name"`
      AddressLine2 string `form:"addressLine2" validate:"required" conform:"name"`
      Telephone    string `form:"telephone" validate:"required" conform:"name"`
      PostCode     string `form:"postCode" validate:"required" conform:"name"`
      City         string `form:"city" validate:"required" conform:"name"`
      Firstname    string `form:"firstname" validate:"required" conform:"name"`
      Lastname     string `form:"lastname" validate:"required" conform:"name"`
      Email        string `form:"email" validate:"required" conform:"name"`
    }
  )
```

* To process your form use "SimpleProcessFormRequest" or "ProcessFormRequest"

### Implementing a FormService

* Write an implementation of the interface domain.FormService.
  This interface describes two methods:
  
  **ParseFormData**
  We recommend to use the following packages:
   * "github.com/leebenson/conform": Helper to modify/sanitize strings
   * "github.com/go-playground/form": Helper to parse urlValues to struct (usage similar to the "json" struct annotation)
  
  **ValidateFormData**
  We recommend to use the following packages:
  
    * package "gopkg.in/go-playground/validator.v9": Helper to validate values in a struct.      
      This "form" package contains a service func "ValidationErrorsToValidationInfo" to use the results of this package.

  Example:
```go
  import (
    "net/url"
  
    "errors"
  
    formlib "github.com/go-playground/form"
    "github.com/leebenson/conform"
    "go.aoe.com/flamingo/core/form/application"
    "go.aoe.com/flamingo/core/form/domain"
    "go.aoe.com/flamingo/framework/web"
    "gopkg.in/go-playground/validator.v9"
  )
  
  type AddressFormService struct{}
  
  // use a single instance of Decoder, it caches struct info
  var decoder *formlib.Decoder
  
  // ParseForm - from FormService interface
  func (form *AddressFormService) ParseFormData(formValues url.Values, ctx web.Context) (interface{}, error) {
    decoder = formlib.NewDecoder()
    var formData AddressFormData
    decoder.Decode(&formData, formValues)
    conform.Strings(&formData)
    return formData, nil
  }
  
  //RunValidation - from FormService interface
  func (form *AddressFormService) ValidateFormData(data interface{}) (*domain.ValidationInfo, error) {
    if formData, ok := data.(AddressFormData); ok {
      validate := validator.New()
      return &application.ValidationErrorsToValidationInfo(validate.Struct(formData)), nil
    } else {
      return nil, errors.New("Cannot convert to AddressFormData")
    }
  }
```
    
  * Use the form in your controller Action:
  
```go
    form, e := formApplicationService.ProcessFormRequest(ctx, new(form.AddressFormService))
    // return on parse error (template need to handle error display)
    if e != nil {
      return cc.Render(ctx, "checkout/checkout", CheckoutViewData{
        Form: form,
      })
    }
    if form.IsValidAndSubmitted() {
      if addressData, ok := form.Data.(form.AddressFormData); ok {
        log.Printf("Do something with your data: %v", addressData)
        return cc.Redirect("checkout.success", nil).With("checkout.success.orderid", "orderid")
      }      
    }
``` 

 * Optional your implementation can also implement the interface "domain.GetDefaultFormData", to be able to prepopulate your form data

# Form validation

Form module gives different ways to attach custom validation into validator.Validate.
To use validator with injected custom validation, just define dingo provider which should
returns instance of validator.Validate and later call it when it's needed.

```go
type (
  validatorProvider func() *validator.Validate
	
  FormService struct {
    validatorProvider     validatorProvider
  }
)

func (fs *FormService) Inject(vp validatorProvider) {
  fs.validatorProvider = vp
}

func (fs *FormService) ValidateFormData(data interface{}) (formDomain.ValidationInfo, error) {
  validate := fs.validatorProvider()
  validationInfo := application.ValidationErrorsToValidationInfo(validate.Struct(formData))
  ...
}
```

## Date field validators

By using Validator Provider, date field validator are automatically injected so they can be
used in the FormData as presented in example:

```go
type FormData struct {
  ...
  DateOfBirth string `form:"dateOfBirth" validate:"required,dateformat,minimumage=18,maximumage=150"`
  ...
}
```

Date format can be changed as part of configuration (default value is "2006-01-02"):

```
form:
  validator:
    dateFormat: 02.01.2006
```

## Custom regex field validators

By using Validator Provider, it's possible to inject simple regex validators just by adapting
configuration:

```
form:
  validator:
    customRegex:
      password: ^[a-z]*$
```

By defining custom regex validator in configuration, it's further possible to use it as field validator,
with same name as provided in configuration ("password"):

```go
type FormData struct {
  ...
  Password 	string `form:"password" validate:"required,password"`
  ...
}
```

## Complex custom field validators

To inject complex field validators it's required to implement domain.FieldValidator interface and at least one
of the following domain.FieldValidatorWithParam or domain.FieldValidatorWithoutParam interfaces.
Following example shows custom field validator which implements both of the interfaces, but it's also enough
to implement only one of them.

```go
type (
  CustomMinValidator struct {
  	logger flamingo.Logger
  }
)

func (*CustomMinValidator) ValidatorName() string {
  return "custommin"
}

func (v *CustomMinValidator) ValidateWithoutParam(value interface{}) bool {
  number, ok := value.(int)
  if !ok {
    v.logger.WithField("customValidatorValue", fmt.Sprintf("%v", value))
    return false
  }
  return number > 0
}

func (v *CustomMinValidator) ValidateWithParam(param string value interface{}) bool {
  min, err := strconv.Atoi(param)
  if err != nil {
  	panic(err)
  }
  number, ok := value.(int)
  if !ok {
    v.logger.WithField("customValidatorValue", fmt.Sprintf("%v", value))
    return false
  }
  return number > min
}
```

To attach custom validator, simply inject it by using dingo injector:

```go
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*domain.FieldValidator)(nil)).To(&CustomMinValidator{})
}
```

Final usage of this custom field validator is:
```go
type FormData struct {
  ...
  SomeField       string `form:"someField" validate:"custommin"` // calls method ValidateWithoutParam
  SomeOtherField  string `form:"someOtherField" validate:"custommin=10"` // calls method ValidateWithParam, where param is 10
  ...
}
```
