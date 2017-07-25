package pugast

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type ObjectType string

const (
	STRING = "string"
	NUMBER = "number"
	ARRAY  = "array"
	MAP    = "map"
	FUNC   = "func"
	BOOL   = "bool"
	NIL    = "nil"
)

type Object interface {
	Type() ObjectType
	Field(name string) Object
	String() string
}

func convert(in interface{}) Object {
	//log.Printf("Converting %#v", in)
	if in == nil {
		return Nil{}
	}

	if in, ok := in.(Object); ok {
		return in
	}

	val, ok := in.(reflect.Value)

	if !ok {
		val = reflect.ValueOf(in)
	}

	if !val.IsValid() {
		return Nil{}
	}

	if err, ok := in.(error); ok && err != nil {
		return String(fmt.Sprintf("%+v", err))
	}

	switch val.Kind() {
	case reflect.Slice:
		newval := &Array{
			items: make([]Object, val.Len()),
		}
		for i := 0; i < val.Len(); i++ {
			newval.items[i] = convert(val.Index(i))
		}
		return newval

	case reflect.Map:
		newval := &Map{
			Items: make(map[Object]Object, val.Len()),
		}
		for _, k := range val.MapKeys() {
			newval.Items[convert(k)] = convert(val.MapIndex(k))
		}
		return newval

	case reflect.Struct:
		newval := make(map[string]interface{})
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).CanInterface() {
				n := Fixtype(val.Field(i).Interface())
				newval[val.Type().Field(i).Name] = n
			}
		}

		for i := 0; i < val.NumMethod(); i++ {
			newval[val.Type().Method(i).Name] = convert(val.Type().Method(i).Func.Interface())
		}

		return convert(newval)

	case reflect.String:
		return String(val.String())

	case reflect.Interface:
		return convert(val.Interface())

	case reflect.Float32, reflect.Float64:
		return Number(val.Float())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return Number(float64(val.Int()))

	case reflect.Func:
		return &Func{fnc: val.Interface()}

		// TODO meh?
	case reflect.Ptr:
		//log.Println("PTR", val, in, val.Pointer(), val.Elem())
		//if val.Elem().IsValid() {
		//	if o, ok := val.Elem().Interface().(Object); !ok {
		//		return convert(o)
		//	}
		//}
		return Nil{}

	case reflect.Bool:
		return Bool(val.Bool())
	}

	panic(fmt.Sprintf("Cannot convert %#v %T %s %s", val, val, val.Type(), val.Kind()))
}

// Func

type Func struct {
	fnc interface{}
}

func (f *Func) Type() ObjectType         { return FUNC }
func (f *Func) Field(name string) Object { return Nil{} }
func (f *Func) String() string           { return fmt.Sprintf("%v", f.fnc) }

// Array

type Array struct {
	items []Object
}

func (a *Array) Type() ObjectType { return ARRAY }
func (a *Array) String() string {
	tmp := make([]string, len(a.items))
	for i, v := range a.items {
		tmp[i] = v.String()
	}
	return strings.Join(tmp, " ")
}
func (a *Array) Field(name string) Object {
	switch name {
	case "length":
		return &Func{fnc: a.Length}

	case "indexOf":
		return &Func{fnc: a.IndexOf}

	case "join":
		return &Func{fnc: a.Join}

	case "push":
		return &Func{fnc: a.Push}
	}

	panic("field not found")
}

func (a *Array) Length() Object {
	return Number(len(a.items))
}

func (a *Array) IndexOf(what interface{}) Object {
	for i, w := range a.items {
		if reflect.DeepEqual(w, what) {
			return Number(i)
		}
	}
	return Number(-1)
}

func (a *Array) Join(sep string) Object {
	var aa []string

	for _, v := range a.items {
		aa = append(aa, v.String())
	}

	return String(strings.Join(aa, sep))
}

func (a *Array) Push(what Object) Object {
	a.items = append(a.items, what)
	return Nil{}
}

// Map

type Map struct {
	Items map[Object]Object
}

func (m *Map) Type() ObjectType { return MAP }
func (m *Map) String() string {
	b, err := m.MarshalJSON()
	return err.Error() + string(b)
}
func (m *Map) Field(field string) Object {
	if field == "__assign" {
		return &Func{fnc: func(k, v interface{}) Object {
			m.Items[convert(k)] = convert(v)
			return Nil{}
		}}
	}

	if i, ok := m.Items[String(field)]; ok {
		return i
	}
	return Nil{}
}

func (m *Map) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{}, len(m.Items))
	for k, v := range m.Items {
		tmp[k.String()] = v
	}
	return json.Marshal(tmp)
}

// String

type String string

func (s String) Type() ObjectType { return STRING }
func (s String) String() string   { return string(s) }
func (s String) Field(field string) Object {
	switch field {
	case "charAt":
		return &Func{fnc: s.CharAt}
	case "toUpperCase":
		return &Func{fnc: s.ToUpperCase}
	case "slice":
		return &Func{fnc: s.Slice}
	}
	return Nil{}
}

func (s String) CharAt(pos int) string {
	if pos >= len(s) {
		return ""
	}
	return string(s[pos])
}

func (s String) ToUpperCase() string {
	return strings.ToUpper(string(s))
}

func (s String) Slice(from int) string {
	return string(s[from:])
}

// Number

type Number float64

func (n Number) Type() ObjectType    { return NUMBER }
func (n Number) Field(string) Object { return Nil{} }
func (n Number) String() string      { return fmt.Sprintf("%.0f", float64(n)) }

// Bool

type Bool bool

func (b Bool) Type() ObjectType    { return BOOL }
func (b Bool) Field(string) Object { return Nil{} }
func (b Bool) String() string      { return fmt.Sprintf("%v", bool(b)) }

// Nil

type Nil struct{}

func (n Nil) Type() ObjectType    { return NIL }
func (n Nil) Field(string) Object { return Nil{} }
func (n Nil) String() string      { return "" }
