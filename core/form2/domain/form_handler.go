package domain

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/framework/web"
)

type (
	// FormHandler is interface for defining main form processor which provider instance of Form domain
	FormHandler interface {
		// HandleUnsubmittedForm as method for returning Form instance which is not submitted
		HandleUnsubmittedForm(ctx context.Context, req *web.Request) (*Form, error)
		// HandleSubmittedForm as method for returning Form instance which is submitted
		HandleSubmittedForm(ctx context.Context, req *web.Request) (*Form, error)
		// HandleForm as method for returning Form instance with state depending on fact if there was form submission or not
		HandleForm(ctx context.Context, req *web.Request) (*Form, error)
	}

	// NamedFormInstance is interface for defining all form services and extensions with names
	NamedFormInstance interface {
		// Name as method for defining name of form service and extension
		Name() string
	}

	// FormExtensionWithName is interface for defining all form extensions with names
	FormExtensionWithName interface {
		NamedFormInstance
	}

	// FormServiceWithName is interface for defining all form services with names
	FormServiceWithName interface {
		NamedFormInstance
	}

	// FormDataProvider is interface for defining all form services which creates form data
	FormDataProvider interface {
		// GetFormData as method for defining form data
		GetFormData(ctx context.Context, req *web.Request) (interface{}, error)
	}

	// DefaultFormDataProvider is interface for defining default form data provider
	// used in case when there is no custom form data provider defined
	DefaultFormDataProvider interface {
		FormDataProvider
	}

	// FormDataProviderWithName is interface for defining all form data provider with names
	FormDataProviderWithName interface {
		NamedFormInstance
		FormDataProvider
	}

	// FormDataDecoder is interface for defining all form services which process http request and transform it into form data
	FormDataDecoder interface {
		// Decode as method for transforming http request body into form data
		Decode(ctx context.Context, req *web.Request, values url.Values, formData interface{}) (interface{}, error)
	}

	// DefaultFormDataDecoder is interface for defining default form data decoder
	// used in case when there is no custom form data decoder defined
	DefaultFormDataDecoder interface {
		FormDataDecoder
	}

	// FormDataDecoderWithName is interface for defining all form data decoder with names
	FormDataDecoderWithName interface {
		NamedFormInstance
		FormDataDecoder
	}

	// FormDataValidator is interface for defining all form services which validates form data
	FormDataValidator interface {
		// Validate as method for validating form data
		Validate(ctx context.Context, req *web.Request, validatorProvider ValidatorProvider, formData interface{}) (*ValidationInfo, error)
	}

	// DefaultFormDataValidator is interface for defining default form data validator
	// used in case when there is no custom form data validator defined
	DefaultFormDataValidator interface {
		FormDataValidator
	}

	// FormDataValidatorWithName is interface for defining all form data validators with names
	FormDataValidatorWithName interface {
		NamedFormInstance
		FormDataValidator
	}

	// CompleteFormService is interface for defining all form services which can acts as provider, decoder and validator
	CompleteFormService interface {
		FormDataProvider
		FormDataDecoder
		FormDataValidator
	}

	// CompleteFormServiceWithName is interface for defining all form services with names
	// which can acts as provider, decoder and validator
	CompleteFormServiceWithName interface {
		NamedFormInstance
		CompleteFormService
	}
)
