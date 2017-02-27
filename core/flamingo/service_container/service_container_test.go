package service_container

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	DI_root struct{}

	DI_wants_root struct {
		Root *DI_root `inject:""`
		mark string
	}

	DI_test struct {
		*DI_wants_root   `inject:""`
		Test             *DI_root `inject:"root"`
		postinjectmarker bool
	}
)

func (d *DI_test) PostInject() {
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

		var test, test_named = new(DI_test), new(DI_root)

		Describe("Dependency Injection behaviour", func() {
			Context("When Working with test dependency Tree", func() {
				BeforeEach(func() {
					test, test_named = new(DI_test), new(DI_root)

					serviceContainer.Register(test, "tag.1")
					serviceContainer.Register(test_named)
					serviceContainer.RegisterNamed("root", test_named, "tag.1", "tag.2")
				})

				It("Should resove the whole graph", func() {
					Expect(test.Root).To(Equal(test_named))
					Expect(test.Test).To(Equal(test_named))
					Expect(test.DI_wants_root).To(BeAssignableToTypeOf(new(DI_wants_root)))
				})

				It("Should have tags", func() {
					Expect(serviceContainer.GetByTag("tag.1")).To(ConsistOf(test, test_named))
					Expect(serviceContainer.GetByTag("tag.2")).To(ConsistOf(test_named))
					Expect(serviceContainer.GetByTag("tag.3")).To(BeEmpty())
				})

				It("Should call PostInject methods", func() {
					Expect(test.postinjectmarker).To(BeTrue())
				})
			})

			Context("When dealing with already existing objects", func() {
				var di_old, di_new = &DI_wants_root{mark: "old"}, &DI_wants_root{mark: "new"}

				BeforeEach(func() {
					test, test_named = new(DI_test), new(DI_root)
					di_old, di_new = &DI_wants_root{mark: "old"}, &DI_wants_root{mark: "new"}
					serviceContainer.Register(test)
					serviceContainer.RegisterNamed("root", test_named)
					serviceContainer.Register(di_old)
					serviceContainer.Register(di_new)
				})

				It("Should have replaced di_old with the newer di_new object", func() {
					Expect(test.DI_wants_root).ToNot(Equal(di_old))
					Expect(test.DI_wants_root).To(Equal(di_new))
					Expect(test.DI_wants_root.mark).To(Equal("new"))
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
