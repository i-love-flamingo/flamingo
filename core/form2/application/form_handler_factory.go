package application

import (
	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// FormHandlerFactory as interface for simpler creation of form handler instance
	FormHandlerFactory interface {
		// CreateSimpleFormHandler as method for creating the simplest form handler instance which uses
		// default form data provider, decoder and validator
		CreateSimpleFormHandler() domain.FormHandler
		// CreateFormHandlerWithFormService as method for creating customized form handler.
		// Form service must implement at least one of the provider, decoder or validator interface, and it's methods
		// are used to override default form data provider, decoder and validator.
		// From extensions must implement at least one of the provider, decoder or validator interface, and they are
		// used to add additional form functionality, like validation which is attached to final validation info.
		// Form extensions can be passed as instances or by their names, which reflect named extensions injected via dingo injector.
		CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler
		// CreateFormHandlerWithFormServices as method for creating customized form handler.
		// It expect instances provider, decoder or validator interface, and it's methods
		// are used to override default form data provider, decoder and validator.
		// If nil values are passed for provider, decoder or validator, default form data provider, decoder and validator
		// are used.
		// From extensions must implement at least one of the provider, decoder or validator interface, and they are
		// used to add additional form functionality, like validation which is attached to final validation info.
		// Form extensions can be passed as instances or by their names, which reflect named extensions injected via dingo injector.
		CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler
		// GetFormHandlerBuilder returns FomHandlerBuilder for creating more complex instances of form handler.
		GetFormHandlerBuilder() FormHandlerBuilder
	}

	// FormHandlerFactoryImpl as actual implementation of FormHandlerFactory interface
	FormHandlerFactoryImpl struct {
		namedFormServices        map[string]domain.FormService
		namedFormDataProviders   map[string]domain.FormDataProvider
		namedFormDataDecoders    map[string]domain.FormDataDecoder
		namedFormDataValidators  map[string]domain.FormDataValidator
		namedFormExtensions      map[string]domain.FormExtension
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
		validatorProvider        domain.ValidatorProvider
		logger                   flamingo.Logger
	}
)

func (f *FormHandlerFactoryImpl) Inject(
	s map[string]domain.FormService,
	p map[string]domain.FormDataProvider,
	d map[string]domain.FormDataDecoder,
	v map[string]domain.FormDataValidator,
	e map[string]domain.FormExtension,
	dp domain.DefaultFormDataProvider,
	dd domain.DefaultFormDataDecoder,
	dv domain.DefaultFormDataValidator,
	vp domain.ValidatorProvider,
	l flamingo.Logger,
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
	f.logger = l
}

// CreateSimpleFormHandler as method for creating the simplest form handler instance which uses
// default form data provider, decoder and validator
func (f *FormHandlerFactoryImpl) CreateSimpleFormHandler() domain.FormHandler {
	return f.GetFormHandlerBuilder().Build()
}

// CreateFormHandlerWithFormService as method for creating customized form handler.
// Form service must implement at least one of the provider, decoder or validator interface, and it's methods
// are used to override default form data provider, decoder and validator.
// From extensions must implement at least one of the provider, decoder or validator interface, and they are
// used to add additional form functionality, like validation which is attached to final validation info.
// Form extensions can be passed as instances or by their names, which reflect named extensions injected via dingo injector.
func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormService(formService interface{}, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().SetFormService(formService)
	f.attachExtensions(builder, formExtensions...)
	return builder.Build()
}

// CreateFormHandlerWithFormServices as method for creating customized form handler.
// It expect instances provider, decoder or validator interface, and it's methods
// are used to override default form data provider, decoder and validator.
// If nil values are passed for provider, decoder or validator, default form data provider, decoder and validator
// are used.
// From extensions must implement at least one of the provider, decoder or validator interface, and they are
// used to add additional form functionality, like validation which is attached to final validation info.
// Form extensions can be passed as instances or by their names, which reflect named extensions injected via dingo injector.
func (f *FormHandlerFactoryImpl) CreateFormHandlerWithFormServices(formDataProvider domain.FormDataProvider, formDataDecoder domain.FormDataDecoder, formDataValidator domain.FormDataValidator, formExtensions ...interface{}) domain.FormHandler {
	builder := f.GetFormHandlerBuilder().
		SetFormDataProvider(formDataProvider).
		SetFormDataDecoder(formDataDecoder).
		SetFormDataValidator(formDataValidator)
	f.attachExtensions(builder, formExtensions...)
	return builder.Build()
}

// GetFormHandlerBuilder returns FomHandlerBuilder for creating more complex instances of form handler.
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
		logger:                   f.logger,
	}
}

// attachExtensions method for attaching form extension to the list of extensions.
// It expects string as form extension's name or actual instance of form extension
func (f *FormHandlerFactoryImpl) attachExtensions(builder FormHandlerBuilder, formExtensions ...interface{}) {
	for _, item := range formExtensions {
		if name, ok := item.(string); ok {
			builder.AddNamedFormExtension(name)
		} else {
			builder.AddFormExtension(item)
		}
	}
}
