package dependencyinjection

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	// App defines an example app
	App struct {
		*Router         `inject:""`
		CtxFactory      CtxFactory `inject:""`
		NamedCtxFactory CtxFactory `inject:"namedctxfactory"`
		TestRouter      Router     `inject:"inline"`
		PrivateRouter   *Router    `inject:"private"`
		NamedRouter     *Router    `inject:"namedrouter"`
		TestSubscriber  Subscriber `inject:""`
	}

	// Router is a dependency of App
	Router struct {
		*Test `inject:""`
		*App  `inject:"app"`
	}

	// CtxFactory creates ctx-instances
	CtxFactory func(arg1, arg2 string) Ctx

	// Ctx is something with an event router
	Ctx interface {
		EventRouter
	}

	// CtxImpl is an implementation of Ctx with it's own dependencies
	CtxImpl struct {
		EventRouter     `inject:"private"`
		SuperDependency `inject:"private"`
	}

	// EventRouter is an example event system
	EventRouter interface {
		Dispatch(event string)
		Subscriber() []Subscriber
	}

	// EventRouterImpl implements the EventRouter
	EventRouterImpl struct {
		subscriber []Subscriber
	}

	// Subscriber is something which can be used by EventRouter
	Subscriber interface {
		Events() []string
	}

	// SubscriberImpl is an example Subscriber
	SubscriberImpl struct {
		Name string
	}

	// SuperDependency is something which is also a Subscriber
	SuperDependency interface {
		Test() string
	}

	// SuperDependencyImpl is an implementation of SuperDependency
	SuperDependencyImpl struct {
		nothing string // prevent compiler optimization
	}

	// Test is just a test struct
	Test struct {
		nothing string // prevent compiler optimization
	}
)

// CompilerPass to retrieve subscriber
func (ev *EventRouterImpl) CompilerPass(c *Container) {
	for _, o := range c.GetTagged("subscriber") {
		ev.subscriber = append(ev.subscriber, o.Value.(Subscriber))
	}
}

// Dispatch to implement EventRouter
func (ev *EventRouterImpl) Dispatch(event string) {}

// Subscriber to test and implement EventRouter
func (ev *EventRouterImpl) Subscriber() []Subscriber {
	return ev.subscriber
}

// Events to implement Subscriber
func (si *SubscriberImpl) Events() []string {
	return []string{
		"test",
		"test2",
	}
}

// Test to implement SuperDependency
func (sdi *SuperDependencyImpl) Test() string {
	return "SuperDependencyImpl"
}

// Events to implement Subscriber
func (sdi *SuperDependencyImpl) Events() []string {
	return []string{
		"test",
	}
}

var _ = Describe("DependencyInjection Package", func() {
	Describe("Example App Behaviour", func() {
		var container = NewContainer()

		// CtxFactory
		container.Register(CtxFactory(func(arg1, arg2 string) Ctx { return new(CtxImpl) }))

		// Named CtxFactory
		container.RegisterNamed("namedctxfactory", CtxFactory(func(arg1, arg2 string) Ctx { return new(CtxImpl) }))

		// EventRouter Factory
		container.RegisterFactory(func() EventRouter { return new(EventRouterImpl) })

		// Subscriber+SuperDependency
		container.RegisterFactory(func() SuperDependency { return new(SuperDependencyImpl) }, "subscriber")

		// Named Router
		container.RegisterNamed("namedrouter", new(Router))

		// Subscriber
		var testsubscriber = new(SubscriberImpl)
		container.Register(testsubscriber, "subscriber")

		Context("App will be resolved", func() {
			container.RegisterNamed("app", new(App), "app-tag")
			app := container.Get("app").(*App)

			It("Should have all dependencies properly resolved", func() {
				Expect(app).ToNot(BeNil())
				Expect(app).To(Equal(container.Get("app").(*App)))

				Expect(app.Router).ToNot(BeNil())
				Expect(app.Router == app.PrivateRouter).ToNot(BeTrue())
				Expect(app.Router.App).To(Equal(app))

				Expect(app.TestSubscriber).To(Equal(testsubscriber))

				Expect(app.CtxFactory).ToNot(BeNil())
				Expect(app.NamedCtxFactory).ToNot(BeNil())
				Expect(app.CtxFactory("a1", "a2")).ToNot(BeNil())
				Expect(app.NamedCtxFactory("a1", "a2")).ToNot(BeNil())

				ctx1 := app.CtxFactory("a1", "a2").(Ctx)
				ctx2 := app.CtxFactory("a1", "a2").(Ctx)
				Expect(ctx1).To(BeAssignableToTypeOf(&CtxImpl{}))
				Expect(ctx2).To(BeAssignableToTypeOf(&CtxImpl{}))
				Expect(ctx1 == ctx2).ToNot(BeTrue())

				Expect(ctx1.(*CtxImpl).SuperDependency).ToNot(BeNil())
				Expect(ctx1.(*CtxImpl).SuperDependency).To(Equal(ctx2.(*CtxImpl).SuperDependency))
				Expect(ctx1.(*CtxImpl).SuperDependency.(*SuperDependencyImpl) == ctx2.(*CtxImpl).SuperDependency.(*SuperDependencyImpl)).ToNot(BeTrue())

				Expect(ctx1.Subscriber()).To(ConsistOf(new(SubscriberImpl), new(SuperDependencyImpl)))
			})
		})
	})

	Describe("Error cases", func() {
		It("Should panic for concrete types", func() {
			Expect(func() {
				var container = NewContainer()
				container.RegisterNamed("string", "string")
				container.Get("string")
			}).To(Panic())
		})

		It("Should panic for unknown services", func() {
			Expect(func() {
				var container = NewContainer()
				container.Get("unknown")
			}).To(Panic())
		})

		It("Should panic for unknown factories", func() {
			Expect(func() {
				var container = NewContainer()
				container.RegisterNamed("test", new(CtxImpl))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for private inline structs", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					Test `inject:"private"`
				}
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for unknown functions", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					CtxFactory `inject:""`
				}
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for more than one possible interface candidates", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					Subscriber Subscriber `inject:""`
				}
				container.Register(new(SubscriberImpl))
				container.Register(new(SuperDependencyImpl))
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for no assignable interfaces", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					Subscriber Subscriber `inject:""`
				}
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for not injectable values", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					Test string `inject:""`
				}
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for not known named dependencies", func() {
			Expect(func() {
				var container = NewContainer()
				type test struct {
					SubscriberImpl `inject:"subscriber"`
				}
				container.RegisterNamed("test", new(test))
				container.Get("test")
			}).To(Panic())
		})

		It("Should panic for more than one result for functions", func() {
			Expect(func() {
				var container = NewContainer()
				container.Register(func() (string, string) { return "a", "b" })
			}).To(Panic())

			Expect(func() {
				var container = NewContainer()
				container.RegisterNamed("test", func() (string, string) { return "a", "b" })
			}).To(Panic())
		})
	})

	Context("Simple resolve", func() {
		It("Should resolve dependencies on request", func() {
			var ctximpl = new(CtxImpl)
			var container = NewContainer()

			container.RegisterFactory(func() EventRouter { return new(EventRouterImpl) })
			container.RegisterFactory(func() SuperDependency { return new(SuperDependencyImpl) })

			Expect(ctximpl.EventRouter).To(BeNil())
			Expect(ctximpl.SuperDependency).To(BeNil())

			container.Resolve(ctximpl)

			Expect(ctximpl.EventRouter).ToNot(BeNil())
			Expect(ctximpl.SuperDependency).ToNot(BeNil())
		})
	})
})

func TestServiceContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DependencyInjection Suite")
}
