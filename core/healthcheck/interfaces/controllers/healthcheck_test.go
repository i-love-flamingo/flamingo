package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
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
		statusProvider statusProvider
	}
	type args struct {
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
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
				request: nil,
			},
			want: "{\"services\":[{\"name\":\"test\",\"alive\":true,\"details\":\"alive\"}]}",
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
				request: nil,
			},
			want: "{\"services\":[{\"name\":\"test\",\"alive\":false,\"details\":\"not alive\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := &Healthcheck{}
			controller.Inject(tt.fields.statusProvider)

			recorder := httptest.NewRecorder()
			controller.ServeHTTP(recorder, tt.args.request)

			resp := recorder.Result()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, string(body))
		})
	}
}

func TestController_Ping(t *testing.T) {
	controller := &Ping{}

	recorder := httptest.NewRecorder()
	controller.ServeHTTP(recorder, nil)

	resp := recorder.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK", string(body))
}
