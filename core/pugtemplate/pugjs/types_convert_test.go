package pugjs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testConvertStruct1 struct {
		Str string
		Num int
	}

	testConvertInterfaceEmpty interface{}

	testConvertInterface1 interface {
		Method() string
	}

	testConvertPrimitive string
)

func (tcs *testConvertStruct1) PtrMethod() string {
	return "PtrMethod String"
}

func (tcs testConvertStruct1) Method() string {
	return "Method String"
}

func (tcp testConvertPrimitive) Method() string {
	return "primitive implementation"
}

func TestConvert(t *testing.T) {
	t.Run("Complex structs", func(t *testing.T) {
		tcs1s := testConvertStruct1{Str: "TestStr", Num: 1337}
		tcs1 := convert(tcs1s).(*Map).Items

		assert.Equal(t, tcs1[String("str")], String("TestStr"))
		assert.Equal(t, tcs1[String("num")], Number(1337))
		assert.Equal(t, tcs1[String("method")].(*Func).fnc.Type(), reflect.TypeOf(tcs1s.Method))
		assert.NotContains(t, tcs1, String("ptrMethod"))
	})

	t.Run("Pointer", func(t *testing.T) {
		tcs2s := &testConvertStruct1{Str: "TestStr", Num: 1337}
		tcs2 := convert(tcs2s).(*Map).Items

		assert.Equal(t, tcs2[String("str")], String("TestStr"))
		assert.Equal(t, tcs2[String("num")], Number(1337))
		assert.Equal(t, tcs2[String("method")].(*Func).fnc.Type(), reflect.TypeOf(tcs2s.Method))
		assert.Equal(t, tcs2[String("ptrMethod")].(*Func).fnc.Type(), reflect.TypeOf(tcs2s.PtrMethod))

		assert.Equal(t, Nil{}, convert((*testConvertStruct1)(nil)))
	})

	t.Run("Pointer/Interfaces", func(t *testing.T) {
		// explicit empty interface
		tcs1s := testConvertInterfaceEmpty(testConvertStruct1{Str: "TestStr", Num: 1337})
		tcs1 := convert(tcs1s).(*Map).Items

		assert.Equal(t, tcs1[String("str")], String("TestStr"))
		assert.Equal(t, tcs1[String("num")], Number(1337))
		assert.Equal(t, tcs1[String("method")].(*Func).fnc.Type(), reflect.TypeOf(tcs1s.(testConvertStruct1).Method))
		assert.NotContains(t, tcs1, String("ptrMethod"))

		// interface on struct
		tcs2s := testConvertInterface1(testConvertStruct1{Str: "TestStr", Num: 1337})
		tcs2 := convert(tcs2s).(*Map).Items

		assert.Equal(t, tcs2[String("str")], String("TestStr"))
		assert.Equal(t, tcs2[String("num")], Number(1337))
		assert.Equal(t, tcs2[String("method")].(*Func).fnc.Type(), reflect.TypeOf(tcs2s.(testConvertStruct1).Method))
		assert.NotContains(t, tcs2, String("ptrMethod"))
	})

	t.Run("Primitive types", func(t *testing.T) {
		testmaps := []interface{}{
			map[string]interface{}{"foo": "bar", "xxx": 1},
			map[string]interface{}{},
		}

		testfuncs := []interface{}{
			func() {},
			func(string) string { return "" },
			func(int, string, int) {},
		}

		teststructs := []interface{}{
			struct{ Foo, Bar string }{Foo: "foofoo", Bar: "barbar"},
		}

		expected := []struct{ in, out interface{} }{
			// Special Cases
			{nil, Nil{}},                                // nil
			{String("a"), String("a")},                  // object -> object
			{reflect.ValueOf(String("a")), String("a")}, // reflect.Value(object) -> object
			{reflect.Value{}, Nil{}},                    // invalid reflect
			{errors.New("test"), String("Error: test")}, // errors

			// Strings
			{"foo", String("foo")},
			{"", String("")},
			{"a b c -da0sdoa0wdw", String("a b c -da0sdoa0wdw")},

			// Numbers
			{0, Number(0)},
			{1, Number(1)},
			{1.2, Number(1.2)},
			{-1111, Number(-1111)},
			{uint8(1), Number(1)},
			{complex(1, 1), Nil{}},

			// Channel
			{make(chan bool), Nil{}},

			// Bool
			{true, Bool(true)},
			{false, Bool(false)},

			// Arrays
			{[]string{"foo", "bar"}, &Array{items: []Object{String("foo"), String("bar")}, o: []string{"foo", "bar"}}},
			{[]interface{}{1, "bar", nil}, &Array{items: []Object{Number(1), String("bar"), Nil{}}, o: []interface{}{1, "bar", nil}}},

			// Maps
			{testmaps[0], &Map{Items: map[Object]Object{String("foo"): String("bar"), String("xxx"): Number(1)}, o: testmaps[0]}},
			{testmaps[1], &Map{Items: map[Object]Object{}, o: testmaps[1]}},

			// Functions
			{testfuncs[0], &Func{fnc: reflect.ValueOf(testfuncs[0])}},
			{testfuncs[1], &Func{fnc: reflect.ValueOf(testfuncs[1])}},
			{testfuncs[2], &Func{fnc: reflect.ValueOf(testfuncs[2])}},

			// Structs
			{teststructs[0], &Map{Items: map[Object]Object{String("foo"): String("foofoo"), String("bar"): String("barbar")}, o: teststructs[0]}},
		}

		for _, e := range expected {
			assert.Equal(t, e.out, Convert(e.in))
		}
	})
}
