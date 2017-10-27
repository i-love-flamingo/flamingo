package pugjs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
)

// FuncMap is the default runtime funcmap for pugast templates
var funcmap = FuncMap{
	"__": func(s ...string) string { return strings.Join(s, "::") },

	"__op__add":   runtimeAdd,
	"__op__inc":   runtimeInc,
	"__op__sub":   runtimeSub,
	"__op__mul":   runtimeMul,
	"__op__quo":   runtimeQuo,
	"__op__slash": runtimeQuo,
	"__op__mod":   runtimeRem,
	"__op__eql":   runtimeEql,

	"__op__lt": runtimeLss,
	"__op__gt": func(x, y interface{}) bool {
		return !runtimeLss(x, y) && !runtimeEql(x, y)
	},
	"__op__gte": func(x, y interface{}) bool {
		return !runtimeLss(x, y)
	},
	"__op__lte": func(x, y interface{}) bool {
		return runtimeLss(x, y) || runtimeEql(x, y)
	},
	"__op__neq": func(x, y interface{}) bool {
		return !runtimeEql(x, y)
	},

	"neq": func(x, y interface{}) bool { return !runtimeEql(x, y) },
	"__tryindex": func(obj, key interface{}) interface{} {
		arr, ok := obj.(*Array)
		idx, ok2 := key.(int)
		if ok && ok2 {
			if len(arr.items) <= idx {
				return Nil{}
			}
			return arr.items[idx]
		}

		if obj, ok := obj.(Object); ok {
			return obj.Field(convert(key).String())
		}

		vo, _ := indirect(reflect.ValueOf(obj))
		k := int(reflect.ValueOf(key).Int())
		if !vo.IsValid() {
			return nil
		}
		if vo.Len() > k {
			return vo.Index(k).Interface()
		}
		return nil
	},
	"json":      runtimeJSON,
	"unescaped": runtimeUnescaped,
	"null":      func() interface{} { return Nil{} },
	"__Range": func(args ...int) Object {
		var res []int
		var m, o int
		if len(args) == 1 {
			m = args[0]
			o = 0
		} else {
			m = args[1]
			o = args[0]
		}

		for i := o; i < m; i++ {
			res = append(res, i)
		}
		return convert(res)
	},
	"__range_helper__": func(o Object) interface{} {
		switch o := o.(type) {
		case *Map:
			return o.Items
		case *Array:
			return o.items
		case String:
			return string(o)
		}
		return nil
	},
	"__range_helper_keys__": func(o Object) []interface{} {
		var res []interface{}
		switch o := o.(type) {
		case *Map:
			for k := range o.Items {
				res = append(res, k)
			}
		case *Array:
			for i := range o.items {
				res = append(res, i)
			}
		}
		return res
	},
	"raw": func(s ...interface{}) template.HTML { return template.HTML(fmt.Sprint(s...)) },
	"__str": func(l ...interface{}) string {
		var res string
		for _, s := range l {
			res += convert(s).String()
		}
		if len(res) > 1 {
			return " " + strings.TrimSpace(res)
		}
		return ""
	},
	"__op__array": func(a ...interface{}) Object { return convert(a) },
	"__op__map": func(a ...interface{}) Object {
		m := make(map[interface{}]interface{}, len(a)/2)
		for i := 0; i < len(a); i += 2 {
			m[a[i]] = a[i+1]
		}
		return convert(m)
	},
	"__op__map_params": func(a ...interface{}) Object {
		m := make(map[interface{}]interface{}, len(a)/2)
		for i := 0; i < len(a); i += 2 {
			if _, ok := m[a[i]]; ok {
				if x, ok := m[a[i]].([]interface{}); ok {
					m[a[i]] = append(x, a[i+1])
				} else {
					m[a[i]] = []interface{}{m[a[i]], a[i+1]}
				}
			} else {
				m[a[i]] = a[i+1]
			}
		}
		return convert(m)
	},
	"__attr": func(k string, v interface{}, e bool) []Attribute {
		if v, ok := v.(Bool); ok {
			b := v.True()
			return []Attribute{{Name: k, BoolVal: &b}}
		}
		if v, ok := v.(bool); ok {
			return []Attribute{{Name: k, BoolVal: &v}}
		}
		if v, ok := v.(Object); ok {
			return []Attribute{{Name: k, Val: JavaScriptExpression(v.String()), MustEscape: e}}
		}
		if v, ok := v.(string); ok {
			return []Attribute{{Name: k, Val: JavaScriptExpression(string(v)), MustEscape: e}}
		}
		return []Attribute{{Name: k, Val: JavaScriptExpression(fmt.Sprintf("%t", v)), MustEscape: e}}
	},
	"__attrs": func(attrs ...*Array) (res string) {
		type tmpattr struct {
			mustEscape bool
			val        string
			bool       *bool
		}
		a := make(map[string][]tmpattr)
		var order []string
		for _, list := range attrs {
		attrloop:
			for _, attr := range list.items {
				if attr == nil {
					continue
				}
				name := string(attr.(*Map).Items[String("Name")].(String))
				mustEscape := bool(attr.(*Map).Items[String("MustEscape")].(Bool))
				var val string
				att := tmpattr{mustEscape: mustEscape}
				if _, ok := attr.(*Map).Items[String("BoolVal")].(Bool); attr.(*Map).Items[String("BoolVal")] != nil && ok {
					b := attr.(*Map).Items[String("BoolVal")].(Bool).True()
					att.bool = &b
					if mustEscape {
						val = name
					} else {
						val = `"` + name + `"`
					}
				} else {
					val = string(attr.(*Map).Items[String("Val")].(String))
				}
				att.val = val
				if _, ok := a[name]; ok {
					if name == "class" {
						for _, s := range a[name] {
							if s == att {
								// we already now this attribute value for class, continue
								continue attrloop
							}
						}
						a[name] = append(a[name], att)
					} else {
						a[name] = []tmpattr{att}
					}
				} else {
					a[name] = []tmpattr{att}
					order = append(order, name)
				}
			}
		}
	renderloop:
		for _, attr := range order {
			var tmp string
			for _, val := range a[attr] {
				if val.bool != nil && !*val.bool {
					continue renderloop
				}
				if len(tmp) > 0 {
					tmp += ` `
				}
				if val.mustEscape {
					tmp += template.HTMLEscapeString(val.val)
				} else if val.val[0] == '"' {
					tmp += val.val[1 : len(val.val)-1]
				}
			}
			res += ` ` + attr + `="` + strings.TrimSpace(tmp) + `"`
		}
		return
	},
	"__and_attrs": func(x *Map) (res []Attribute) {
		for k, v := range x.Items {
			if b, ok := v.(Bool); ok {
				boolval := b.True()
				res = append(res, Attribute{Name: k.String(), Val: JavaScriptExpression(v.String()), MustEscape: true, BoolVal: &boolval})
			} else {
				res = append(res, Attribute{Name: k.String(), Val: JavaScriptExpression(v.String()), MustEscape: true})
			}
		}
		return
	},
	"__if": func(test, left, right interface{}) interface{} {
		if t, ok := IsTrue(test); ok && t {
			return left
		}
		return right
	},

	"parseInt": func(num Object, base Number) Number {
		f, _ := strconv.ParseFloat(num.String(), 64)
		n, err := strconv.ParseInt(strconv.Itoa(int(f)), int(base), 64)
		if err != nil {
			panic(err)
		}
		return Number(n)
	},

	"__freeze": func(name string) Nil {
		return Nil{}
	},
}

func runtimeAdd(l, r interface{}) Object {
	x := convert(l)
	y := convert(r)

	switch x := x.(type) {
	case String:
		return String(x.String() + y.String())

	case Number:
		switch y := y.(type) {
		case Number:
			return Number(x + y)

		case String:
			f, _ := strconv.ParseFloat(y.String(), 64)
			return Number(float64(x) + f)
		}
	}
	return Nil{}
}

func runtimeInc(x interface{}) int {
	vx := reflect.ValueOf(x)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			return int(vx.Int() + 1)
		}
	case reflect.Float32, reflect.Float64:
		{
			return int(vx.Float() + 1)
		}
	}

	return 0
}

func runtimeSub(i ...interface{}) interface{} {
	y := i[0]
	var x interface{}
	if len(i) > 1 {
		x = i[1]
	} else {
		x = 0
	}
	x, y = y, x
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return int(vx.Int() - vy.Int())
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

func runtimeMul(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return int(vx.Int() * vy.Int())
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

func runtimeQuo(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return float64(vx.Int()) / float64(vy.Int())
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

func runtimeRem(x, y interface{}) interface{} {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	switch vx.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return int(vx.Int() % vy.Int())
			}
		}
	case reflect.Float64, reflect.Float32:
		{
			switch vy.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
				return int64(vx.Float()) % vy.Int()
			case reflect.Float32, reflect.Float64:
				return int64(vx.Float()) % int64(vy.Float())
			}
		}
	}

	return "<nil>"
}

func runtimeEql(x, y interface{}) bool {
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

func runtimeLss(x, y interface{}) bool {
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

func runtimeJSON(x interface{}) (res template.JS, err error) {
	bres, err := json.Marshal(x)
	res = template.JS(string(bres))
	return
}

func runtimeUnescaped(x string) interface{} {
	return template.HTML(x)
}
