package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	root := new(Area)

	require.NoError(t, os.Setenv("TEST1", "test-value"))
	require.NoError(t, os.Setenv("TEST4", "injected"))

	t.Run("config dir does not exist", func(t *testing.T) {
		err := Load(root, "not-existing")
		require.NoError(t, err)
		assert.Equal(t, Map{"area": ""}, root.Configuration.Flat())
	})

	t.Run("valid config files", func(t *testing.T) {
		err := Load(root, "testdata/valid")
		assert.NoError(t, err)
		assert.Contains(t, root.Configuration.Flat(), "area")
		assert.Contains(t, root.Configuration.Flat(), "foo")
		assert.Contains(t, root.Configuration.Flat(), "foo.bar")
		assert.Contains(t, root.Configuration.Flat(), "foo.bar.test")
		assert.Equal(t, Shim(1.0, true), Shim(root.Configuration.Get("foo.bar.test")))

		assert.Equal(t, Shim("test-value", true), Shim(root.Configuration.Get("env.var.test1")))
		assert.Equal(t, Shim(nil, true), Shim(root.Configuration.Get("env.var.test2")))
		assert.Equal(t, Shim("default", true), Shim(root.Configuration.Get("env.var.test3")))
		assert.Equal(t, Shim("injected", true), Shim(root.Configuration.Get("env.var.test4")))
	})

	t.Run("valid config files with dev context", func(t *testing.T) {
		require.NoError(t, os.Setenv("CONTEXT", "dev"))
		err := Load(root, "testdata/valid")
		assert.NoError(t, err)
		assert.Contains(t, root.Configuration.Flat(), "area")
		assert.Contains(t, root.Configuration.Flat(), "foo")
		assert.NotContains(t, root.Configuration.Flat(), "foo.bar")

		assert.Equal(t, Shim(nil, true), Shim(root.Configuration.Get("foo")))
	})

	t.Run("valid config files with additional config", func(t *testing.T) {
		require.NoError(t, flagSet.Set("flamingo-config", "baz: bam"))
		require.NoError(t, flagSet.Set("flamingo-config", "foo.bar.test: 'hello'"))
		err := Load(root, "testdata/valid")
		assert.NoError(t, err)
		assert.Contains(t, root.Configuration.Flat(), "area")
		assert.Contains(t, root.Configuration.Flat(), "foo")
		assert.Contains(t, root.Configuration.Flat(), "foo.bar")
		assert.Contains(t, root.Configuration.Flat(), "foo.bar.test")
		assert.Contains(t, root.Configuration.Flat(), "baz")

		assert.Equal(t, Shim("hello", true), Shim(root.Configuration.Get("foo.bar.test")))
		assert.Equal(t, Shim("bam", true), Shim(root.Configuration.Get("baz")))
	})

	t.Run("invalid config file", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = Load(root, "testdata/invalid")
		})
	})

	t.Run("valid config file with invalid additional config", func(t *testing.T) {
		assert.Panics(t, func() {
			require.NoError(t, flagSet.Set("flamingo-config", "baz"))

			_ = Load(root, "testdata/valid")
		})
	})

}

func Shim(a, b interface{}) []interface{} {
	return []interface{}{a, b}
}
