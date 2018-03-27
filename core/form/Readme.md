## Form Package

This package provides helper to use forms in your interfaces.

### Usage

* Add your Data Representation of your form to your package ("/interfaces/controller/form")

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
