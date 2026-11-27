package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/config"
)

type brandTarget struct {
	Name string `json:"title"`
	Size int
}

// MapInto follows encoding/json.
func TestMapInto_HonoursJSONTagsAndCaseInsensitiveNames(t *testing.T) {
	t.Parallel()

	entry := config.Map{"title": "x", "SIZE": 2}

	t.Run("map", func(t *testing.T) {
		t.Parallel()

		var got brandTarget

		require.NoError(t, entry.MapInto(&got))
		assert.Equal(t, brandTarget{Name: "x", Size: 2}, got)
	})

	t.Run("slice of structs", func(t *testing.T) {
		t.Parallel()

		var got []brandTarget

		require.NoError(t, config.Slice{entry}.MapInto(&got))
		assert.Equal(t, []brandTarget{{Name: "x", Size: 2}}, got)
	})

	t.Run("slice of strings", func(t *testing.T) {
		t.Parallel()

		var got []string

		require.NoError(t, config.Slice{"a", "b"}.MapInto(&got))
		assert.Equal(t, []string{"a", "b"}, got)
	})
}
