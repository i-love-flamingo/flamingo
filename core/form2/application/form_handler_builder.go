package application

import (
	"fmt"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// FormHandlerBuilder as interface for complex creation of form handler instance
	FormHandlerBuilder interface {
		// SetFormService sets form service instance and overrides provider, decoder and validator if
		// it implements theirs interfaces. If it doesn't implements any of those interfaces it panics.
		SetFormService(formService interface{}) FormHandlerBuilder
		// SetNamedFormService sets form service instance by searching named form service provided via dingo injector.
		// It panics if there is no injected form service with that name.
		// It overrides provider, decoder and validator if it implements theirs interfaces.
		// If it doesn't implements any of those interfaces it panics.
		SetNamedFormService(name string) FormHandlerBuilder
		// SetFormDataProvider sets form data provider instance and overrides default one.
		SetFormDataProvider(formDataProvider domain.FormDataProvider) FormHandlerBuilder
		// SetNamedFormDataProvider sets form data provider by searching named provider provided via dingo injector.
		// It panics if there is no injected form data provider with that name.
		// It sets form data provider instance and overrides default one.
		SetNamedFormDataProvider(name string) FormHandlerBuilder
		// SetFormDataDecoder sets form data decoder instance and overrides default one.
		SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) FormHandlerBuilder
		// SetNamedFormDataDecoder sets form data decoder by searching named decoder provided via dingo injector.
		// It panics if there is no injected form data decoder with that name.
		// It sets form data decoder instance and overrides default one.
		SetNamedFormDataDecoder(name string) FormHandlerBuilder
		// SetFormDataValidator sets form data validator instance and overrides default one.
		SetFormDataValidator(formDataValidator domain.FormDataValidator) FormHandlerBuilder
		// SetNamedFormDataValidator sets form data decoder by searching named decoder validator via dingo injector.
		// It panics if there is no injected form data validator with that name.
		// It sets form data validator instance and overrides default one.
		SetNamedFormDataValidator(name string) FormHandlerBuilder
		// AddFormExtension adds form extension to the list of form extensions.
		AddFormExtension(formExtension interface{}) FormHandlerBuilder
		// AddNamedFormExtension adds form extension by searching named extension via dingo injector.
		// It panics if there is no injected form extension with that name.
		AddNamedFormExtension(name string) FormHandlerBuilder
		// Build creates new instance of FormHandler interface
		Build() domain.FormHandler
	}

	// formHandlerBuilderImpl as actual implementation of FormHandlerBuilder interface
	formHandlerBuilderImpl struct {
		namedFormServices        map[string]domain.FormService
		namedFormDataProviders   map[string]domain.FormDataProvider
		namedFormDataDecoders    map[string]domain.FormDataDecoder
		namedFormDataValidators  map[string]domain.FormDataValidator
		namedFormExtensions      map[string]domain.FormExtension
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
		validatorProvider        domain.ValidatorProvider
		logger                   flamingo.Logger

		formDataProvider  domain.FormDataProvider
		formDataDecoder   domain.FormDataDecoder
		formDataValidator domain.FormDataValidator
		formExtensions    []interface{}
	}
)

// SetNamedFormService sets form service instance by searching named form service provided via dingo injector.
// It panics if there is no injected form service with that name.
// It overrides provider, decoder and validator if it implements theirs interfaces.
// If it doesn't implements any of those interfaces it panics.
func (b *formHandlerBuilderImpl) SetNamedFormService(name string) FormHandlerBuilder {
	if service, ok := b.namedFormServices[name]; ok {
		return b.SetFormService(service)
	}

	panic(fmt.Sprintf(`there is no FormService with name "%q"`, name))
}

// SetFormService sets form service instance and overrides provider, decoder and validator if
// it implements theirs interfaces. If it doesn't implements any of those interfaces it panics.
func (b *formHandlerBuilderImpl) SetFormService(formService interface{}) FormHandlerBuilder {
	set := false
	if provider, ok := formService.(domain.FormDataProvider); ok {
		b.SetFormDataProvider(provider)
		set = true
	}
	if decoder, ok := formService.(domain.FormDataDecoder); ok {
		b.SetFormDataDecoder(decoder)
		set = true
	}
	if validator, ok := formService.(domain.FormDataValidator); ok {
		b.SetFormDataValidator(validator)
		set = true
	}
	if !set {
		panic("FormService doesn't implement any of FormDataProvider, FormDataDecoder or FormDataValidator interfaces")
	}
	return b
}

// SetNamedFormDataProvider sets form data provider by searching named provider provided via dingo injector.
// It panics if there is no injected form data provider with that name.
// It sets form data provider instance and overrides default one.
func (b *formHandlerBuilderImpl) SetNamedFormDataProvider(name string) FormHandlerBuilder {
	if service, ok := b.namedFormDataProviders[name]; ok {
		return b.SetFormDataProvider(service)
	}

	panic(fmt.Sprintf(`there is no FormDataProvider with name "%q"`, name))
}

// SetFormDataProvider sets form data provider instance and overrides default one.
func (b *formHandlerBuilderImpl) SetFormDataProvider(formDataProvider domain.FormDataProvider) FormHandlerBuilder {
	b.formDataProvider = formDataProvider

	return b
}

// SetNamedFormDataDecoder sets form data decoder by searching named decoder provided via dingo injector.
// It panics if there is no injected form data decoder with that name.
// It sets form data decoder instance and overrides default one.
func (b *formHandlerBuilderImpl) SetNamedFormDataDecoder(name string) FormHandlerBuilder {
	if service, ok := b.namedFormDataDecoders[name]; ok {
		return b.SetFormDataDecoder(service)
	}

	panic(fmt.Sprintf(`there is no FormDataDecoder with name "%q"`, name))
}

// SetFormDataDecoder sets form data decoder instance and overrides default one.
func (b *formHandlerBuilderImpl) SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) FormHandlerBuilder {
	b.formDataDecoder = formDataDecoder

	return b
}

// SetNamedFormDataValidator sets form data decoder by searching named decoder validator via dingo injector.
// It panics if there is no injected form data validator with that name.
// It sets form data validator instance and overrides default one.
func (b *formHandlerBuilderImpl) SetNamedFormDataValidator(name string) FormHandlerBuilder {
	if service, ok := b.namedFormDataValidators[name]; ok {
		return b.SetFormDataValidator(service)
	}

	panic(fmt.Sprintf(`there is no FormDataValidator with name "%q"`, name))
}

// SetFormDataValidator sets form data validator instance and overrides default one.
func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) FormHandlerBuilder {
	b.formDataValidator = formDataValidator

	return b
}

// AddNamedFormExtension adds form extension by searching named extension via dingo injector.
// It panics if there is no injected form extension with that name.
func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) FormHandlerBuilder {
	if service, ok := b.namedFormExtensions[name]; ok {
		return b.AddFormExtension(service)
	}

	panic(fmt.Sprintf(`there is no FormExtension with name "%q"`, name))
}

// AddFormExtension adds form extension to the list of form extensions.
func (b *formHandlerBuilderImpl) AddFormExtension(formExtension interface{}) FormHandlerBuilder {
	implements := false
	if _, ok := formExtension.(domain.FormDataProvider); ok {
		implements = true
	}
	if _, ok := formExtension.(domain.FormDataDecoder); ok {
		implements = true
	}
	if _, ok := formExtension.(domain.FormDataValidator); ok {
		implements = true
	}
	if !implements {
		panic("FormExtension doesn't implement any of FormDataProvider, FormDataDecoder or FormDataValidator interfaces")
	}

	b.formExtensions = append(b.formExtensions, formExtension)

	return b
}

// Build creates new instance of FormHandler interface
func (b *formHandlerBuilderImpl) Build() domain.FormHandler {
	formDataProvider := b.formDataProvider
	if formDataProvider == nil {
		formDataProvider = b.defaultFormDataProvider
	}

	formDataDecoder := b.formDataDecoder
	if formDataDecoder == nil {
		formDataDecoder = b.defaultFormDataDecoder
	}

	formDataValidator := b.formDataValidator
	if formDataValidator == nil {
		formDataValidator = b.defaultFormDataValidator
	}

	return &formHandlerImpl{
		defaultFormDataProvider:  b.defaultFormDataProvider,
		defaultFormDataDecoder:   b.defaultFormDataDecoder,
		defaultFormDataValidator: b.defaultFormDataValidator,
		formDataProvider:         b.formDataProvider,
		formDataDecoder:          b.formDataDecoder,
		formDataValidator:        b.formDataValidator,
		formExtensions:           b.formExtensions,
		validatorProvider:        b.validatorProvider,
		logger:                   b.logger,
	}
}
