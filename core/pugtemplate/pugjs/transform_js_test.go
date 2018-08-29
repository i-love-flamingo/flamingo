package pugjs

import (
	"testing"

	"flamingo.me/flamingo/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

func TestJsExpr(t *testing.T) {
	var s = newRenderState("/", true, nil, flamingo.NullLogger{})

	t.Run("JsExpr modes", func(t *testing.T) {
		t.Run("With raw, wrap", func(t *testing.T) {
			assert.Equal(t, `{{ $a := 1 -}}`, s.JsExpr(`var a = 1`, true, true))
			assert.Panics(t, func() { s.JsExpr(`[1,2,`, false, true) })
		})

		t.Run("With raw, not wrap", func(t *testing.T) {
			assert.Equal(t, `$a := 1`, s.JsExpr(`var a = 1`, false, true))
			assert.Panics(t, func() { s.JsExpr(`[1,2,`, false, true) })
		})

		t.Run("With not raw, wrap", func(t *testing.T) {
			assert.Equal(t, `{{(__op__map "key" "value")}}`, s.JsExpr(`{"key": "value"}`, true, false))
			assert.Panics(t, func() { s.JsExpr(`[1,2,`, false, false) })
		})

		t.Run("With not raw, not wrap", func(t *testing.T) {
			assert.Equal(t, `(__op__map "key" "value")`, s.JsExpr(`{"key": "value"}`, false, false))
			assert.Panics(t, func() { s.JsExpr(`var a = 1`, false, false) })
		})
	})

	t.Run("Function renderExpression", func(t *testing.T) {
		assert.Equal(t, "", s.JsExpr(``, true, true))

		t.Run("Transpile Identifier", func(t *testing.T) {
			assert.Equal(t, `{{$testknown | __pug__html}}`, s.JsExpr(`testknown`, true, true))
			assert.Equal(t, `{{$testknown | __pug__html}}`, s.JsExpr(`testknown`, true, true))

			s.rawmode = true
			assert.Equal(t, `{{$testknown}}`, s.JsExpr(`testknown`, true, true))
			s.rawmode = false
		})

		t.Run("Transpile String Literal", func(t *testing.T) {
			assert.Equal(t, `{{(__str "foo" $a " $" (__op__add 1 2) )}}`, s.JsExpr(`"foo${a} \$${1+2}"`, true, false))

			assert.Equal(t, `test`, s.JsExpr(`"test"`, true, false))
			assert.Equal(t, `"test"`, s.JsExpr(`"test"`, false, false))
			assert.Equal(t, `&lt;test&gt;`, s.JsExpr(`"<test>"`, true, false))
			assert.Equal(t, `"<test>"`, s.JsExpr(`"<test>"`, false, false))
		})

		t.Run("Transpile Array Literal", func(t *testing.T) {
			assert.Equal(t, `{{(__op__array 1 2 3)}}`, s.JsExpr(`[1, 2, 3]`, true, false))
			assert.Equal(t, `(__op__array 1 2 3)`, s.JsExpr(`[1, 2, 3]`, false, false))
		})

		t.Run("Transpile Boolean expression", func(t *testing.T) {
			assert.Equal(t, `true`, s.JsExpr(`true`, false, false))
			assert.Equal(t, `false`, s.JsExpr(`false`, false, false))
		})

		t.Run("Transpile Map Literal", func(t *testing.T) {
			assert.Equal(t, `{{(__op__map "key" 1 "key2" (__op__map "key1" (__op__array (__op__add 1 2) 3 4)))}}`, s.JsExpr(`{"key": 1, "key2": {"key1": [1+2, 3, 4]}}`, true, false))
			assert.Equal(t, `(__op__map "key" 1 "key2" (__op__map "key1" (__op__array (__op__add 1 2) 3 4)))`, s.JsExpr(`{"key": 1, "key2": {"key1": [1+2, 3, 4]}}`, false, false))
		})

		t.Run("Transpile Null Literal", func(t *testing.T) {
			assert.Equal(t, `{{null}}`, s.JsExpr(`null`, true, false))
			assert.Equal(t, ``, s.JsExpr(`null`, false, false))
		})

		t.Run("Transpile Dot Expression", func(t *testing.T) {
			assert.Equal(t, `$a.b`, s.JsExpr(`a.b`, false, false))

			s.rawmode = true
			assert.Equal(t, `{{$a.b}}`, s.JsExpr(`a.b`, true, false))
			s.rawmode = false
		})

		t.Run("Transpile Conditional Expression", func(t *testing.T) {
			//assert.Equal(t, (s.JsExpr(`a ? b : c`, false, false)).To(Equal(`{{if $a}}{{$b}}{{else}}{{$c}}{{end}}`))
			//assert.Equal(t, (s.JsExpr(`a ? b : null`, false, false)).To(Equal(`{{if $a}}{{$b}}{{end}}`))
		})

		t.Run("Transpile Binary Expressions", func(t *testing.T) {
			assert.Equal(t, `{{(__op__b_and $a $b)}}`, s.JsExpr(`a & b`, true, false))
		})

		t.Run("Transpile Call Expressions", func(t *testing.T) {
			s.funcs = FuncMap{"foo": func(int, int) {}}
			assert.Equal(t, `{{(foo (__op__add 1 2)) | __pug__html}}`, s.JsExpr(`foo(1+2)`, true, false))
		})

		t.Run("Transpile Assign Expressions", func(t *testing.T) {
			assert.Equal(t, `{{ $a := 1 -}}`, s.JsExpr(`a = 1`, true, false))
		})

		t.Run("Transpile Sequence Expression", func(t *testing.T) {
			assert.Equal(t, `(__op__array 1 2 3)`, s.JsExpr(`1,2,3`, true, false))
		})

		t.Run("Transpile Bracket Expression", func(t *testing.T) {
			assert.Equal(t, `{{(__pug__index (__pug__index $a 0) (__pug__index $b 1)) | __pug__html}}`, s.JsExpr(`a[0][b[1]]`, true, false))
		})
	})

	t.Run("Known Bugs", func(t *testing.T) {
		assert.Equal(t, `(__str "background-image:url(" $brand.heroImage.url ")")`, s.JsExpr("`background-image:url(${brand.heroImage.url})`", false, false))
	})
}
