package config

import (
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	root := new(Area)

	err := Load(root, "not-existing")
	assert.NoError(t, err)
	assert.Equal(t, Map{"area": ""}, root.Configuration.Flat())

	err = Load(root, "test")
	assert.NoError(t, err)
	assert.Contains(t, root.Configuration.Flat(), "area")
	assert.Contains(t, root.Configuration.Flat(), "foo")
	assert.Contains(t, root.Configuration.Flat(), "foo.bar")
	assert.Contains(t, root.Configuration.Flat(), "foo.bar.test")

	os.Setenv("CONTEXT", "dev")
	err = Load(root, "test")
	assert.NoError(t, err)
	assert.Contains(t, root.Configuration.Flat(), "area")
	assert.Contains(t, root.Configuration.Flat(), "foo")
}
