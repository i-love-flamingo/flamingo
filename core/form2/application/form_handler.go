package application

import (
	"context"

	"flamingo.me/flamingo/core/form2/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	FormHandlerImpl struct {
		formDataProvider  domain.FormDataProvider
		formDataDecoder   domain.FormDataDecoder
		formDataValidator domain.FormDataValidator
		formExtensionList []interface{}
	}
)

var _ domain.FormHandler = &FormHandlerImpl{}

func (h *FormHandlerImpl) GetForm(ctx context.Context, req *web.Request) (*domain.Form, error) {
	return nil, nil
}

func (h *FormHandlerImpl) HandleRequest(ctx context.Context, req *web.Request) (*domain.Form, error) {
	return nil, nil
}
