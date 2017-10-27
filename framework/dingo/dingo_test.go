package dingo

//import (
//	"testing"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//)

type (
	Interface interface {
		Test() int
	}

	InterfaceSub Interface

	InterfaceImpl1 struct {
		i   int
		foo string
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

	PreTestModule struct{}
)

func InterfaceProvider(str string) Interface {
	return &InterfaceImpl1{foo: str}
}

func (ptm *PreTestModule) Configure(injector *Injector) {
	injector.Bind((*string)(nil)).ToInstance("Hello World")
}

func (tm *TestModule) Configure(injector *Injector) {
	injector.Bind((*Interface)(nil)).To((*InterfaceSub)(nil))
	injector.Bind((*InterfaceSub)(nil)).To(InterfaceImpl1{})
	injector.Bind((*Interface)(nil)).AnnotatedWith("test").To(InterfaceImpl2{})

	injector.Bind((*Interface)(nil)).AnnotatedWith("provider").ToProvider(InterfaceProvider)
	injector.Bind((*Interface)(nil)).AnnotatedWith("instance").ToInstance(new(InterfaceImpl2))

	injector.Bind(TestSingleton{}).AsEagerSingleton()
}

func (if1 *InterfaceImpl1) Test() int {
	return 1
}

func (if2 *InterfaceImpl2) Test() int {
	return 2
}

// AOP
type (
	AopInterface interface {
		Test() string
	}

	AopImpl struct{}

	AopDep struct {
		A AopInterface `inject:""`
	}

	AopInterceptor1 struct {
		AopInterface
	}

	AopInterceptor2 struct {
		AopInterface
	}

	AopModule struct{}
)

func (m *AopModule) Configure(injector *Injector) {
	injector.Bind((*AopInterface)(nil)).To(AopImpl{})

	injector.BindInterceptor((*AopInterface)(nil), AopInterceptor1{})
	injector.BindInterceptor((*AopInterface)(nil), AopInterceptor2{})
}

func (a *AopImpl) Test() string {
	return "Test"
}

func (a *AopInterceptor1) Test() string {
	return a.AopInterface.Test() + " 1"
}

func (a *AopInterceptor2) Test() string {
	return a.AopInterface.Test() + " 2"
}
//
//var _ = Describe("Dingo Test", func() {
//	Context("Simple resolve", func() {
//		It("Should resolve dependencies on request", func() {
//			injector := NewInjector(new(PreTestModule), new(TestModule))
//
//			var iface Interface
//			iface = injector.GetInstance(new(Interface)).(Interface)
//
//			Expect(iface.Test()).To(Equal(1))
//
//			var dt DepTest = *injector.GetInstance(new(DepTest)).(*DepTest)
//
//			Expect(dt.Iface.Test()).To(Equal(1))
//			Expect(dt.Iface2.Test()).To(Equal(2))
//
//			var dt2 DepTest
//			injector.requestInjection(&dt2)
//
//			Expect(dt2.Iface.Test()).To(Equal(1))
//			Expect(dt2.Iface2.Test()).To(Equal(2))
//
//			Expect(dt.IfaceProvided.Test()).To(Equal(1))
//			Expect(dt.IfaceInstance.Test()).To(Equal(2))
//
//			Expect(dt.IfaceProvider().Test()).To(Equal(1))
//			Expect(dt.IfaceProvided.(*InterfaceImpl1).foo).To(Equal("Hello World"))
//		})
//
//		It("Should resolve scopes", func() {
//			injector := NewInjector(new(TestModule))
//
//			Expect(injector.GetInstance(TestSingleton{})).To(Equal(injector.GetInstance(TestSingleton{})))
//		})
//	})
//
//	Context("AOP", func() {
//		It("Should intercept interfaces", func() {
//			injector := NewInjector(new(AopModule))
//
//			var dep AopDep
//			injector.requestInjection(&dep)
//
//			Expect(dep.A.Test()).To(Equal("Test 1 2"))
//		})
//	})
//})
//
//func TestServiceContainer(t *testing.T) {
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "Dingo Suite")
//}
