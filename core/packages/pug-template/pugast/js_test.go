package pugast

import "testing"

type Case struct {
	Test, Expected string
	Raw            bool
}

var cases = []Case{
	{"var a = 1", "$a := 1", true},
	{"1 + 2", "(__op__add 1 2)", false},
	{"[1, 2, 3, 4, 5]", "(__op__array 1 2 3 4 5)", false},
	{`{"key": "value", "key2": [1, 2, 3]}`, `(__op__map  "key" "value" "key2" (__op__array 1 2 3))`, false},
	{`Func(1, 2, 3)`, `(Func 1 2 3)`, false},
}

func TestJsExpr(t *testing.T) {
	for _, c := range cases {
		res := JsExpr(c.Test, false, c.Raw)
		if res != c.Expected {
			t.Errorf("%s: %s != %s", c.Test, res, c.Expected)
		}
	}
}
