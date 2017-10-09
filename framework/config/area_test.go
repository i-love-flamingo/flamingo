package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := make(Map)

	m.Add(Map{
		"foo": "bar",
	})

	m.Add(Map{
		"foo": "aaa",
	})

	assert.Equal(t, "aaa", m["foo"])

	assert.Panics(t, func() {
		m.Add(Map{
			"foo.bar": "a",
		})
	})

	m.Add(Map{
		"b.a": "a",
		"b.b": "b",
	})

	assert.Equal(t, Map{"a": "a", "b": "b"}, m["b"])

	m.Add(Map{
		"b": Map{
			"a": "a",
			"b": "b",
		},
	})

	m.Add(Map{
		"b.a": "c",
	})

	assert.Equal(t, "c", m["b"].(Map)["a"])

	assert.Panics(t, func() {
		m.Add(Map{
			"b": "a",
		})
	})

	m.Add(Map{
		"x": Map{
			"x":   "x",
			"y.z": "a",
		},
	})

	assert.Equal(t, "x", m["x"].(Map)["x"])
	assert.Equal(t, Map{"z": "a"}, m["x"].(Map)["y"])
}

func TestMapMarshalTo(t *testing.T) {
	type resultType struct {
		Key    string
		Number int
		Flag   bool
		Map    map[string]interface{}
		Sub    struct {
			Foo    string
			Subsub struct {
				Bar string
			}
		}
	}

	//fill the config map according to the resultType struct
	m := make(Map)

	m.Add(Map{
		"key":    "value",
		"number": "5",
		"flag":   true,
	})
	m.Add(Map{
		"sub.foo":        "baz",
		"sub.subsub.bar": "myvalue",
	})
	m.Add(Map{
		"map.a":   "a",
		"map.b":   "b",
		"map.c":   "c",
		"map.d.a": "da",
		"map.d.b": "db",
		"map.e":   "e",
		"map.f":   "f",
	})


	var result resultType

	err := m.MarshalTo(&result)
	assert.Nil(t, err)

	assert.Equal(
		t,
		resultType{
			Key:    "value",
			Number: 5,
			Flag:   true,
			Map: map[string]interface{}{
				"a": "a",
				"b": "b",
				"c": "c",
				"d": map[string]interface{}{
					"a": "da",
					"b": "db",
				},
				"e": "e",
				"f": "f",
			},
			Sub: struct {
				Foo    string
				Subsub struct{ Bar string }
			}{
				Foo:    "baz",
				Subsub: struct{ Bar string }{Bar: "myvalue"},
			},
		},
		result,
		"result of marshalling is wrong",
	)

}
