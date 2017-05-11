package template_functions

import (
	"math"
	"reflect"
)

type (
	// MathLib is exported as a template function
	MathLib struct{}

	// Math is our Javascript's Math equivalent
	Math struct{}
)

// Name alias for use in template
func (ml MathLib) Name() string {
	return "Math"
}

// Func as implementation of debug method
func (ml MathLib) Func() interface{} {
	return func() Math {
		return Math{}
	}
}

// Ceil rounds a value up to the next biggest integer
func (m Math) Ceil(x interface{}) int64 {
	if reflect.TypeOf(x).Kind() == reflect.Int64 {
		x = float64(reflect.ValueOf(x).Int())
	}
	return int64(math.Ceil(x.(float64)))
}

// Min gets the minimum value
func (m Math) Min(x ...interface{}) (res int64) {
	for _, v := range x {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			v = reflect.ValueOf(v).Int()
		}
		if v.(int64) < res {
			res = v.(int64)
		}
	}
	return
}

// Max gets the maximum value
func (m Math) Max(x ...interface{}) (res int64) {
	for _, v := range x {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			v = reflect.ValueOf(v).Int()
		}
		if v.(int64) > res {
			res = v.(int64)
		}
	}
	return
}
