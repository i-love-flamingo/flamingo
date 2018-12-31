package application

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	formHandlerImpl struct {
		formDataProvider  domain.FormDataProvider
		formDataDecoder   domain.FormDataDecoder
		formDataValidator domain.FormDataValidator
		formExtensions    []interface{}
		validatorProvider domain.ValidatorProvider
	}
)

var _ domain.FormHandler = &formHandlerImpl{}

func (h *formHandlerImpl) HandleForm(ctx context.Context, req *web.Request) (*domain.Form, error) {
	submitted := req.Request().Method == http.MethodPost

	form, err := h.buildForm(ctx, req, submitted)
	if err != nil {
		return nil, err
	}

	if submitted {
		return h.handleSubmittedForm(ctx, req, form)
	}

	return form, nil
}

func (h *formHandlerImpl) HandleUnsubmittedForm(ctx context.Context, req *web.Request) (*domain.Form, error) {
	return h.buildForm(ctx, req, false)
}

func (h *formHandlerImpl) HandleSubmittedForm(ctx context.Context, req *web.Request) (*domain.Form, error) {
	form, err := h.buildForm(ctx, req, true)
	if err != nil {
		return nil, err
	}

	return h.handleSubmittedForm(ctx, req, form)
}

func (h *formHandlerImpl) buildForm(ctx context.Context, req *web.Request, submitted bool) (*domain.Form, error) {
	formData, err := h.formDataProvider.GetFormData(ctx, req)
	if err != nil {
		return nil, err
	}

	form := domain.NewForm(submitted, h.extractValidationRules(formData))
	form.Data = formData

	return &form, nil
}

func (h *formHandlerImpl) handleSubmittedForm(ctx context.Context, req *web.Request, form *domain.Form) (*domain.Form, error) {
	values, err := h.getPostValues(req)
	if err != nil {
		return nil, err
	}

	formData, err := h.formDataDecoder.Decode(ctx, req, *values, form.Data)
	if err != nil {
		return nil, err
	}
	form.Data = formData

	validationInfo, err := h.formDataValidator.Validate(ctx, req, h.validatorProvider, formData)
	if err != nil {
		return nil, err
	}
	form.ValidationInfo = *validationInfo

	err = h.processExtensions(ctx, req, *values, form)
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (h *formHandlerImpl) extractValidationRules(formData interface{}) map[string][]domain.ValidationRule {
	validationRules := map[string][]domain.ValidationRule{}

	if formData == nil {
		return validationRules
	}

	typeOf := reflect.TypeOf(formData)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	if typeOf.Kind() != reflect.Struct {
		return validationRules
	}

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		validationTag := field.Tag.Get("validate")
		if validationTag == "" {
			continue
		}

		name := field.Tag.Get("form")
		if name == "-" || name == "" {
			continue
		}

		tags := strings.Split(validationTag, ",")
		for _, tag := range tags {
			values := strings.Split(tag, "=")
			if len(values) == 0 {
				continue
			}
			if values[0] == "omitempty" || values[0] == "" {
				continue
			}

			validationRule := domain.ValidationRule{
				Name: values[0],
			}
			if len(values) > 1 {
				validationRule.Value = values[1]
			}

			validationRules[name] = append(validationRules[name], validationRule)
		}
	}

	return validationRules
}

func (h *formHandlerImpl) getPostValues(r *web.Request) (*url.Values, error) {
	err := r.Request().ParseForm()
	if err != nil {
		return nil, err
	}

	return &r.Request().Form, nil
}

func (h *formHandlerImpl) processExtensions(ctx context.Context, req *web.Request, values url.Values, form *domain.Form) error {
	for _, formExtension := range h.formExtensions {
		err := h.processExtension(ctx, req, values, formExtension, form)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *formHandlerImpl) processExtension(ctx context.Context, req *web.Request, values url.Values, formExtension interface{}, form *domain.Form) error {
	var formData interface{}
	var err error

	if provider, ok := formExtension.(domain.FormDataProvider); ok {
		formData, err = provider.GetFormData(ctx, req)
		if err != nil {
			return err
		}
	}

	if decoder, ok := formExtension.(domain.FormDataDecoder); ok {
		formData, err = decoder.Decode(ctx, req, values, formData)
		if err != nil {
			return err
		}
	}

	if validator, ok := formExtension.(domain.FormDataValidator); ok {
		validationInfo, err := validator.Validate(ctx, req, h.validatorProvider, formData)
		if err != nil {
			return err
		}

		form.ValidationInfo.AppendGeneralErrors(validationInfo.GetGeneralErrors())
		form.ValidationInfo.AppendFieldErrors(validationInfo.GetErrorsForAllFields())
	}

	return nil
}
