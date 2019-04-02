package router

import (
	"go.aoe.com/flamingo/framework/web"

	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testController(web.Context) web.Response {
	return &web.ContentResponse{}
}

var _ = Describe("Registry Test", func() {
	Context("Utils", func() {
		Context("parseHandler", func() {
			It("should treat empty params properly", func() {
				var handler = parseHandler("foo.bar")
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(BeEmpty())
			})

			It("should treat params properly", func() {
				var handler = parseHandler("foo.bar(foo, bar)")
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(2))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: false, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("bar", &param{optional: false, value: ""}))
			})

			It("should treat optional params properly", func() {
				var handler = parseHandler("foo.bar(foo?, bar?)")
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(2))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: true, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("bar", &param{optional: true, value: ""}))
			})

			It("should treat hardcoded params properly", func() {
				var handler = parseHandler(`foo.bar(foo="bar", x="y")`)
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(2))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: false, value: "bar"}))
				Expect(handler.params).To(HaveKeyWithValue("x", &param{optional: false, value: "y"}))
			})

			It("should treat default value params properly", func() {
				var handler = parseHandler(`foo.bar(foo?="bar")`)
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(1))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: true, value: "bar"}))
			})

			It("should treat complexer params properly", func() {
				var handler = parseHandler(`foo.bar(a, b?, x="a", y ,z, foo ?= "bar")`)
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(6))
				Expect(handler.params).To(HaveKeyWithValue("a", &param{optional: false, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("b", &param{optional: true, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("x", &param{optional: false, value: "a"}))
				Expect(handler.params).To(HaveKeyWithValue("y", &param{optional: false, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("z", &param{optional: false, value: ""}))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: true, value: "bar"}))
			})

			It("should treat escaped values properly", func() {
				var handler = parseHandler(`foo.bar(foo?="\"bar")`)
				Expect(handler.handler).To(Equal("foo.bar"))
				Expect(handler.params).To(HaveLen(1))
				Expect(handler.params).To(HaveKeyWithValue("foo", &param{optional: true, value: `"bar`}))
			})
		})
	})

	Context("API", func() {
		var registry = NewRegistry()
		It("Should create a new router", func() {
			Expect(registry).ToNot(BeNil())
		})

		It("Should allow controller registration", func() {
			registry.Handle("page.view", testController)
		})

		It("Should make path registration easy", func() {
			registry.Route("/", `page.view(page="home")`)                 // hardcoded
			registry.Route("/page/:page", `page.view(page)`)              // extract via param
			registry.Route("/homepage/:page", `page.view(page?="home2")`) // extract via param, default value
			registry.Route("/page", `page.view(page="page")`)             // extract from GET
			registry.Route("/page2", `page.view(page?="page2")`)          // extract from GET
			registry.Route("/mustget", `page.view(page)`)                 // extract from GET
		})

		It("Should reverse routes properly", func() {
			Expect(registry.Reverse("page.view", map[string]string{"page": "home"})).To(Equal("/"))
			Expect(registry.Reverse("page.view", map[string]string{})).To(Equal("/homepage/home2"))
			Expect(registry.Reverse("page.view", map[string]string{"page": "foo"})).To(Equal("/page/foo"))
		})

		It("Should match paths", func() {
			By("route parameters")
			controller, params := registry.Match("/homepage/home2")
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "home2"))

			By("hardcoded")
			controller, params = registry.Match("/")
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "home"))

			By("optional")
			controller, params = registry.Match("/page2")
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "page2"))

			By("optional")
			controller, params = registry.Match("/page")
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "page"))
		})

		It("Should match HTTP Requests", func() {
			By("Default")
			request, _ := http.NewRequest("GET", "/page2", nil)
			controller, params, _ := registry.MatchRequest(request)
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "page2"))

			By("GET Parameter")
			request, _ = http.NewRequest("GET", "/page2?page=foo", nil)
			controller, params, _ = registry.MatchRequest(request)
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "foo"))

			By("Missing GET Parameter")
			request, _ = http.NewRequest("GET", "/mustget", nil)
			controller, params, _ = registry.MatchRequest(request)
			Expect(controller).To(BeNil())
			Expect(params).To(BeNil())

			By("Mandatory GET Parameter")
			request, _ = http.NewRequest("GET", "/mustget?page=foo", nil)
			controller, params, _ = registry.MatchRequest(request)
			Expect(controller).ToNot(BeNil())
			Expect(params).To(HaveLen(1))
			Expect(params).To(HaveKeyWithValue("page", "foo"))
		})

		It("should render get requests if possible", func() {
			registry.Handle("page.get", testController)
			registry.Route("/path_mustget", `page.get(page)`)

			registry.Handle("page.get2", testController)
			registry.Route("/path_mustget2", `page.get2(page?="test")`)

			path, err := registry.Reverse("page.get", map[string]string{"page": "test"})
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal("/path_mustget?page=test"))

			path, err = registry.Reverse("page.get2", map[string]string{"page": "test"})
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal("/path_mustget2"))

			path, err = registry.Reverse("page.get2", map[string]string{"page": "nottest"})
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal("/path_mustget2?page=nottest"))
		})
	})
})
