package prefixrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontRouter(t *testing.T) {
	var fr = NewFrontRouter()

	createRouterHandler := func(prefix string) (string, routerHandler) {
		return prefix, routerHandler{
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/", "/test":
					_, err := w.Write([]byte(fmt.Sprintf("Found: %s %s", prefix, r.URL.Path)))
					assert.NoError(t, err)
				default:
					_, err := w.Write([]byte(fmt.Sprintf("Not found: %s %s", prefix, r.URL.Path)))
					assert.NoError(t, err)
				}
			}),
		}
	}

	fr.Add(createRouterHandler("/prefix1"))

	fr.Add(createRouterHandler("/prefix2"))

	fr.Add(createRouterHandler("test.com/prefix1"))

	fr.SetFinalFallbackHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Default"))
	}))

	t.Run("Request Routing", func(t *testing.T) {
		t.Run("Should Match Host before Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			request.Host = "test.com"

			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Found: test.com/prefix1 /test", string(body))
		})

		t.Run("Should Match Host before Prefix but 404 for double prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/prefix1/test", nil)
			request.Host = "test.com"

			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Not found: test.com/prefix1 /prefix1/test", string(body))
		})

		t.Run("Should Match Prefix after Host", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Found: /prefix1 /test", string(body))
		})

		t.Run("Should Match Prefix after Host but 404 for double prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Not found: /prefix1 /prefix1/test", string(body))
		})

		t.Run("Should Match Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Found: /prefix2 /test", string(body))
		})

		t.Run("Should Match Prefix but 404 for double prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2/prefix2/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Not found: /prefix2 /prefix2/test", string(body))
		})

		t.Run("Should Match just the Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Found: /prefix2 /", string(body))
		})

		t.Run("Should Match just the Prefix but 404 for double prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2/prefix2", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Not found: /prefix2 /prefix2", string(body))
		})

		t.Run("Should have a default", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/UNKNOWN/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "Default", string(body))
		})

		t.Run("Should use 404 if no default is set", func(t *testing.T) {
			fr = NewFrontRouter()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, 404)
			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, "", string(body))
		})

		t.Run("Primary handler", func(t *testing.T) {
			fr = NewFrontRouter()
			fr.primaryHandlers = []OptionalHandler{
				&rootRedirectHandler{
					redirectTarget: "http://flamingo.me/flamingo",
				},
			}
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, http.StatusTemporaryRedirect)
			assert.Equal(t, recorder.Header().Get("Location"), "http://flamingo.me/flamingo")
		})

		t.Run("Fallback handler", func(t *testing.T) {
			fr = NewFrontRouter()
			fr.fallbackHandlers = []OptionalHandler{
				&rootRedirectHandler{
					redirectTarget: "http://flamingo.me/flamingo",
				},
			}
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, http.StatusTemporaryRedirect)
			assert.Equal(t, recorder.Header().Get("Location"), "http://flamingo.me/flamingo")
		})
	})
}
