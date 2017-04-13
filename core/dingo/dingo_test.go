package dingo

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	Interface interface {
		Test() int
	}

	InterfaceSub Interface

	InterfaceImpl1 struct {
		i int
	}

	InterfaceImpl2 struct {
		i int
	}

	DepTest struct {
		Iface  Interface `inject:""`
		Iface2 Interface `inject:"test"`

		IfaceProvider func() Interface `inject:""`
		IfaceProvided Interface        `inject:"provider"`
		IfaceInstance Interface        `inject:"instance"`
	}

	TestSingleton struct {
		i int
	}

	TestModule struct{}
)

func init() {
	var _ Interface = &InterfaceImpl1{}
	var _ Interface = &InterfaceImpl2{}
}

func InterfaceProvider() Interface {
	return new(InterfaceImpl1)
}

func (tm *TestModule) Configure(injector *Injector) {
	injector.Bind((*Interface)(nil)).To((*InterfaceSub)(nil))
	injector.Bind((*InterfaceSub)(nil)).To(InterfaceImpl1{})
	injector.Bind((*Interface)(nil)).AnnotatedWith("test").To(InterfaceImpl2{})

	injector.Bind((*Interface)(nil)).AnnotatedWith("provider").ToProvider(InterfaceProvider)
	injector.Bind((*Interface)(nil)).AnnotatedWith("instance").ToInstance(new(InterfaceImpl2))

	injector.Bind(TestSingleton{}).In(Singleton)
}

func (if1 *InterfaceImpl1) Test() int {
	return 1
}

func (if2 *InterfaceImpl2) Test() int {
	return 2
}

var _ = Describe("Dingo Test", func() {
	Context("Simple resolve", func() {
		It("Should resolve dependencies on request", func() {
			injector := NewInjector(new(TestModule))

			var iface Interface
			iface = injector.GetInstance(new(Interface)).(Interface)

			Expect(iface.Test()).To(Equal(1))

			var dt DepTest = *injector.GetInstance(new(DepTest)).(*DepTest)

			Expect(dt.Iface.Test()).To(Equal(1))
			Expect(dt.Iface2.Test()).To(Equal(2))

			var dt2 DepTest
			injector.RequestInjection(&dt2)

			Expect(dt2.Iface.Test()).To(Equal(1))
			Expect(dt2.Iface2.Test()).To(Equal(2))

			log.Printf("%#v\n", dt)

			Expect(dt.IfaceProvided.Test()).To(Equal(1))
			Expect(dt.IfaceInstance.Test()).To(Equal(2))
			Expect(dt.IfaceProvider().Test()).To(Equal(1))
		})

		It("Should resolve scopes", func() {
			injector := NewInjector(new(TestModule))

			log.Printf("%p\n", injector.GetInstance(TestSingleton{}))
			log.Printf("%p\n", injector.GetInstance(TestSingleton{}))

			panic(1)
		})
	})
})

func TestServiceContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dingo Suite")
}
