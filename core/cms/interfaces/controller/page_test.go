package controller

import (
	"context"
	"errors"
	"flamingo/core/cms/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

type (
	MockPageService struct{}
)

func (m *MockPageService) Get(ctx context.Context, identifier string) (*domain.Page, error) {
	if identifier == "fail" {
		return nil, errors.New("failed")
	}

	if identifier == "notfound" {
		return nil, nil
	}

	return &domain.Page{
		Identifier: identifier,
	}, nil
}

func TestViewController_Get(t *testing.T) {
	ctx := web.NewContext()

	errorAware := new(mocks.ErrorAware)
	errorAware.On("Error", ctx, mock.Anything).Return(nil)
	errorAware.On("ErrorNotFound", ctx, mock.Anything).Return(nil)

	renderAware := new(mocks.RenderAware)
	renderAware.On("Render", ctx, mock.AnythingOfType("string"), mock.Anything).Return(nil)

	vc := &ViewController{
		PageService: new(MockPageService),
		ErrorAware:  errorAware,
		RenderAware: renderAware,
	}

	ctx.LoadParams(router.P{"name": "fail"})
	vc.Get(ctx)
	errorAware.AssertCalled(t, "Error", ctx, mock.Anything)

	ctx.LoadParams(router.P{"name": "notfound"})
	vc.Get(ctx)
	errorAware.AssertCalled(t, "ErrorNotFound", ctx, mock.Anything)

	expectedTpl := "cms/cms"
	ctx.LoadParams(router.P{"name": "page"})
	vc.Get(ctx)
	renderAware.AssertCalled(t, "Render", ctx, expectedTpl, mock.Anything)

	expectedTpl = "cms/template"
	ctx.LoadParams(router.P{"name": "page", "template": "cms/template"})
	vc.Get(ctx)
	renderAware.AssertCalled(t, "Render", ctx, expectedTpl, mock.Anything)
}
