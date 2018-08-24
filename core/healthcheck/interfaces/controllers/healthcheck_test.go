package controllers

import (
	"context"
	"testing"

	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"

	"flamingo.me/flamingo/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/framework/web/responder/mocks"
)

type (
	testStatus struct {
		alive bool
		text  string
	}
)

func (t *testStatus) Status() (alive bool, details string) {
	return t.alive, t.text
}

func TestHealthcheck_Get(t *testing.T) {
	type fields struct {
		jsonAware      responder.JSONAware
		statusProvider StatusProvider
	}
	type args struct {
		ctx     context.Context
		request *web.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   response
	}{
		{
			name: "alive",
			fields: fields{
				jsonAware: new(mocks.JSONAware),
				statusProvider: func() map[string]healthcheck.Status {
					return map[string]healthcheck.Status{
						"test": &testStatus{true, "alive"},
					}
				},
			},
			args: args{
				ctx: context.Background(),
				request: &web.Request{
					Values: nil,
				},
			},
			want: response{[]service{{Name: "test", Alive: true, Details: "alive"}}},
		},
		{
			name: "not alive",
			fields: fields{
				jsonAware: new(mocks.JSONAware),
				statusProvider: func() map[string]healthcheck.Status {
					return map[string]healthcheck.Status{
						"test": &testStatus{false, "not alive"},
					}
				},
			},
			args: args{
				ctx: context.Background(),
				request: &web.Request{
					Values: nil,
				},
			},
			want: response{Services: []service{{Name: "test", Alive: false, Details: "not alive"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &Healthcheck{}
			controller.Inject(tt.fields.jsonAware, tt.fields.statusProvider)

			tt.fields.jsonAware.(*mocks.JSONAware).On("JSON", tt.want).Once().Return(nil)

			_ = controller.Get(tt.args.ctx, tt.args.request)

			tt.fields.jsonAware.(*mocks.JSONAware).AssertExpectations(t)
		})
	}
}
