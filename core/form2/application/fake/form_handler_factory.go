package fake

import (
	"flamingo.me/flamingo/core/form2/application"
	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/core/form2/domain/mocks"
)

type (
	FormHandlerFactoryImpl struct {
		formHandler *mocks.FormHandler
	}
)

func New(formHandler *mocks.FormHandler) application.FormHandlerFactory {
	return &FormHandlerFactoryImpl{
		formHandler: formHandler,
	}
}

func (f *FormHandlerFactoryImpl) CreateSimpleFormHandler() domain.FormHandler {
	return f.formHandler
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler {
	return f.formHandler
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler {
	return f.formHandler
}

func (f *FormHandlerFactoryImpl) GetFormHandlerBuilder() application.FormHandlerBuilder {
	return &formHandlerBuilderImpl{
		formHandler: f.formHandler,
	}
}
