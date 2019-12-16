package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParents(t *testing.T) {
	area := NewArea("root", nil, NewArea("c1", nil))
	assert.NoError(t, area.Configuration.Add(Map{"key1": "1"}))
	assert.NoError(t, area.Configuration.Add(Map{"key2": "2"}))

	child := area.Childs[0]
	assert.NoError(t, child.Configuration.Add(Map{"key1": "c1"}))

	assert.True(t, area.HasConfigKey("key1"))
	assert.True(t, area.HasConfigKey("key2"))
	assert.False(t, area.HasConfigKey("key3"))
	assert.True(t, child.HasConfigKey("key1"))
	assert.True(t, child.HasConfigKey("key2"))
	assert.False(t, child.HasConfigKey("key3"))

	v, ok := area.Config("key1")
	assert.True(t, ok)
	assert.Equal(t, "1", v)

	v, ok = area.Config("key2")
	assert.True(t, ok)
	assert.Equal(t, "2", v)

	_, ok = area.Config("key3")
	assert.False(t, ok)

	v, ok = child.Config("key1")
	assert.True(t, ok)
	assert.Equal(t, "c1", v)

	v, ok = child.Config("key2")
	assert.True(t, ok)
	assert.Equal(t, "2", v)

	_, ok = child.Config("key3")
	assert.False(t, ok)
}
