package config

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	root := new(Area)

	os.Setenv("TEST1", "test-value")
	os.Setenv("TEST4", "injected")

	err := Load(root, "not-existing")
	assert.NoError(t, err)
	assert.Equal(t, Map{"area": ""}, root.Configuration.Flat())

	err = Load(root, "test")
	assert.NoError(t, err)
	assert.Contains(t, root.Configuration.Flat(), "area")
	assert.Contains(t, root.Configuration.Flat(), "foo")
	assert.Contains(t, root.Configuration.Flat(), "foo.bar")
	assert.Contains(t, root.Configuration.Flat(), "foo.bar.test")

	assert.Equal(t, Shim("test-value", true) , Shim(root.Configuration.Get("env.var.test1")))
	assert.Equal(t, Shim(nil, true), Shim(root.Configuration.Get("env.var.test2")))
	assert.Equal(t, Shim("default", true), Shim(root.Configuration.Get("env.var.test3")))
	assert.Equal(t, Shim("injected", true), Shim(root.Configuration.Get("env.var.test4")))

	os.Setenv("CONTEXT", "dev")
	err = Load(root, "test")
	assert.NoError(t, err)
	assert.Contains(t, root.Configuration.Flat(), "area")
	assert.Contains(t, root.Configuration.Flat(), "foo")
}

func Shim(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}
