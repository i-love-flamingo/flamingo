package pugast

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JS Expression transpiling", func() {
	var p = NewPugAst("/")

	Describe("JsExpr modes", func() {
		Context("With raw, wrap", func() {
			It("Should treat code as escaped blocks of JavaScript", func() {
				Expect(p.JsExpr(`var a = 1`, true, true)).To(Equal(`{{- $a := 1 -}}`))
			})

			It("Should fail with panic on invalid code", func() {
				Expect(func() { p.JsExpr(`[1,2,`, false, true) }).To(Panic())
			})
		})

		Context("With raw, not wrap", func() {
			It("Should treat code as escaped blocks of JavaScript", func() {
				Expect(p.JsExpr(`var a = 1`, false, true)).To(Equal(`$a := 1`))
			})

			It("Should fail with panic on invalid code", func() {
				Expect(func() { p.JsExpr(`[1,2,`, false, true) }).To(Panic())
			})
		})

		Context("With not raw, wrap", func() {
			It("Should treat code as escaped blocks of JavaScript", func() {
				Expect(p.JsExpr(`{"key": "value"}`, true, false)).To(Equal(`{{(__op__map "key" "value")}}`))
			})

			It("Should fail with panic on invalid code", func() {
				Expect(func() { p.JsExpr(`[1,2,`, false, false) }).To(Panic())
			})
		})

		Context("With not raw, not wrap", func() {
			It("Should treat code as escaped blocks of JavaScript", func() {
				Expect(p.JsExpr(`{"key": "value"}`, false, false)).To(Equal(`(__op__map "key" "value")`))
			})

			It("Should fail with panic on variable statements", func() {
				Expect(func() { p.JsExpr(`var a = 1`, false, false) }).To(Panic())
			})
		})
	})

	Describe("Function renderExpression", func() {
		It("Should not fail for empty input", func() {
			Expect(p.JsExpr(``, true, true)).To(Equal(""))
		})

		Context("Transpile Identifier", func() {
			It("Should transpile it correctly if it is known", func() {
				Expect(p.JsExpr(`testknown`, true, true)).To(Equal(`{{$testknown}}`))
			})
			It("Should transpile it correctly if it is not known", func() {
				Expect(p.JsExpr(`testknown`, true, true)).To(Equal(`{{$testknown}}`))
			})
			It("Should make it raw if rawmode is on", func() {
				p.rawmode = true
				Expect(p.JsExpr(`testknown`, true, true)).To(Equal(`{{$testknown | raw}}`))
				p.rawmode = false
			})
		})

		Context("Transpile String Literal", func() {
			It("Should interpolate if necessary", func() {
				Expect(p.JsExpr(`"foo${a} \$${1+2}"`, true, false)).To(Equal(`{{(s "foo" $a " $" (__op__add 1 2) )}}`))
			})

			It("Should strip unnecessary template stuff for raw strings without interpolation", func() {
				Expect(p.JsExpr(`"test"`, true, false)).To(Equal(`test`))
				Expect(p.JsExpr(`"test"`, false, false)).To(Equal(`"test"`))
				Expect(p.JsExpr(`"<test>"`, true, false)).To(Equal(`&lt;test&gt;`))
				Expect(p.JsExpr(`"<test>"`, false, false)).To(Equal(`"<test>"`))
			})
		})

		Context("Transpile Array Literal", func() {
			It("Should map arrays to __op__array", func() {
				Expect(p.JsExpr(`[1, 2, 3]`, true, false)).To(Equal(`{{(__op__array 1 2 3)}}`))
				Expect(p.JsExpr(`[1, 2, 3]`, false, false)).To(Equal(`(__op__array 1 2 3)`))
			})
		})

		Context("Transpile Boolean expression", func() {
			It("Should be true and false", func() {
				Expect(p.JsExpr(`true`, false, false)).To(Equal(`true`))
				Expect(p.JsExpr(`false`, false, false)).To(Equal(`false`))
			})
		})

		Context("Transpile Map Literal", func() {
			It("Should map objects to __op__map", func() {
				Expect(p.JsExpr(`{"key": 1, "key2": {"key1": [1+2, 3, 4]}}`, true, false)).To(Equal(`{{(__op__map "key" 1 "key2" (__op__map "key1" (__op__array (__op__add 1 2) 3 4)))}}`))
				Expect(p.JsExpr(`{"key": 1, "key2": {"key1": [1+2, 3, 4]}}`, false, false)).To(Equal(`(__op__map "key" 1 "key2" (__op__map "key1" (__op__array (__op__add 1 2) 3 4)))`))
			})
		})

		Context("Transpile Null Literal", func() {
			It("Should be null if wrapped", func() {
				Expect(p.JsExpr(`null`, true, false)).To(Equal(`{{null}}`))
			})
			It("Should be empty if not wrapped", func() {
				Expect(p.JsExpr(`null`, false, false)).To(Equal(``))
			})
		})

		Context("Transpile Dot Expression", func() {
			It("Should use dot-notation", func() {
				Expect(p.JsExpr(`a.b`, false, false)).To(Equal(`$a.b`))
			})
			It("Should be raw and escaped if rawmode and wrap is set", func() {
				p.rawmode = true
				Expect(p.JsExpr(`a.b`, true, false)).To(Equal(`{{$a.b | raw}}`))
				p.rawmode = false
			})
		})

		Context("Transpile Conditional Expression", func() {
			It("Should become and if-else statement", func() {
				Expect(p.JsExpr(`a ? b : c`, false, false)).To(Equal(`{{if $a}}{{$b}}{{else}}{{$c}}{{end}}`))
			})
			It("Should eliminate null else branches", func() {
				Expect(p.JsExpr(`a ? b : null`, false, false)).To(Equal(`{{if $a}}{{$b}}{{end}}`))
			})
		})

		Context("Transpile Binary Expressions", func() {
			It("Should handle &-operator", func() {
				Expect(p.JsExpr(`a & b`, true, false)).To(Equal(`{{(__op__b_and $a $b)}}`))
			})
		})

		Context("Transpile Call Expressions", func() {
			It("Should transform js-call-syntax to go template call syntax", func() {
				p.FuncMap = FuncMap{"foo": func(int, int) {}}
				Expect(p.JsExpr(`foo(1+2)`, true, false)).To(Equal(`{{(foo (__op__add 1 2))}}`))
			})
		})

		Context("Transpile Assign Expressions", func() {
			It("Should assign expressions to variables", func() {
				Expect(p.JsExpr(`a = 1`, true, false)).To(Equal(`{{- $a := 1 -}}`))
			})
		})

		Context("Transpile Sequence Expression", func() {
			It("Should use __op__array and not wrap", func() {
				Expect(p.JsExpr(`1,2,3`, true, false)).To(Equal(`(__op__array 1 2 3)`))
			})
		})

		Context("Transpile Bracket Expression", func() {
			It("Should use the index function to access the specified element", func() {
				Expect(p.JsExpr(`a[0][b[1]]`, true, false)).To(Equal(`{{(index (index $a 0) (index $b 1))}}`))
			})
		})
	})

	Describe("Known Bugs", func() {
		It("Should render brand.heroImage.url to $brand.heroImage.url", func() {
			Expect(p.JsExpr("`background-image:url(${brand.heroImage.url})`", false, false)).To(Equal(`(s "background-image:url(" $brand.heroImage.url ")")`))
		})
	})
})

func TestJsExpr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JsExpr Suite")
}
