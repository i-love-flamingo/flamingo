package fake

import (
	"flamingo.me/flamingo/v3/core/form2/application"
	"flamingo.me/flamingo/v3/core/form2/domain"
	"flamingo.me/flamingo/v3/core/form2/domain/mocks"
)

type (
	// FormHandlerFactoryImpl defines faked implementation of FormHandlerFactory interface used for unit testing
	FormHandlerFactoryImpl struct {
		formHandler *mocks.FormHandler
	}
)

// New returns faked implementation of FormHandlerFactory interface which should deliver mocked domain.FormHandler instance
func New(formHandler *mocks.FormHandler) application.FormHandlerFactory {
	return &FormHandlerFactoryImpl{
		formHandler: formHandler,
	}
}

// CreateSimpleFormHandler returns mocked instance of domain.FormHandler interface
func (f *FormHandlerFactoryImpl) CreateSimpleFormHandler() domain.FormHandler {
	return f.formHandler
}

// CreateFormHandlerWithFormService returns mocked instance of domain.FormHandler interface
func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormService(domain.FormService, ...string) domain.FormHandler {
	return f.formHandler
}

// CreateFormHandlerWithFormServices returns mocked instance of domain.FormHandler interface
func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormServices(domain.FormDataProvider, domain.FormDataDecoder, domain.FormDataValidator, ...string) domain.FormHandler {
	return f.formHandler
}

// GetFormHandlerBuilder returns faked instance of FormHandlerBuilder interface
func (f *FormHandlerFactoryImpl) GetFormHandlerBuilder() application.FormHandlerBuilder {
	return &formHandlerBuilderImpl{
		formHandler: f.formHandler,
	}
}
