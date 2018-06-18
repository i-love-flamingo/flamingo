package prefixrouter

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FrontRouter", func() {
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

	Context("Request Routing", func() {
		It("Should Match Host before Prefix", func() {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			Expect(ioutil.ReadAll(recorder.Result().Body)).To(Equal([]byte(`Host1, Prefix 1`)))
		})

		It("Should Match Prefix", func() {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2/test", nil)
			fr.ServeHTTP(recorder, request)

			Expect(ioutil.ReadAll(recorder.Result().Body)).To(Equal([]byte(`Prefix 2`)))
		})

		It("Should Match just the Prefix", func() {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix2", nil)
			fr.ServeHTTP(recorder, request)

			Expect(ioutil.ReadAll(recorder.Result().Body)).To(Equal([]byte(`Prefix 2`)))
		})

		It("Should have a default", func() {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/UNKNOWN/test", nil)
			fr.ServeHTTP(recorder, request)

			Expect(ioutil.ReadAll(recorder.Result().Body)).To(Equal([]byte(`Default`)))
		})

		It("Should use 404 if no default is set", func() {
			fr = NewFrontRouter()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/prefix1/test", nil)
			fr.ServeHTTP(recorder, request)

			Expect(ioutil.ReadAll(recorder.Result().Body)).To(Equal([]byte(``)))
			Expect(recorder.Result().StatusCode).To(Equal(404))
		})
	})
})

func TestFrontRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FrontRouter Suite")
}
