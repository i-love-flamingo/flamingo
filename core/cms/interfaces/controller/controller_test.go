package controller

import (
	"context"
	"errors"
	"flamingo/core/cms/domain"
	"flamingo/framework/router"
	"flamingo/framework/testutil"
	"flamingo/framework/web"
	"testing"
)

type (
	MockPageService  struct{}
	MockBlockService struct{}
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
	var expectedTpl string

	vc := &ViewController{
		PageService: new(MockPageService),
		ErrorAware: &testutil.MockErrorAware{
			CbError: func(context web.Context, err error) web.Response {
				_, testerr := (&MockPageService{}).Get(nil, "fail")
				if err.Error() != testerr.Error() {
					t.Errorf("wrong error %v returned", err)
				}
				return nil
			},
			CbErrorNotFound: func(context web.Context, error error) web.Response {
				return nil
			},
		},
		RenderAware: &testutil.MockRenderAware{
			CbRender: func(context web.Context, tpl string, data interface{}) web.Response {
				if tpl != expectedTpl {
					t.Errorf("Expected template %q != %q", expectedTpl, tpl)
				}
				return nil
			},
		},
	}
	ctx := web.NewContext()

	ctx.LoadParams(router.P{"name": "fail"})

	res := vc.Get(ctx)
	if res != nil {
		t.Errorf("Expected mocked response to be nil, not %v", res)
	}

	ctx.LoadParams(router.P{"name": "notfound"})

	res = vc.Get(ctx)
	if res != nil {
		t.Errorf("Expected mocked response to be nil, not %v", res)
	}

	expectedTpl = "cms/view"
	ctx.LoadParams(router.P{"name": "page"})

	res = vc.Get(ctx)
	if res != nil {
		t.Errorf("Expected mocked response to be nil, not %v", res)
	}

	expectedTpl = "cms/template"
	ctx.LoadParams(router.P{"name": "page", "template": "cms/template"})

	res = vc.Get(ctx)
	if res != nil {
		t.Errorf("Expected mocked response to be nil, not %v", res)
	}
}

func (m *MockBlockService) Get(ctx context.Context, identifier string) (*domain.Block, error) {
	return &domain.Block{
		Identifier: identifier,
	}, nil
}

func TestDataController_Data(t *testing.T) {
	dc := &DataController{
		BlockService: new(MockBlockService),
	}
	ctx := web.NewContext()

	ctx.LoadParams(router.P{"block": "test"})

	block, ok := dc.Data(ctx).(*domain.Block)
	if !ok {
		t.Error("Result not a block")
	}

	if block == nil {
		t.Error("Block is nil")
	}

	if block.Identifier != "test" {
		t.Errorf("Wrong identifier %q, expected test", block.Identifier)
	}
}
