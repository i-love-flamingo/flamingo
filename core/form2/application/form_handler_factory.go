package application

import "flamingo.me/flamingo/core/form2/domain"

type (
	FormHandlerFactory interface {
		CreateSimpleFormHandler() domain.FormHandler
		CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler
		CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler
		GetFormHandlerBuilder() FormHandlerBuilder
	}

	FormHandlerFactoryImpl struct {
		formServices             []domain.FormServiceWithName
		formDataProviders        []domain.FormDataProviderWithName
		formDataDecoders         []domain.FormDataDecoderWithName
		formDataValidators       []domain.FormDataValidatorWithName
		formExtensions           []domain.FormExtensionWithName
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
	}
)

func (f *FormHandlerFactoryImpl) Inject(
	s []domain.FormServiceWithName,
	p []domain.FormDataProviderWithName,
	d []domain.FormDataDecoderWithName,
	v []domain.FormDataValidatorWithName,
	e []domain.FormExtensionWithName,
	dp domain.DefaultFormDataProvider,
	dd domain.DefaultFormDataDecoder,
	dv domain.DefaultFormDataValidator,
) {
	f.formServices = s
	f.formDataProviders = p
	f.formDataDecoders = d
	f.formDataValidators = v
	f.formExtensions = e
	f.defaultFormDataProvider = dp
	f.defaultFormDataDecoder = dd
	f.defaultFormDataValidator = dv
}

func (f *FormHandlerFactoryImpl) CreateSimpleFormHandler() domain.FormHandler {
	return f.GetFormHandlerBuilder().Build()
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().SetFormService(formService)
	f.attachExtensions(builder, formExtensions)
	return builder.Build()
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().
		SetFormDataProvider(formDataProvider).
		SetFormDataDecoder(formDataDecoder).
		SetFormDataValidator(formDataValidator)
	f.attachExtensions(builder, formExtensions)
	return builder.Build()
}

func (f *FormHandlerFactoryImpl) GetFormHandlerBuilder() FormHandlerBuilder {
	return &FormHandlerBuilderImpl{
		formServices:             f.formServices,
		formDataProviders:        f.formDataProviders,
		formDataDecoders:         f.formDataDecoders,
		formDataValidators:       f.formDataValidators,
		formExtensions:           f.formExtensions,
		defaultFormDataProvider:  f.defaultFormDataProvider,
		defaultFormDataDecoder:   f.defaultFormDataDecoder,
		defaultFormDataValidator: f.defaultFormDataValidator,
	}
}

func (f *FormHandlerFactoryImpl) attachExtensions(builder FormHandlerBuilder, formExtensions ...interface{}) {
	for _, item := range formExtensions {
		if name, ok := item.(string); ok {
			builder.AddNamedFormExtension(name)
		} else {
			builder.AddFormExtension(item)
		}
	}
}
