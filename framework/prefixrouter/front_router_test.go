package prefixrouter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrontRouter(t *testing.T) {
	var fr = NewFrontRouter()

	fr.Add("/prefix1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Prefix 1"))
	}))

	fr.Add("/prefix2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Prefix 2"))
	}))

	fr.Add("example.com/prefix1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Host1, Prefix 1"))
	}))

	fr.SetFinalFallbackHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Default"))
	}))

	t.Run("Request Routing", func(t *testing.T) {
		t.Run("Should Match Host before Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com/prefix1/test", nil)

			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Host1, Prefix 1`))
		})

		t.Run("Should Match Host before Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example2.com/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Prefix 1`))
		})

		t.Run("Should Match Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com/prefix2/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Prefix 2`))
		})

		t.Run("Should Match just the Prefix", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com/prefix2", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Prefix 2`))
		})

		t.Run("Should have a default", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com/UNKNOWN/test", nil)
			fr.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(`Default`))
		})

		t.Run("Should use 404 if no default is set", func(t *testing.T) {
			fr = NewFrontRouter()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "http://example.com/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Result().StatusCode, 404)
			body, err := ioutil.ReadAll(recorder.Result().Body)
			assert.NoError(t, err)
			assert.Equal(t, body, []byte(``))
		})
	})
}
