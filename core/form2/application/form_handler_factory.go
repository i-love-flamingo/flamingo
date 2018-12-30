package application

import "flamingo.me/flamingo/core/form2/domain"

type (
	FormHandlerFactory interface {
		CreateSimpleFormHandler() domain.FormHandler
		CreateFormHandlerWithFormService(formService interface{}) domain.FormHandler
		CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formDataMapper domain.FormDataMapper) domain.FormHandler
		GetFormHandlerBuilder() FormHandlerBuilder
	}
)
