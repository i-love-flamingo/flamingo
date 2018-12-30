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
		namedFormServices        []domain.FormServiceWithName
		namedFormDataProviders   []domain.FormDataProviderWithName
		namedFormDataDecoders    []domain.FormDataDecoderWithName
		namedFormDataValidators  []domain.FormDataValidatorWithName
		namedFormExtensions      []domain.FormExtensionWithName
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
		validatorProvider        domain.ValidatorProvider
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
	vp domain.ValidatorProvider,
) {
	f.namedFormServices = s
	f.namedFormDataProviders = p
	f.namedFormDataDecoders = d
	f.namedFormDataValidators = v
	f.namedFormExtensions = e
	f.defaultFormDataProvider = dp
	f.defaultFormDataDecoder = dd
	f.defaultFormDataValidator = dv
	f.validatorProvider = vp
}

func (f *FormHandlerFactoryImpl) CreateSimpleFormHandler() domain.FormHandler {
	return f.GetFormHandlerBuilder().Build()
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().SetFormService(formService)
	f.attachExtensions(builder, formExtensions...)
	return builder.Build()
}

func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().
		SetFormDataProvider(formDataProvider).
		SetFormDataDecoder(formDataDecoder).
		SetFormDataValidator(formDataValidator)
	f.attachExtensions(builder, formExtensions...)
	return builder.Build()
}

func (f *FormHandlerFactoryImpl) GetFormHandlerBuilder() FormHandlerBuilder {
	return &formHandlerBuilderImpl{
		namedFormServices:        f.namedFormServices,
		namedFormDataProviders:   f.namedFormDataProviders,
		namedFormDataDecoders:    f.namedFormDataDecoders,
		namedFormDataValidators:  f.namedFormDataValidators,
		namedFormExtensions:      f.namedFormExtensions,
		defaultFormDataProvider:  f.defaultFormDataProvider,
		defaultFormDataDecoder:   f.defaultFormDataDecoder,
		defaultFormDataValidator: f.defaultFormDataValidator,
		validatorProvider:        f.validatorProvider,
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
