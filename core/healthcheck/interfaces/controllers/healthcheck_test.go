package controllers

import (
	"context"
	"testing"

	"strings"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
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

func TestController_Healthcheck(t *testing.T) {
	type fields struct {
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
			controller.Inject(&web.Responder{}, tt.fields.statusProvider)

			result := controller.Healthcheck(tt.args.ctx, tt.args.request)
			response, ok := result.(*web.DataResponse)
			assert.True(t, ok)
			assert.Equal(t, tt.want, response.Data)
		})
	}
}

func TestController_Ping(t *testing.T) {
	controller := &Healthcheck{}
	controller.Inject(&web.Responder{}, nil)

	result := controller.Ping(nil, nil)
	response, ok := result.(*web.HTTPResponse)
	assert.True(t, ok)
	assert.Equal(t, strings.NewReader("OK"), response.Body)
}
