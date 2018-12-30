package domain

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/framework/web"
)

type (
	FormHandler interface {
		GetForm(ctx context.Context, req *web.Request) (*Form, error)
		HandleRequest(ctx context.Context, req *web.Request) (*Form, error)
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
		GetFormData(ctx context.Context, req *web.Request) (interface{}, error)
	}

	DefaultFormDataProvider interface {
		FormDataProvider
	}

	FormDataProviderWithName interface {
		NamedFormInstance
		FormDataProvider
	}

	FormDataDecoder interface {
		Decode(ctx context.Context, req *web.Request, values url.Values, formData interface{}) (interface{}, error)
	}

	DefaultFormDataDecoder interface {
		FormDataDecoder
	}

	FormDataDecoderWithName interface {
		NamedFormInstance
		FormDataDecoder
	}

	FormDataValidator interface {
		Validate(ctx context.Context, req *web.Request, validatorProvider ValidatorProvider, formData interface{}) (*ValidationInfo, error)
	}

	DefaultFormDataValidator interface {
		FormDataValidator
	}

	FormDataValidatorWithName interface {
		NamedFormInstance
		FormDataValidator
	}
)
