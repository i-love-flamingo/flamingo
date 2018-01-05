package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func SubTestMapDeepmerge(t *testing.T) {
	m := make(Map)

	m.Add(Map{
		"foo.bar": "bar",
	})

	m.Add(Map{
		"foo.bar": "bar2",
	})

	assert.Equal(t, Map{
		"foo": Map{
			"bar": "bar2",
		},
	}, m)
}

func TestNilValuesRemoveData(t *testing.T) {
	config := make(Map)

	cfg := readConfig(t, "test/config.yml")
	config.Add(cfg)
	cfg = readConfig(t, "test/config_dev.yml")
	config.Add(cfg)

	assert.Equal(t, Map{"foo": nil}, config)
}

func readConfig(t *testing.T, configName string) Map {
	config, err := ioutil.ReadFile(configName)
	assert.NoError(t, err)

	config = []byte(regex.ReplaceAllStringFunc(
		string(config),
		func(a string) string { return os.Getenv(regex.FindStringSubmatch(a)[1]) },
	))

	cfg := make(Map)
	err = yaml.Unmarshal(config, &cfg)
	assert.NoError(t, err)
	return cfg
}

func TestMapDeepmerge(t *testing.T) {
	t.Run("run 1", SubTestMapDeepmerge)
	t.Run("run 2", SubTestMapDeepmerge)
	t.Run("run 3", SubTestMapDeepmerge)
	t.Run("run 4", SubTestMapDeepmerge)
	t.Run("run 5", SubTestMapDeepmerge)
}

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

func TestMapMapInto(t *testing.T) {
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
		"number": 5,
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

	err := m.MapInto(&result)
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

func TestMap_Get(t *testing.T) {
	m := make(Map)
	m.Add(Map{
		"foo.bar.x.y.z": "test",
	})

	val, ok := m.Get("foo.bar.x.y.z")
	assert.True(t, ok)
	assert.Equal(t, "test", val)

	val, ok = m.Get("foo.bar")
	assert.True(t, ok)
	assert.Equal(t, Map{"x": Map{"y": Map{"z": "test"}}}, val)

	val, ok = m.Get("foo.bar.baz")
	assert.False(t, ok)

	val, ok = m.Get("unknown")
	assert.False(t, ok)
}
