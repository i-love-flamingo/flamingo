package fake

import (
	"flamingo.me/flamingo/v3/core/form2/application"
	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/core/form2/domain/mocks"
)

type (
	// formHandlerBuilderImpl defines faked implementation of FormHandlerBuilder interface used for unit testing
	formHandlerBuilderImpl struct {
		formHandler *mocks.FormHandler
	}
)

var _ application.FormHandlerBuilder = &formHandlerBuilderImpl{}

// SetNamedFormService fakes storing of named form service into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetNamedFormService(name string) error {
	return nil
}

// SetFormService fakes storing of form service into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetFormService(formService domain.FormService) error {
	return nil
}

// SetNamedFormDataProvider fakes storing of named form data provider into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetNamedFormDataProvider(name string) error {
	return nil
}

// SetFormDataProvider fakes storing of form data provider into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetFormDataProvider(formDataProvider domain.FormDataProvider) application.FormHandlerBuilder {
	return b
}

// SetNamedFormDataDecoder fakes storing of named form data decoder into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetNamedFormDataDecoder(name string) error {
	return nil
}

// SetFormDataDecoder fakes storing of form data decoder into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) application.FormHandlerBuilder {
	return b
}

// SetNamedFormDataValidator fakes storing of named form data validator into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetNamedFormDataValidator(name string) error {
	return nil
}

// SetFormDataValidator fakes storing of form data validator into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) application.FormHandlerBuilder {
	return b
}

// AddNamedFormExtension fakes storing of named form extension into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) error {
	return nil
}

// AddFormExtension fakes storing of form extension into mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) AddFormExtension(formExtension domain.FormExtension) error {
	return nil
}

// Must fakes storing wrapping of methods that can returns error message.
func (b *formHandlerBuilderImpl) Must(error) application.FormHandlerBuilder {
	return b
}

// Build returns mocked instance of domain.FormHandler.
func (b *formHandlerBuilderImpl) Build() domain.FormHandler {
	return b.formHandler
}
