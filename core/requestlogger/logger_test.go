package requestlogger

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

func TestLogger(t *testing.T) {
	logger := new(logger)

	logSink := new(bytes.Buffer)
	logger.Inject(&flamingo.StdLogger{Logger: *log.New(logSink, "", 0)}, nil)

	recorder := httptest.NewRecorder()
	request := web.CreateRequest(httptest.NewRequest(http.MethodPost, "/test", nil), nil)
	request.Request().Header.Set("X-Business-ID", "business-id")
	request.Request().Header.Set("Referer", "https://example.com/")

	responder := new(web.Responder).Inject(&web.Router{}, flamingo.NullLogger{}, &struct {
		Engine                flamingo.TemplateEngine "inject:\",optional\""
		Debug                 bool                    "inject:\"config:flamingo.debug.mode\""
		TemplateForbidden     string                  "inject:\"config:flamingo.template.err403\""
		TemplateNotFound      string                  "inject:\"config:flamingo.template.err404\""
		TemplateUnavailable   string                  "inject:\"config:flamingo.template.err503\""
		TemplateErrorWithCode string                  "inject:\"config:flamingo.template.errWithCode\""
	}{})

	tests := []struct {
		testCase string
		response web.Result
		regex    string
	}{
		{testCase: "regular response", response: &web.Response{}, regex: "^info: \x1b\\[32mPOST /test 200: 0b in \\d+(\\.\\d+)?..\x1b\\[0m$"},
		{testCase: "http status ok", response: &web.Response{Status: http.StatusOK}, regex: "^info: \x1b\\[32mPOST /test 200: 0b in \\d+(\\.\\d+)?..\x1b\\[0m$"},
		{testCase: "url redirect", response: &web.URLRedirectResponse{Response: web.Response{Status: http.StatusSeeOther}, URL: &url.URL{Path: "/foo"}}, regex: "^info: \x1b\\[34mPOST /test 303: 0b in \\d+(\\.\\d+)?.. \\(-> /foo\\)\x1b\\[0m$"},
		{testCase: "route redirect", response: responder.RouteRedirect("/route", nil), regex: "^info: \x1b\\[34mPOST /test 303: 0b in \\d+(\\.\\d+)?.. \\(-> /route\\)\x1b\\[0m$"},
		{testCase: "http notfound", response: &web.Response{Status: http.StatusNotFound}, regex: "^info: \x1b\\[33mPOST /test 404: 0b in \\d+(\\.\\d+)?..\x1b\\[0m$"},
		{testCase: "internal server errror", response: &web.Response{Status: http.StatusInternalServerError}, regex: "^info: \x1b\\[31mPOST /test 500: 0b in \\d+(\\.\\d+)?..\x1b\\[0m$"},
		{testCase: "error response", response: &web.ServerErrorResponse{Error: fmt.Errorf("test error")}, regex: "^info: \x1b\\[31mPOST /test 500: 5b in \\d+(\\.\\d+)?.. \\(Error: test error\\)\x1b\\[0m$"},
		{testCase: "unknown status code", response: &web.Response{Status: 1}, regex: "^info: \x1b\\[90mPOST /test 1: 0b in \\d+(\\.\\d+)?..\x1b\\[0m$"},
	}

	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			logSink.Reset()
			assert.NoError(t, logger.Filter(context.Background(), request, nil, web.NewFilterChain(func(ctx context.Context, req *web.Request, w http.ResponseWriter) web.Result {
				return test.response
			})).Apply(context.Background(), recorder))
			logs := strings.Split(logSink.String(), "\n")
			require.Len(t, logs, 3)
			assert.Regexp(t, "^WithFields map\\[accesslog:1 businessId:business-id client_ip:192.0.2.1:1234 referer:https://example.com/ response_code:\\d+ response_time:\\d+(\\.\\d+)?..]$", logs[0])
			assert.Regexp(t, test.regex, logs[1])
		})
	}
}

func TestHumanBytes(t *testing.T) {
	assert.Equal(t, "100b", humanBytes(100))
	assert.Equal(t, "100kb", humanBytes(100000))
}
