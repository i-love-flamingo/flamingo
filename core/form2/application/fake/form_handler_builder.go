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

func (b *formHandlerBuilderImpl) SetNamedFormService(name string) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetFormService(formService interface{}) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataProvider(name string) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetFormDataProvider(formDataProvider domain.FormDataProvider) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataDecoder(name string) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataValidator(name string) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) AddFormExtension(formExtension interface{}) application.FormHandlerBuilder {
	return b
}

func (b *formHandlerBuilderImpl) Build() domain.FormHandler {
	return b.formHandler
}
