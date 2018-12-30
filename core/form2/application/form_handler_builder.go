package application

import (
	"flamingo.me/flamingo/core/form2/domain"
	"fmt"
)

type (
	FormHandlerBuilder interface {
		SetFormService(formService interface{}) FormHandlerBuilder
		SetNamedFormService(name string) FormHandlerBuilder
		SetFormDataProvider(formDataProvider domain.FormDataProvider) FormHandlerBuilder
		SetNamedFormDataProvider(name string) FormHandlerBuilder
		SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) FormHandlerBuilder
		SetNamedFormDataDecoder(name string) FormHandlerBuilder
		SetFormDataValidator(formDataValidator domain.FormDataValidator) FormHandlerBuilder
		SetNamedFormDataValidator(name string) FormHandlerBuilder
		AddFormExtension(formExtension interface{}) FormHandlerBuilder
		AddNamedFormExtension(name string) FormHandlerBuilder
		Build() domain.FormHandler
	}

	formHandlerBuilderImpl struct {
		formServices             []domain.FormServiceWithName
		formDataProviders        []domain.FormDataProviderWithName
		formDataDecoders         []domain.FormDataDecoderWithName
		formDataValidators       []domain.FormDataValidatorWithName
		formExtensions           []domain.FormExtensionWithName
		defaultFormDataProvider  domain.DefaultFormDataProvider
		defaultFormDataDecoder   domain.DefaultFormDataDecoder
		defaultFormDataValidator domain.DefaultFormDataValidator
		validatorProvider        domain.ValidatorProvider

		formDataProvider  domain.FormDataProvider
		formDataDecoder   domain.FormDataDecoder
		formDataValidator domain.FormDataValidator
		formExtensionList []interface{}
	}
)

func (b *formHandlerBuilderImpl) SetNamedFormService(name string) FormHandlerBuilder {
	for _, service := range b.formServices {
		if name == service.Name() {
			return b.SetFormService(service)
		}
	}

	panic(fmt.Sprintf(`there is no FormService with name "%s"`, name))
}

func (b *formHandlerBuilderImpl) SetFormService(formService interface{}) FormHandlerBuilder {
	set := false
	if provider, ok := formService.(domain.FormDataProvider); ok {
		b.SetFormDataProvider(provider)
		set = true
	}
	if decoder, ok := formService.(domain.FormDataDecoder); ok {
		b.SetFormDataDecoder(decoder)
		set = true
	}
	if validator, ok := formService.(domain.FormDataValidator); ok {
		b.SetFormDataValidator(validator)
		set = true
	}
	if !set {
		panic("FormService doesn't implement any of FormDataProvider, FormDataDecoder or FormDataValidator interfaces")
	}
	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataProvider(name string) FormHandlerBuilder {
	for _, provider := range b.formDataProviders {
		if name == provider.Name() {
			return b.SetFormDataProvider(provider)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataProvider with name "%s"`, name))
}

func (b *formHandlerBuilderImpl) SetFormDataProvider(formDataProvider domain.FormDataProvider) FormHandlerBuilder {
	b.formDataProvider = formDataProvider

	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataDecoder(name string) FormHandlerBuilder {
	for _, decoder := range b.formDataDecoders {
		if name == decoder.Name() {
			return b.SetFormDataDecoder(decoder)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataDecoder with name "%s"`, name))
}

func (b *formHandlerBuilderImpl) SetFormDataDecoder(formDataDecoder domain.FormDataDecoder) FormHandlerBuilder {
	b.formDataDecoder = formDataDecoder

	return b
}

func (b *formHandlerBuilderImpl) SetNamedFormDataValidator(name string) FormHandlerBuilder {
	for _, validator := range b.formDataValidators {
		if name == validator.Name() {
			return b.SetFormDataValidator(validator)
		}
	}

	panic(fmt.Sprintf(`there is no FormDataValidator with name "%s"`, name))
}

func (b *formHandlerBuilderImpl) SetFormDataValidator(formDataValidator domain.FormDataValidator) FormHandlerBuilder {
	b.formDataValidator = formDataValidator

	return b
}

func (b *formHandlerBuilderImpl) AddNamedFormExtension(name string) FormHandlerBuilder {
	for _, extension := range b.formExtensions {
		if name == extension.Name() {
			return b.AddFormExtension(extension)
		}
	}

	panic(fmt.Sprintf(`there is no FormExtension with name "%s"`, name))
}

func (b *formHandlerBuilderImpl) AddFormExtension(formExtension interface{}) FormHandlerBuilder {
	implements := false
	if _, ok := formExtension.(domain.FormDataProvider); ok {
		implements = true
	}
	if _, ok := formExtension.(domain.FormDataDecoder); ok {
		implements = true
	}
	if _, ok := formExtension.(domain.FormDataValidator); ok {
		implements = true
	}
	if !implements {
		panic("FormExtension doesn't implement any of FormDataProvider, FormDataDecoder or FormDataValidator interfaces")
	}

	b.formExtensionList = append(b.formExtensionList, formExtension)

	return b
}

func (b *formHandlerBuilderImpl) Build() domain.FormHandler {
	formDataProvider := b.formDataProvider
	if formDataProvider == nil {
		formDataProvider = b.defaultFormDataProvider
	}

	formDataDecoder := b.formDataDecoder
	if formDataDecoder == nil {
		formDataDecoder = b.defaultFormDataDecoder
	}

	formDataValidator := b.formDataValidator
	if formDataValidator == nil {
		formDataValidator = b.defaultFormDataValidator
	}

	return &formHandlerImpl{
		formDataProvider:  formDataProvider,
		formDataDecoder:   formDataDecoder,
		formDataValidator: formDataValidator,
		formExtensionList: b.formExtensionList,
		validatorProvider: b.validatorProvider,
	}
}
