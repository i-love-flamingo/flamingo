package redirects

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"testing"

	"flamingo.me/flamingo/v3/core/redirects/infrastructure"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/router"
	"flamingo.me/flamingo/v3/framework/router/mocks"
	"flamingo.me/flamingo/v3/framework/web"
	responderMocks "flamingo.me/flamingo/v3/framework/web/responder/mocks"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

type (
	filterMocker func(ctx context.Context, r *web.Request, w http.ResponseWriter) web.Response
)

func (fnc filterMocker) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	return fnc(ctx, r, w)
}

func Test_redirector_TryServeHTTP(t *testing.T) {
	type fields struct {
		redirectDataMap map[string]infrastructure.CsvContent
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantModifyHeader bool
		want             bool
		wantStatus       int
		wantLocation     string
		wantErr          bool
	}{
		{
			name: "valid 302 redirect",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusFound,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri",
				},
			},
			wantModifyHeader: true,
			want:             false,
			wantStatus:       http.StatusFound,
			wantLocation:     "redirect",
			wantErr:          false,
		},
		{
			name: "valid 301 redirect",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusMovedPermanently,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri",
				},
			},
			wantModifyHeader: true,
			want:             false,
			wantStatus:       http.StatusMovedPermanently,
			wantLocation:     "redirect",
			wantErr:          false,
		},
		{
			name: "error 410 redirect",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusGone,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri",
				},
			},
			wantModifyHeader: false,
			want:             false,
			wantStatus:       http.StatusGone,
			wantLocation:     "",
			wantErr:          false,
		},
		{
			name: "unsupported code",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: 999,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri",
				},
			},
			wantModifyHeader: false,
			want:             false,
			wantStatus:       http.StatusNotFound,
			wantLocation:     "",
			wantErr:          false,
		},
		{
			name: "redirect not in map",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: 999,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri2",
				},
			},
			wantModifyHeader: false,
			want:             true,
			wantStatus:       0,
			wantLocation:     "",
			wantErr:          true,
		},
		{
			name: "no redirect map",
			fields: fields{
				redirectDataMap: nil,
			},
			args: args{
				req: &http.Request{
					RequestURI: "the-uri2",
				},
			},
			wantModifyHeader: false,
			want:             true,
			wantStatus:       0,
			wantLocation:     "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redirector{
				logger:          flamingo.NullLogger{},
				redirectDataMap: tt.fields.redirectDataMap,
			}

			rwMock := &mocks.ResponseWriter{}
			header := http.Header{}
			if tt.wantModifyHeader {
				rwMock.On("Header").Return(header)
			}
			if !tt.wantErr {
				rwMock.On("WriteHeader", tt.wantStatus).Once()
			}

			got, err := r.TryServeHTTP(rwMock, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("redirector.TryServeHTTP() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("redirector.TryServeHTTP() = %v, want %v", got, tt.want)
			}

			if gotLocation := header.Get("Location"); tt.wantLocation != gotLocation {
				t.Errorf("redirector location = %v, want %v", gotLocation, tt.wantLocation)
			}

			rwMock.AssertExpectations(t)
		})
	}
}

func Test_redirector_Filter(t *testing.T) {
	type fields struct {
		redirectDataMap map[string]infrastructure.CsvContent
		redirect        *responderMocks.RedirectAware
		error           *responderMocks.ErrorAware
	}
	type args struct {
		req http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantChainCalled int32
	}{
		{
			name: "valid 302 redirect",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusFound,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
				redirect: withRedirect("RedirectURL", "redirect"),
			},
			args: args{
				req: http.Request{
					RequestURI: "the-uri",
				},
			},
		},
		{
			name: "valid 301 redirect",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusMovedPermanently,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
				redirect: withRedirect("RedirectPermanentURL", "redirect"),
			},
			args: args{
				req: http.Request{
					RequestURI: "the-uri",
				},
			},
		},
		{
			name: "error 410",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusGone,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
				error: withError("ErrorWithCode", errors.New("page is gone"), 410),
			},
			args: args{
				req: http.Request{
					RequestURI: "the-uri",
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				redirectDataMap: map[string]infrastructure.CsvContent{
					"the-uri": {
						HTTPStatusCode: http.StatusNotFound,
						OriginalPath:   "orig",
						RedirectTarget: "redirect",
					},
				},
				error: withError("ErrorNotFound", errors.New("page not found"), 0),
			},
			args: args{
				req: http.Request{
					RequestURI: "the-uri",
				},
			},
		},
		{
			name: "no file",
			fields: fields{
				redirectDataMap: nil,
			},
			args: args{
				req: http.Request{
					RequestURI: "the-uri",
				},
			},
			wantChainCalled: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &redirector{
				RedirectAware:   tt.fields.redirect,
				ErrorAware:      tt.fields.error,
				logger:          flamingo.NullLogger{},
				redirectDataMap: tt.fields.redirectDataMap,
			}

			rwMock := &mocks.ResponseWriter{}

			var chainCalled int32

			var f filterMocker = func(ctx context.Context, r *web.Request, w http.ResponseWriter) web.Response {
				atomic.AddInt32(&chainCalled, 1)
				return nil
			}

			chain := &router.FilterChain{
				Filters: []router.Filter{f},
			}

			r.Filter(
				context.Background(),
				web.RequestFromRequest(&tt.args.req, web.NewSession(&sessions.Session{})),
				rwMock,
				chain,
			)

			if tt.fields.redirect != nil {
				tt.fields.redirect.AssertExpectations(t)
			}
			if tt.fields.error != nil {
				tt.fields.error.AssertExpectations(t)
			}

			if chainCalled != tt.wantChainCalled {
				t.Errorf("number of chain.Next calls is %v, expected %v", chainCalled, tt.wantChainCalled)
			}
		})
	}
}

func withRedirect(method, location string) *responderMocks.RedirectAware {
	m := &responderMocks.RedirectAware{}
	m.On(method, location).Return(nil)

	return m
}

func withError(method string, err error, status int) *responderMocks.ErrorAware {
	m := &responderMocks.ErrorAware{}
	if status > 0 {
		m.On(method, mock.AnythingOfType("*context.emptyCtx"), err, status).Return(nil)
	} else {
		m.On(method, mock.AnythingOfType("*context.emptyCtx"), err).Return(nil)
	}

	return m
}

func Test_newRedirector(t *testing.T) {
	type args struct {
		redirectData *infrastructure.RedirectData
	}
	tests := []struct {
		name               string
		args               args
		wantedRedirectData map[string]infrastructure.CsvContent
	}{
		{
			name: "valid map",
			args: args{
				redirectData: infrastructure.NewRedirectData(
					&struct {
						RedirectsCsv string `inject:"config:redirects.csv"`
					}{RedirectsCsv: filepath.Join("testdata", "redirect.csv")},
					flamingo.NullLogger{}),
			},
			wantedRedirectData: map[string]infrastructure.CsvContent{
				"foo": {
					HTTPStatusCode: 301,
					OriginalPath:   "foo",
					RedirectTarget: "bar",
				},
				"baz": {
					HTTPStatusCode: 302,
					OriginalPath:   "baz",
					RedirectTarget: "bam",
				},
			},
		},
		{
			name: "loop map",
			args: args{
				redirectData: infrastructure.NewRedirectData(
					&struct {
						RedirectsCsv string `inject:"config:redirects.csv"`
					}{RedirectsCsv: filepath.Join("testdata", "redirect-loop.csv")},
					flamingo.NullLogger{}),
			},
			wantedRedirectData: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redAware := &responderMocks.RedirectAware{}
			errAware := &responderMocks.ErrorAware{}
			logger := flamingo.NullLogger{}

			got := newRedirector(redAware, errAware, logger, tt.args.redirectData)

			if !reflect.DeepEqual(tt.wantedRedirectData, got.redirectDataMap) {
				t.Errorf("unexpected redirect data map: got %#v, want %#v", got.redirectDataMap, tt.wantedRedirectData)
			}

		})
	}
}
