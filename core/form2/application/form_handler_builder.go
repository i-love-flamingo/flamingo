package application

import "flamingo.me/flamingo/core/form2/domain"

type (
	FormHandlerBuilder interface {
		SetFormService(formService interface{})
		SetNamedFormService(name string)
		SetFormDataProvider(formDataProvider domain.FormDataProvider)
		SetNamedFormDataProvider(name string)
		SetFormDataDecoder(formDataDecoder domain.FormDataDecoder)
		SetNamedFormDataDecoder(name string)
		SetFormDataValidator(formDataValidator domain.FormDataValidator)
		SetNamedFormDataValidator(name string)
		SetFormDataMapper(formDataMapper domain.FormDataMapper)
		SetNamedFormDataMapper(name string)
		AddFormExtension(formExtension interface{})
		AddNamedFormExtension(name string)
		Build() domain.FormHandler
	}
)
