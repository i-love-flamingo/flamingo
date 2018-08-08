package templatefunctions

import (
	"math"
	"reflect"
)

type (
	// JsMath is exported as a template function
	JsMath struct{}

	// Math is our Javascript's Math equivalent
	Math struct{}
)

// Func as implementation of debug method
func (ml JsMath) Func() interface{} {
	return func() Math {
		return Math{}
	}
}

// Ceil rounds a value up to the next biggest integer
func (m Math) Ceil(x interface{}) int {
	if reflect.TypeOf(x).Kind() == reflect.Int {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Int64 {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Float64 {
		x = float64(reflect.ValueOf(x).Float())
	}
	return int(math.Ceil(x.(float64)))
}

// Min gets the minimum value
func (m Math) Min(x ...interface{}) (res float64) {
	res = float64(math.MaxFloat64)
	for _, v := range x {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			v = float64(reflect.ValueOf(v).Int())
		} else if reflect.TypeOf(v).Kind() == reflect.Int64 {
			v = float64(reflect.ValueOf(v).Int())
		} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
			v = float64(reflect.ValueOf(v).Float())
		}
		if v.(float64) < res {
			res = v.(float64)
		}
	}
	return
}

// Max gets the maximum value
func (m Math) Max(x ...interface{}) (res float64) {
	res = float64(math.SmallestNonzeroFloat64)
	for _, v := range x {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			v = float64(reflect.ValueOf(v).Int())
		} else if reflect.TypeOf(v).Kind() == reflect.Int64 {
			v = float64(reflect.ValueOf(v).Int())
		} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
			v = float64(reflect.ValueOf(v).Float())
		}
		if v.(float64) > res {
			res = v.(float64)
		}
	}
	return
}

// Trunc drops the decimals
func (m Math) Trunc(x interface{}) int {
	if reflect.TypeOf(x).Kind() == reflect.Int {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Int64 {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Float64 {
		x = float64(reflect.ValueOf(x).Float())
	}
	return int(math.Trunc(x.(float64)))
}

func round(n float64) float64 {
	if n >= 0.5 {
		return math.Trunc(n + 0.5)
	}
	if n <= -0.5 {
		return math.Trunc(n - 0.5)
	}
	if math.IsNaN(n) {
		return math.NaN()
	}
	return 0
}

// Round rounds a value to the nearest integer
func (m Math) Round(x interface{}) int {
	if reflect.TypeOf(x).Kind() == reflect.Int {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Int64 {
		x = float64(reflect.ValueOf(x).Int())
	} else if reflect.TypeOf(x).Kind() == reflect.Float64 {
		x = float64(reflect.ValueOf(x).Float())
	} else {
		return 0
	}
	return int(round(x.(float64)))
}
