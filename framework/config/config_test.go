package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/stretchr/testify/assert"
)

func readConfig(t *testing.T, configName string) Map {
	config, err := os.ReadFile(configName)
	assert.NoError(t, err)

	config = []byte(envRegex.ReplaceAllStringFunc(
		string(config),
		func(a string) string { return os.Getenv(envRegex.FindStringSubmatch(a)[1]) },
	))

	cfg := make(Map)
	err = yaml.Unmarshal(config, &cfg)
	assert.NoError(t, err)
	return cfg
}

func TestNilValuesRemoveData(t *testing.T) {
	config := make(Map)

	cfg := readConfig(t, "testdata/valid/config.yml")
	assert.NoError(t, config.Add(cfg))
	cfg = readConfig(t, "testdata/valid/config_dev.yml")
	assert.NoError(t, config.Add(cfg))

	fooValue, present := config.Get("foo")
	assert.Equal(t, nil, fooValue)
	assert.Equal(t, true, present)
}

func TestMapDeepmerge(t *testing.T) {
	subTestMapDeepmerge := func(t *testing.T) {
		m := make(Map)

		assert.NoError(t, m.Add(Map{
			"foo.bar": "bar",
		}))

		assert.NoError(t, m.Add(Map{
			"foo.bar": "bar2",
		}))

		assert.Equal(t, Map{
			"foo": Map{
				"bar": "bar2",
			},
		}, m)
	}

	t.Run("run 1", subTestMapDeepmerge)
	t.Run("run 2", subTestMapDeepmerge)
	t.Run("run 3", subTestMapDeepmerge)
	t.Run("run 4", subTestMapDeepmerge)
	t.Run("run 5", subTestMapDeepmerge)
}

func TestMap(t *testing.T) {
	m := make(Map)

	assert.NoError(t, m.Add(Map{
		"foo": "bar",
	}))

	assert.NoError(t, m.Add(Map{
		"foo": "aaa",
	}))

	assert.Equal(t, "aaa", m["foo"])

	assert.Error(t, m.Add(Map{
		"foo.bar": "a",
	}))

	assert.NoError(t, m.Add(Map{
		"b.a": "a",
		"b.b": "b",
	}))

	assert.Equal(t, Map{"a": "a", "b": "b"}, m["b"])

	assert.NoError(t, m.Add(Map{
		"b": Map{
			"a": "a",
			"b": "b",
		},
	}))

	assert.NoError(t, m.Add(Map{
		"b.a": "c",
	}))

	assert.Equal(t, "c", m["b"].(Map)["a"])

	assert.Error(t, m.Add(Map{
		"b": "a",
	}))

	assert.NoError(t, m.Add(Map{
		"x": Map{
			"x":   "x",
			"y.z": "a",
		},
	}))

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

	assert.NoError(t, m.Add(Map{
		"key":    "value",
		"number": 5,
		"flag":   true,
	}))
	assert.NoError(t, m.Add(Map{
		"sub.foo":        "baz",
		"sub.subsub.bar": "myvalue",
	}))
	assert.NoError(t, m.Add(Map{
		"map.a":   "a",
		"map.b":   "b",
		"map.c":   "c",
		"map.d.a": "da",
		"map.d.b": "db",
		"map.e":   "e",
		"map.f":   "f",
	}))

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
	assert.NoError(t, m.Add(Map{
		"foo.bar.x.y.z": "test",
	}))

	val, ok := m.Get("foo.bar.x.y.z")
	assert.True(t, ok)
	assert.Equal(t, "test", val)

	val, ok = m.Get("foo.bar")
	assert.True(t, ok)
	assert.Equal(t, Map{"x": Map{"y": Map{"z": "test"}}}, val)

	_, ok = m.Get("foo.bar.baz")
	assert.False(t, ok)

	_, ok = m.Get("unknown")
	assert.False(t, ok)
}

func TestMap_Flat(t *testing.T) {
	tests := []struct {
		name      string
		m         Map
		overwrite Map
		want      Map
	}{
		{
			name: "overwrite",
			m: Map{
				"tri.tra":     "tral",
				"foo.bar.baz": "DEFAULT",
				"foo.bar.bam": "",
			},
			overwrite: Map{"foo.bar.baz": "OVERWRITE"},
			want: Map{
				"foo":         Map{"bar": Map{"bam": "", "baz": "OVERWRITE"}},
				"foo.bar":     Map{"bam": "", "baz": "OVERWRITE"},
				"foo.bar.baz": "OVERWRITE",
				"foo.bar.bam": "",
				"tri":         Map{"tra": "tral"},
				"tri.tra":     "tral",
			},
		},
	}
	for _, tt := range tests {
		// run each case multiple times
		for i := 0; i < 20; i++ {
			t.Run(tt.name, func(t *testing.T) {
				assert.NoError(t, tt.m.Add(tt.overwrite))

				got := tt.m.Flat()

				if len(got) != len(tt.want) {
					t.Errorf("number of entries different, got %d, want %d", len(got), len(tt.want))
				}

				for key, value := range got {
					want, found := tt.want[key]
					if !found {
						t.Errorf("Key %v is missing in expected data", key)
					}
					if !reflect.DeepEqual(value, want) {
						t.Errorf("key %v is %v, want %v", key, value, want)
					}
				}
			})
		}
	}
}

func TestMap_Add(t *testing.T) {
	tests := []struct {
		name string
		m    Map
		add  Map
		want Map
	}{
		{
			name: "overwrite",
			m: Map{
				"tri.tra":     "tral",
				"foo.bar.baz": "DEFAULT",
				"foo.bar.bam": "",
			},
			add: Map{"foo.bar.baz": "OVERWRITE"},
			want: Map{
				"foo": Map{"bar": Map{"bam": "", "baz": "OVERWRITE"}},
				"tri": Map{"tra": "tral"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, tt.m.Add(tt.add))

			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("got %v, want %v", tt.m, tt.want)
			}
		})
	}
}
