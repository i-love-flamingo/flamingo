package application

import (
	"flamingo.me/flamingo/core/form2/domain"
	"fmt"
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
		namedFormServices        []domain.FormServiceWithName
		namedFormDataProviders   []domain.FormDataProviderWithName
		namedFormDataDecoders    []domain.FormDataDecoderWithName
		namedFormDataValidators  []domain.FormDataValidatorWithName
		namedFormExtensions      []domain.FormExtensionWithName
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
		validatorProvider        domain.ValidatorProvider

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
	for _, service := range b.namedFormServices {
		if name == service.Name() {
			return b.SetFormService(service)
		}
	}

	panic(fmt.Sprintf(`there is no FormService with name "%s"`, name))
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
	for _, provider := range b.namedFormDataProviders {
		if name == provider.Name() {
			return b.SetFormDataProvider(provider)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataProvider with name "%s"`, name))
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
	for _, decoder := range b.namedFormDataDecoders {
		if name == decoder.Name() {
			return b.SetFormDataDecoder(decoder)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataDecoder with name "%s"`, name))
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
	for _, validator := range b.namedFormDataValidators {
		if name == validator.Name() {
			return b.SetFormDataValidator(validator)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataValidator with name "%s"`, name))
}

// SetFormDataValidator sets form data validator instance and overrides default one.
func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) FormHandlerBuilder {
	b.formDataValidator = formDataValidator

	return b
}

// AddNamedFormExtension adds form extension by searching named extension via dingo injector.
// It panics if there is no injected form extension with that name.
func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) FormHandlerBuilder {
	for _, extension := range b.namedFormExtensions {
		if name == extension.Name() {
			return b.AddFormExtension(extension)
		}
	}

	panic(fmt.Sprintf(`there is no FormExtension with name "%s"`, name))
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
		formDataProvider:  formDataProvider,
		formDataDecoder:   formDataDecoder,
		formDataValidator: formDataValidator,
		formExtensions:    b.formExtensions,
		validatorProvider: b.validatorProvider,
	}
}
