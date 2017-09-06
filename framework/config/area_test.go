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
