package node

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

var FuncMap = template.FuncMap{
	"asset": func(a string) template.URL { return template.URL("http://CDN/" + a) },
	"__":    func(s ...string) string { return strings.Join(s, "::") },
	"get": func(what string) interface{} {
		if what == "user.name" {
			return "testuser"
		}
		return []map[string]string{{"url": "url1", "name": "item1"}, {"url": "url2", "name": "name2"}}
	},

	"__op__add":   runtime_add,
	"__op__sub":   runtime_sub,
	"__op__mul":   runtime_mul,
	"__op__quo":   runtime_quo,
	"__op__rem":   runtime_rem,
	"__op__mod":   runtime_rem,
	"__op__minus": runtime_minus,
	"__op__plus":  runtime_plus,
	"__op__eql":   runtime_eql,
	"__op__gtr":   runtime_gtr,
	"__op__lss":   runtime_lss,
	"neq":         func(x, y interface{}) bool { return !runtime_eql(x, y) },

	"json":      runtime_json,
	"unescaped": runtime_unescaped,

	"null": func() interface{} { return nil },

	"raw":     func(s string) template.HTML { return template.HTML(s) },
	"tagopen": func(t, p string) template.HTML { return template.HTML(`<` + p + t) },
	"s": func(l ...interface{}) (res string) {
		for _, s := range l {
			vs := reflect.ValueOf(s)
			switch vs.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				{
					res += fmt.Sprintf("%d", vs.Int())
				}
			case reflect.Float32, reflect.Float64:
				{
					res += fmt.Sprintf("%d", vs.Float())
				}
			case reflect.String:
				{
					res += vs.String()
				}
			}
		}
		return
	},

	"__op__array": func(a ...interface{}) []interface{} { return a },
	"__op__map": func(a ...interface{}) map[interface{}]interface{} {
		m := make(map[interface{}]interface{})
		for i := 0; i < len(a); i += 2 {
			m[a[i]] = a[i+1]
		}
		return m
	},
	"attr": func(attr, prefix interface{}) string {
		if attr == nil {
			return ""
		}
		t := strings.Split(attr.(string), " ")
		for k, v := range t {
			t[k] = prefix.(string) + "-" + v
		}
		return strings.Join(t, " ")
	},
	"extend": func(deep bool, m map[interface{}]interface{}, on ...map[interface{}]interface{}) map[interface{}]interface{} {
		for _, o := range on {
			for k, v := range o {
				m[k] = v
			}
		}
		return m
	},
}

func runtime_add(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() + vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) + vy.Float()
			case reflect.String:
				return fmt.Sprintf("%d%s", vx.Int(), vy.String())
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() + float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() + vy.Float()
			case reflect.String:
				return fmt.Sprintf("%f%s", vx.Float(), vy.String())
			}
		}
	case reflect.String:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return fmt.Sprintf("%s%d", vx.String(), vy.Int())
			case reflect.Float32, reflect.Float64:
				return fmt.Sprintf("%s%f", vx.String(), vy.Float())
			case reflect.String:
				return fmt.Sprintf("%s%s", vx.String(), vy.String())
			}
		}
	}

	return "<nil>"
}

func runtime_sub(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() - vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) - vy.Float()
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() - float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() - vy.Float()
			}
		}
	}

	return "<nil>"
}

func runtime_mul(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() * vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) * vy.Float()
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() * float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() * vy.Float()
			}
		}
	}

	return "<nil>"
}

func runtime_quo(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() / vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) / vy.Float()
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() / float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() / vy.Float()
			}
		}
	}

	return "<nil>"
}

func runtime_rem(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() % vy.Int()
			}
		}
	}

	return "<nil>"
}

func runtime_minus(x interface{}) interface{} {
	vx := reflect.ValueOf(x)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return -vx.Int()
	case reflect.Float32, reflect.Float64:
		return -vx.Float()
	}

	return "<nil>"
}

func runtime_plus(x interface{}) interface{} {
	vx := reflect.ValueOf(x)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		return +vx.Int()
	case reflect.Float32, reflect.Float64:
		return +vx.Float()
	}

	return "<nil>"
}

func runtime_eql(x, y interface{}) bool {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() == vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) == vy.Float()
			case reflect.String:
				return fmt.Sprintf("%d", vx.Int()) == vy.String()
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() == float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() == vy.Float()
			case reflect.String:
				return fmt.Sprintf("%f", vx.Float()) == vy.String()
			}
		}
	case reflect.String:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.String() == fmt.Sprintf("%d", vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.String() == fmt.Sprintf("%f", vy.Float())
			case reflect.String:
				return vx.String() == fmt.Sprintf("%s", vy.String())
			}
		}
	case reflect.Bool:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Bool() && vy.Int() != 0
			case reflect.Bool:
				return vx.Bool() == vy.Bool()
			}
		}
	}

	return false
}

func runtime_lss(x, y interface{}) bool {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Int() < vy.Int()
			case reflect.Float32, reflect.Float64:
				return float64(vx.Int()) < vy.Float()
			case reflect.String:
				return fmt.Sprintf("%d", vx.Int()) < vy.String()
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.Float() < float64(vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.Float() < vy.Float()
			case reflect.String:
				return fmt.Sprintf("%f", vx.Float()) < vy.String()
			}
		}
	case reflect.String:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return vx.String() < fmt.Sprintf("%d", vy.Int())
			case reflect.Float32, reflect.Float64:
				return vx.String() < fmt.Sprintf("%f", vy.Float())
			case reflect.String:
				return vx.String() < vy.String()
			}
		}
	}

	return false
}

func runtime_gtr(x, y interface{}) bool {
	return !runtime_lss(x, y) && !runtime_eql(x, y)
}

func runtime_json(x interface{}) (res string, err error) {
	bres, err := json.Marshal(x)
	res = string(bres)
	return
}

func runtime_unescaped(x string) interface{} {
	return template.HTML(x)
}
