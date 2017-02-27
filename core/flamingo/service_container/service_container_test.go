package service_container

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	DIroot struct{}

	DIwantsRoot struct {
		Root *DIroot `inject:""`
		mark string
	}

	DItest struct {
		*DIwantsRoot     `inject:""`
		Test             *DIroot `inject:"root"`
		postinjectmarker bool
	}
)

func (d *DItest) PostInject() {
	d.postinjectmarker = true
}

var _ = Describe("ServiceContainerPackage", func() {
	Describe("RegisterFunc", func() {
		var testfunc = RegisterFunc(func(r *ServiceContainer) {})

		It("Should have a name", func() {
			name, err := testfunc.MarshalText()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(name)).To(Equal("flamingo/core/flamingo/service_container.glob..func1.1.1"))
		})
	})

	Describe("ServiceContainer", func() {
		var serviceContainer *ServiceContainer

		BeforeEach(func() {
			serviceContainer = New()
		})

		JustBeforeEach(func() {
			serviceContainer.Resolve()
		})

		Context("When newly creatd", func() {
			It("Should be a new ServiceContainer", func() {
				Expect(serviceContainer).To(BeAssignableToTypeOf(&ServiceContainer{}))
				Expect(serviceContainer).ToNot(BeNil())
			})
		})

		var test, testNamed = new(DItest), new(DIroot)

		Describe("Dependency Injection behaviour", func() {
			Context("When Working with test dependency Tree", func() {
				BeforeEach(func() {
					test, testNamed = new(DItest), new(DIroot)

					serviceContainer.Register(test, "tag.1")
					serviceContainer.Register(testNamed)
					serviceContainer.RegisterNamed("root", testNamed, "tag.1", "tag.2")
				})

				It("Should resove the whole graph", func() {
					Expect(test.Root).To(Equal(testNamed))
					Expect(test.Test).To(Equal(testNamed))
					Expect(test.DIwantsRoot).To(BeAssignableToTypeOf(new(DIwantsRoot)))
				})

				It("Should have tags", func() {
					Expect(serviceContainer.GetByTag("tag.1")).To(ConsistOf(test, testNamed))
					Expect(serviceContainer.GetByTag("tag.2")).To(ConsistOf(testNamed))
					Expect(serviceContainer.GetByTag("tag.3")).To(BeEmpty())
				})

				It("Should call PostInject methods", func() {
					Expect(test.postinjectmarker).To(BeTrue())
				})
			})

			Context("When dealing with already existing objects", func() {
				var diOld, diNew = &DIwantsRoot{mark: "old"}, &DIwantsRoot{mark: "new"}

				BeforeEach(func() {
					test, testNamed = new(DItest), new(DIroot)
					diOld, diNew = &DIwantsRoot{mark: "old"}, &DIwantsRoot{mark: "new"}
					serviceContainer.Register(test)
					serviceContainer.RegisterNamed("root", testNamed)
					serviceContainer.Register(diOld)
					serviceContainer.Register(diNew)
				})

				It("Should have replaced diOld with the newer diNew object", func() {
					Expect(test.DIwantsRoot).ToNot(Equal(diOld))
					Expect(test.DIwantsRoot).To(Equal(diNew))
					Expect(test.DIwantsRoot.mark).To(Equal("new"))
				})
			})
		})

		Describe("Private Dependencies on runtime", func() {
			Context("When in need of dynamically request-scoped dependencies", func() {
				type (
					Something  struct{}
					PerRequest struct {
						Something *Something `inject:""`
					}
					PerRequest2 struct {
						Something *Something `inject:""`
					}
					Global struct{}
				)

				var something = new(Something)
				var global = new(Global)

				BeforeEach(func() {
					serviceContainer.Register(something)
					serviceContainer.Register(global)

					serviceContainer.Register(PerRequest{}, "tag")
					serviceContainer.Register(PerRequest2{}, "tag")
				})

				It("Should give me new objects as needed with proper dependencies resolved", func() {
					var e = serviceContainer.Create(PerRequest{}).(*PerRequest)
					Expect(e).To(BeAssignableToTypeOf(&PerRequest{}))
					Expect(e.Something).ToNot(BeNil())
					Expect(e.Something).To(Equal(something))

					var e2 = serviceContainer.Create(e).(*PerRequest)
					Expect(e2).To(BeAssignableToTypeOf(&PerRequest{}))
					Expect(e2.Something).ToNot(BeNil())
					Expect(e2.Something).To(Equal(something))
				})
			})
		})
	})
})

func TestServiceContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServiceContainer Suite")
}
