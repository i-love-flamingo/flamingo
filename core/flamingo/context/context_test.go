package context

import (
	"testing"

	"flamingo/core/flamingo/service_container"

	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type (
	TestDependency struct {
		marker string
	}

	TestDependecyChecker struct {
		*TestDependency `inject:""`
	}
)

func Register(marker string) service_container.RegisterFunc {
	return func(sc *service_container.ServiceContainer) {
		sc.Register(&TestDependency{marker: marker})
	}
}

var _ = g.Describe("Context", func() {

	var root = New(
		"root",
		[]service_container.RegisterFunc{
			Register("root"),
		},
		New(
			"main",
			nil,
			New("main1", nil),
			New("main2", nil),
		),
		New(
			"not-main",
			[]service_container.RegisterFunc{
				Register("not-main"),
			},
			New("notmain1", nil),
			New(
				"notmain2",
				nil,
				New("notmain2-1", nil),
				New("notmain2-2", []service_container.RegisterFunc{Register("not-main-deep")}),
			),
		),
	)

	g.Context("Merge Tree Behaviour", func() {
		root.Contexts = map[string]string{
			"b1":   "main/main1",
			"b2":   "main/main2",
			"nb1":  "not-main/notmain1",
			"nb21": "not-main/notmain2/notmain2-1",
			"nb22": "not-main/notmain2/notmain2-2",
		}

		flat := root.GetFlatContexts()

		g.It("Should render correct Contexts", func() {
			Expect(flat).To(HaveKey("main/main1"))
			Expect(flat).To(HaveKey("main/main2"))
			Expect(flat).To(HaveKey("not-main/notmain1"))
			Expect(flat).To(HaveKey("not-main/notmain2/notmain2-1"))
			Expect(flat).To(HaveKey("not-main/notmain2/notmain2-2"))
		})

		g.It("Should have correct DI", func() {
			tester := new(TestDependecyChecker)

			flat["main/main1"].ServiceContainer.Register(tester)
			flat["main/main1"].ServiceContainer.Resolve()
			Expect(tester.marker).To(Equal("root"))

			tester = new(TestDependecyChecker)
			flat["main/main2"].ServiceContainer.Register(tester)
			flat["main/main2"].ServiceContainer.Resolve()
			Expect(tester.marker).To(Equal("root"))

			tester = new(TestDependecyChecker)
			flat["not-main/notmain1"].ServiceContainer.Register(tester)
			flat["not-main/notmain1"].ServiceContainer.Resolve()
			Expect(tester.marker).To(Equal("not-main"))

			tester = new(TestDependecyChecker)
			flat["not-main/notmain2/notmain2-1"].ServiceContainer.Register(tester)
			flat["not-main/notmain2/notmain2-1"].ServiceContainer.Resolve()
			Expect(tester.marker).To(Equal("not-main"))

			tester = new(TestDependecyChecker)
			flat["not-main/notmain2/notmain2-2"].ServiceContainer.Register(tester)
			flat["not-main/notmain2/notmain2-2"].ServiceContainer.Resolve()
			Expect(tester.marker).To(Equal("not-main-deep"))
		})
	})
})

func TestContext(t *testing.T) {
	RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Context Suite")
}
