# Form Package

This package provides helper to use forms in your interfaces layer.

## Usage

Form package provides out-of-the-box support for forms handling. Simply inject FormFactory
into your Controller and you can imidiatelly process any POST request and get values from
http request body.

```go
  package controller
  
  import (
    "context"
  
    "flamingo.me/flamingo/core/form2/application"
    "flamingo.me/flamingo/framework/web"
  )
  
  type MyController struct {
    formHandlerFactory application.FormHandlerFactory
  }
  
  func (c *MyController) Inject(f application.FormHandlerFactory) {
    .formHandlerFactory = f
  }
  
  func (c *MyController) Get(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateSimpleFormHandler()
    // HandleUnsubmittedForm provides default domain.Form instance without performing 
    // http request body processing and form data validation
    form, err := form.Handler.HandleUnsubmittedForm(ctx, req)
    
    // some code
  }
  
  func (c *MyController) Submit(ctx context.Context, req *web.Request) web.Response {
    // some code
      
    formHandler := c.formHandlerFactory.CreateSimpleFormHandler()
    // HandleSubmittedForm provides domain.Form instance after performing 
    // http request body processing and form data validation
    form, err := form.Handler.HandleSubmittedForm(ctx, req)
      
    // some code
  }
    
  func (c *MyController) Action(ctx context.Context, req *web.Request) web.Response {
    // some code
      
    formHandler := c.formHandlerFactory.CreateSimpleFormHandler()
    // HandleForm provides default domain.Form if it's GET http request.
    // If it's POST http request it provides domain.Form instance after performing 
    // http request body processing and form data validation
    form, err := form.Handler.HandleForm(ctx, req)
      
    // some code
  }  

```

Result from "Handle" methods are instances of domain.Form. It contains Data field which
contains parsed http request body represented as instance of map\[string\]string. In case
when form is not submitted and GET http request is processed, default form data is instance of
empty map (without any keys and values).

### Custom Form Data types

It's possible to provide specific custom form data. To do that, first specify data type:

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

Next step is to define a service which implement domain.FormDataProvider interface, which only
includes method "GetFormData":

```go
  type (
    AddressFormDataProvider struct {}
  }
  
  func (p *AddressFormDataProvider) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
    // define address form data with some default values
    return AddressFormData{
      Company:    "Flamingo crew"
      Firstname:  "Awesome",
      Lastname:   "Developer"
    }, nil
  }
```

Final part is to use this new AddressFormDataProvider during creation of domain.FormHandler instance:

```go
  func (c *MyController) Get(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateFormHandlerWithFormService(c.addressFormDataProvider)
    
    // some code
  }  

```

Form service overrides default form processing functionality. There is three functionality which can be
overridden:
* domain.FormDataProvider interface with method GetFormData - it gives support to provide custom form data
* domain.FormDataDecoder interface with method Decode - it gives support for custom form data decoding
* domain.FormDataValidator interface with method Validate - it gives support for custom form data validation

Form service can implement all of those three interfaces, but it must implement at least one of them.
If custom data provider, custom data decoder and custom data validator are defined in different form services,
they can be separately provided to FormHandlerFactory instance. If value NIL is passed for ony of those services
default domain.FormHandler behavior is used.

```go
  func (c *MyController) Get(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateFormHandlerWithFormServices(
      c.addressFormDataProvider,
      c.addressFormDataDecoder,
      c.addressFormDataValidator,
    )
    
    // some code
  }  

```

It's also possible to use FormBuilder instance for more specific domain.FormHandler customization:

```go
  func (c *MyController) Get(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.GetBuilder().
      SetFormDataProvider(c.addressFormDataProvider).
      SetFormDataValidator(c.addressFormDataValidator).
      Build()
    
    // some code
  }  

```

### Custom Form Data decoding

Default domain.FormDataDecoder provides http request body decoding provided by "github.com/go-playground/form"
and string processing provided by "github.com/leebenson/conform". To provide custom form data provider, simply
define new form service type:

```go
  type (
    AddressFormDataDecoder struct {}
  }
  
  func (p *AddressFormDataDecoder) Decode(ctx context.Context, req *web.Request, values url.Values, formData interface{}) (interface{}, error) {
    // some code
  }
```

Finally, it can be provided to instance of domain.FormHandler by using FormHandlerFactory or FormHandlerBuilder:

```go
  func (c *MyController) First(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateFormHandlerWithFormService(c.addressFormDataDecoder)
    
    // some code
  }  
  
  func (c *MyController) Second(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.GetBuilder().
      SetFormDataDecoderr(c.addressFormDataDecoder).
      Build()
    
    // some code
  }  

```

### Custom Form Data validation

Default domain.FormDataValidator provides full struct validation via github.com/go-playground/validator". 
To provide custom form data validator, simply define new form service type:

```go
  type (
    AddressFormDataValidator struct {}
  }
  
  func (p *AddressFormDataValidator) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
    // some code
  }
```

Instance of domain.ValidatorProvider contains extended validator with all injected custom field and struct validators
(check documentation bellow to see how to extend validators).

Finally, it can be provided to instance of domain.FormHandler by using FormHandlerFactory or FormHandlerBuilder:

```go
  func (c *MyController) First(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateFormHandlerWithFormService(c.addressFormDataValidator)
    
    // some code
  }  
  
  func (c *MyController) Second(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.GetBuilder().
      SetFormDataValidator(c.addressFormDataValidator).
      Build()
    
    // some code
  }  

```

### Named form services

Beside defining form services as pure instance by using FormHandlerFactory or FormHandlerBuilder,
all form services can be defined globally, with dingo injector, and later user by their name.
That means, they need to be injected as specific instances of interfaces:
* domain.FormServiceWithName
* domain.FormDataProviderWithName
* domain.FormDataDecoderWithName
* domain.FormDataValidatorWithName

Beside regular methods, those instances must implement method "Name", which provides name of form service
as a string:

```go
  type (
    MyGlobalFormDataProvider struct {}
  }
  
  func (p *MyGlobalFormDataProvider) Name() {
    return "formDataProvider.myCustom"
  }
  
  func (p *MyGlobalFormDataProvider) GetFormData(ctx context.Context, req *web.Request) (interface{}, error) {
    // define address form data with some default values
  }
```

Than, service should be injected via dingo injector:

```go
  func (m *Module) Configure(injector *dingo.Injector) {
    // some code
  
  	injector.BindMulti(new(domain.FormDataProviderWithName)).To(providers.MyGlobalFormDataProvider{})
  	
  	// some code
  }
```

Finally, provider can be defined for instance domain.FormHandler by using FormHandlerBuilder:

```go
  func (c *MyController) Second(ctx context.Context, req *web.Request) web.Response {
      // some code
      
      formHandler := c.formHandlerFactory.GetBuilder().
        SetNamedFormDataProvider("formDataProvider.myCustom").
        Build()
      
      // some code
    } 
```

### Form extensions

Form extensions are smaller form services which can be used with multiple forms. They perform side jobs which is
not reflected as final form data, but it can affect validation results. All form extension can implements at least one
of mentioned interfaces for form services:
* domain.FormDataProvider
* domain.FormDataDecoder
* domain.FormDataValidator

To add some form extensions into domain.FormHandlerInstance there are multiple ways via FormHandlerFactory or FormHandlerBuilder:

```go
  func (c *MyController) First(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.CreateFormHandlerWithFormService(c.formService, c.formExtension, "formExtension.csrfToken")
    
    // some code
  }  
  
  func (c *MyController) Second(ctx context.Context, req *web.Request) web.Response {
    // some code
    
    formHandler := c.formHandlerFactory.GetBuilder().
      SetFormDataDecoder(c.addressFormDataDecoder).
      AddFormExtension(c.formExtension).
      AddNamedFormExtension("formExtension.csrfToken").
      Build()
    
    // some code
  }  

```

To provide custom form extension with specific name which can be used globally, simply define form service with "Name" method:

```go
  type (
    MyGlobalFormExtension struct {}
  }
  
  func (p *MyGlobalFormExtension) Name() {
    return "formExtension.myCustom"
  }
  
  func (p *MyGlobalFormExtension) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
    // define address form data with some default values
  }
```

Finally, define form extension via dingo injector:

```go
  func (m *Module) Configure(injector *dingo.Injector) {
    // some code
  
  	injector.BindMulti(new(domain.FormExtensionWithName)).To(extensions.MyGlobalFormExtension{})
  	
  	// some code
  }
```

# Validation Provider

Form module gives different ways to attach custom validators into validator.Validate instance
by using domain.ValidationProvider instance. Default functionality for instance of domain.FormHandler is
to use all kind of core validator.Validate validations, and in addition to use all domain.FieldValidators
and domain.StructValidators instances. Still, if it's needed to perform different type of validation,
it's possible to provide form service, which implements domain.FormDataValidator interface:

```go
  type (
    AddressFormDataValidator struct {}
  }
  
  func (p *AddressFormDataValidator) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
    // some code
  }

## Additional validators
### Date field validators

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

### Custom regex field validators

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

### Complex custom field validators

To inject complex field validators it's required to implement domain.FieldValidator:

```go
type (
  CustomMinValidator struct {
  	logger flamingo.Logger
  }
)

func (*CustomMinValidator) ValidatorName() string {
  return "custommin"
}

func (v *CustomMinValidator) ValidateField(ctx context.Context, fl validator.FieldLevel) bool {
  // some code
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

### Complex custom struct validators

To inject struct field validators it's required to implement domain.StructValidator:

```go
type (
  FormDataValidator struct {
  	logger flamingo.Logger
  }
)

func (*FormDataValidator) StructType() interface{} {
  return FormData{}
}

func (v *FormDataValidator) ValidateStruct(ctx context.Context, sl validator.StructLevel) {
  // some code
}
```

To attach custom validator, simply inject it by using dingo injector:

```go
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*domain.StructValidator)(nil)).To(&FormDataValidator{})
}
```

Notice: Instance of validator.Validate allows only one struct validation per type. This means
that only last defined struct validator for single type will be use in struct validation.
