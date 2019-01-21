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
		// HandleSubmittedForm as method for returning Form instance which is submitted via POST request
		HandleSubmittedForm(ctx context.Context, req *web.Request) (*Form, error)
		// HandleSubmittedGETForm as method for returning Form instance which is submitted via GET request
		HandleSubmittedGETForm(ctx context.Context, req *web.Request) (*Form, error)
		// HandleForm as method for returning Form instance with state depending on fact if there was form submission or not, via POST request
		HandleForm(ctx context.Context, req *web.Request) (*Form, error)
	}

	// FormExtension is helper interface for form extensions used for binding with dingo injector
	FormExtension interface{}

	// FormService is helper interface for form services used for binding with dingo injector
	FormService interface{}

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

	// CompleteFormService is interface for defining all form services which can acts as provider, decoder and validator
	CompleteFormService interface {
		FormDataProvider
		FormDataDecoder
		FormDataValidator
	}
)
