package domain

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/framework/web"
)

type (
	FormHandler interface {
		GetForm(ctx context.Context, req *web.Request) Form
		HandleRequest(ctx context.Context, req *web.Request) Form
		MapFormData(formData interface{}, mappedData interface{})
	}

	NamedFormInstance interface {
		Name() string
	}

	FormExtensionWithName interface {
		NamedFormInstance
	}

	FormServiceWithName interface {
		NamedFormInstance
	}

	FormDataProvider interface {
		GetFormData(ctx context.Context, req *web.Request) interface{}
	}

	DefaultFormDataProvider interface {
		FormDataProvider
	}

	FormDataProviderWithName interface {
		NamedFormInstance
		FormDataProvider
	}

	FormDataDecoder interface {
		Decode(ctx context.Context, req *web.Request, values *url.Values, formData interface{}) interface{}
	}

	DefaultFormDataDecode interface {
		FormDataDecoder
	}

	FormDataDecoderWithName interface {
		NamedFormInstance
		FormDataDecoder
	}

	FormDataValidator interface {
		Validate(ctx context.Context, req *web.Request, validatorProvider ValidatorProvider, formData interface{}) ValidationInfo
	}

	DefaultFormDataValidator interface {
		FormDataValidator
	}

	FormDataValidatorWithName interface {
		NamedFormInstance
		FormDataValidator
	}

	FormDataMapper interface {
		Map(formData interface{}, mappedData interface{}) bool
	}

	FormDataMapperWithName interface {
		NamedFormInstance
		FormDataMapper
	}
)
