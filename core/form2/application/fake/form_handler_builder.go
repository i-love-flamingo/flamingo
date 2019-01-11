package fake

import (
	"flamingo.me/flamingo/core/form2/application"
	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/core/form2/domain/mocks"
)

type (
	formHandlerBuilderImpl struct {
		formHandler *mocks.FormHandler
	}
)

var _ application.FormHandlerBuilder = &formHandlerBuilderImpl{}

func (b *formHandlerBuilderImpl) SetNamedFormService(name string) error {
	return nil
}

func (b *formHandlerBuilderImpl) SetFormService(formService domain.FormService) error {
	return nil
}

func (b *formHandlerBuilderImpl) SetNamedFormDataProvider(name string) error {
	return nil
}

func (b *formHandlerBuilderImpl) SetFormDataProvider(formDataProvider domain.FormDataProvider) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataDecoder(name string) error {
	return nil
}

func (b *formHandlerBuilderImpl) SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataValidator(name string) error {
	return nil
}

func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) error {
	return nil
}

func (b *formHandlerBuilderImpl) AddFormExtension(formExtension domain.FormExtension) error {
	return nil
}

func (b *formHandlerBuilderImpl) Must(error) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) Build() domain.FormHandler {
	return b.formHandler
}
