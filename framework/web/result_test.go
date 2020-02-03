package web

import (
	"errors"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/flamingo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

var mockTemplateEngine = &mocks.TemplateEngine{}

func getResponder(withRouter bool) Responder {
	var router *Router

	if withRouter == true {
		router = &Router{
			base: &url.URL{Path: "base_path"},
		}
	}

	responder := Responder{}

	responder.Inject(
		router,
		flamingo.NullLogger{},
		&struct {
			Engine                flamingo.TemplateEngine `inject:",optional"`
			Debug                 bool                    `inject:"config:flamingo.debug.mode"`
			TemplateForbidden     string                  `inject:"config:flamingo.template.err403"`
			TemplateNotFound      string                  `inject:"config:flamingo.template.err404"`
			TemplateUnavailable   string                  `inject:"config:flamingo.template.err503"`
			TemplateErrorWithCode string                  `inject:"config:flamingo.template.errWithCode"`
		}{
			Engine:                mockTemplateEngine,
			Debug:                 false,
			TemplateForbidden:     "403_template",
			TemplateNotFound:      "404_template",
			TemplateUnavailable:   "503_template",
			TemplateErrorWithCode: "withErrorCode_template",
		})
	return responder
}

func TestServerErrorWithCodeAndTemplate(t *testing.T) {
	testErr := errors.New("test error")

	t.Run("Router is not nil", func(t *testing.T) {
		responder := getResponder(true)
		require.NotNil(t, responder.router)

		actual := responder.ServerErrorWithCodeAndTemplate(testErr, "403_template", 403)
		assert.Equal(t, 403, int(actual.Response.Status))
		assert.Equal(t, "403_template", actual.Template)

		assert.Equal(t, "base_path", actual.Data.(map[string]interface{})["base"].(string))
		assert.Equal(t, 403, int(actual.Data.(map[string]interface{})["code"].(uint)))
		assert.Equal(t, testErr.Error(), actual.Data.(map[string]interface{})["error"].(string))
	})

	t.Run("Router is nil", func(t *testing.T) {
		responder := getResponder(false)
		require.Nil(t, responder.router)

		actual := responder.ServerErrorWithCodeAndTemplate(testErr, "403_template", 403)
		assert.Equal(t, 403, int(actual.Response.Status))
		assert.Equal(t, "403_template", actual.Template)

		assert.Equal(t, "", actual.Data.(map[string]interface{})["base"].(string))
		assert.Equal(t, 403, int(actual.Data.(map[string]interface{})["code"].(uint)))
		assert.Equal(t, testErr.Error(), actual.Data.(map[string]interface{})["error"].(string))
	})
}
