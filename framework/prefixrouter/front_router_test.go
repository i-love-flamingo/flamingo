package prefixrouter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontRouter(t *testing.T) {
	var fr = NewFrontRouter()

	fr.Add("/prefix1", routerHandler{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Prefix 1"))
	})})

	fr.Add("/prefix2", routerHandler{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Prefix 2"))
	})})

	fr.Add("/prefix22", routerHandler{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Prefix 22"))
	})})

	fr.Add("test.com/prefix1", routerHandler{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Host1, Prefix 1"))
	})})

	fr.Add("test.com/prefix11", routerHandler{handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Host1, Prefix 11"))
	})})

	fr.SetFinalFallbackHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Default"))
	}))

	t.Run("Request Routing", func(t *testing.T) {
		t.Run("Should Match Host before Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			request.Host = "test.com"

			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, []byte(`Host1, Prefix 1`), body)
		})

		t.Run("Should Match Host before Prefix and longer prefix first", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix11/test", nil)
			request.Host = "test.com"

			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, []byte(`Host1, Prefix 11`), body)
		})

		t.Run("Should Match Prefix without Host", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, []byte(`Prefix 1`), body)
		})

		t.Run("Should Match Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, []byte(`Prefix 2`), body)
		})

		t.Run("Should Match just the Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2", nil)
			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Prefix 2`))
		})

		t.Run("Should Match longer Prefix first", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix22/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Prefix 22`))
		})

		t.Run("Should have a default", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/UNKNOWN/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Default`))
		})

		t.Run("Should use 404 if no default is set", func(t *testing.T) {
			emptyFR := NewFrontRouter()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			emptyFR.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, 404)
			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(``))
		})

		t.Run("Malformed url should lead to 404", func(t *testing.T) {
			// only path
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1", nil)
			request.RequestURI = "/prefix1%20HTTP%2F1.1%0D%0ASomeheader:%20value%0D%0A%0D%0A/test"
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, 404)
			body, err := io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(``))

			// host + path
			recorder = httptest.NewRecorder()
			request = httptest.NewRequest("GET", "/prefix1", nil)
			request.Host = "test.com"
			request.RequestURI = "/prefix1%20HTTP%2F1.1%0D%0ASomeheader:%20value%0D%0A%0D%0A/test"
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, 404)
			body, err = io.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(``))
		})
	})
}

func TestFrontRouter_AddDuplicate(t *testing.T) {
	var fr = NewFrontRouter()

	assert.PanicsWithValue(t,
		`prefixrouter: duplicate handler registration on prefix "/prefix" from areas "area1" and "area2"`,
		func() {
			fr.Add("/prefix", routerHandler{area: "area1"})
			fr.Add("/prefix", routerHandler{area: "area2"})
		})
}
