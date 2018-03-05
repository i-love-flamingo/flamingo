package router

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathHandling(t *testing.T) {
	t.Run("Path Handling", func(t *testing.T) {
		t.Run("Incoming Request", func(t *testing.T) {
			var path = NewPath(`/path/to/:something/$id<[0-9]+>/*foo`)

			t.Run("Should properly parse request path's", func(t *testing.T) {
				assert.NotNil(t, path)
				var match = path.Match(`/path/to/something123/445566/foo/bar`)
				assert.NotNil(t, match)
				assert.Equal(t, "something123", match.Values["something"])
				assert.Equal(t, "445566", match.Values["id"])
				assert.Equal(t, "foo/bar", match.Values["foo"])
			})

			t.Run("Should properly render the path", func(t *testing.T) {
				p, err := path.Render(map[string]string{
					"something": "something123",
					"id":        "445566",
					"foo":       "foo/bar",
				}, map[string]struct{}{})

				assert.NoError(t, err)
				assert.Equal(t, `/path/to/something123/445566/foo/bar`, p)

				_, err = path.Render(map[string]string{
					"something": "something123",
					"id":        "aaa",
					"foo":       "foo/bar",
				}, map[string]struct{}{})
				assert.EqualError(t, err, `param id in wrong format`)
			})
		})

		t.Run("Fixed part matching", func(t *testing.T) {
			t.Run("Should match the whole part", func(t *testing.T) {
				var path = NewPath(`/path/to/something`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/something`))
				assert.Empty(t, path.Match(`/path/to/something`).Values)
				assert.Nil(t, path.Match(`/pat/to/something`))
			})
		})

		t.Run("Param part matching", func(t *testing.T) {
			t.Run("Should match the different parameters", func(t *testing.T) {
				var path = NewPath(`/path/to/:something/:else/end`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/something/else/end`))
				assert.Equal(t, "something", path.Match(`/path/to/something/else/end`).Values["something"])
				assert.Equal(t, "else", path.Match(`/path/to/something/else/end`).Values["else"])
				assert.Equal(t, "foo", path.Match(`/path/to/foo/bar/end`).Values["something"])
				assert.Equal(t, "bar", path.Match(`/path/to/foo/bar/end`).Values["else"])

				assert.Nil(t, path.Match(`/path/to/foo/bar`))
				assert.Equal(t, "", path.Match(`/path/to///end`).Values["something"])
				assert.Equal(t, "", path.Match(`/path/to///end`).Values["else"])
			})

			t.Run("Edge-case at the end", func(t *testing.T) {
				var path = NewPath(`/path/to/:something`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/something`))
				assert.Equal(t, "something", path.Match(`/path/to/something`).Values["something"])
				assert.Equal(t, "foo", path.Match(`/path/to/foo`).Values["something"])
				assert.Nil(t, path.Match(`/path/to`))
				assert.NotNil(t, path.Match(`/path/to/`))
				assert.Equal(t, "", path.Match(`/path/to/`).Values["something"])
				assert.Equal(t, "", path.Match(`/path/to//`).Values["something"])
				assert.Nil(t, path.Match(`/path/to///`))
			})
		})

		t.Run("Wildcard part matching", func(t *testing.T) {
			t.Run("Should match the wildcard", func(t *testing.T) {
				var path = NewPath(`/path/to/*something`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/foo/bar`))
				assert.Equal(t, "foo/bar", path.Match(`/path/to/foo/bar`).Values["something"])
				assert.NotNil(t, path.Match(`/path/to/`))
				assert.Equal(t, "", path.Match(`/path/to/`).Values["something"])
				assert.NotNil(t, path.Match(`/path/to/foo/bar/`))
				assert.Equal(t, "foo/bar/", path.Match(`/path/to/foo/bar/`).Values["something"])
			})
		})

		t.Run("Regex part matching", func(t *testing.T) {
			t.Run("Should find the regexp's", func(t *testing.T) {
				var path = NewPath(`/path/to/$<[0-9]+>/$id<[0-9]+>`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/10/20`))
				assert.Equal(t, "20", path.Match(`/path/to/10/20`).Values["id"])
				assert.Len(t, path.Match(`/path/to/10/20`).Values, 1)
				assert.Nil(t, path.Match(`/path/to/10/`))
				assert.Nil(t, path.Match(`/path/to//`))
				assert.Nil(t, path.Match(`/path/to//20`))

				path = NewPath(`/path/to/$id<[0-9]+>/$id2<[0-9]+>`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/10/20`))
				assert.Equal(t, "10", path.Match(`/path/to/10/20`).Values["id"])
				assert.Equal(t, "20", path.Match(`/path/to/10/20`).Values["id2"])
				assert.Len(t, path.Match(`/path/to/10/20`).Values, 2)
			})

			t.Run("Should match fullpath regex", func(t *testing.T) {
				var path = NewPath(`/path/to/$foo<.*>`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/path/to/10/20`))
				assert.Equal(t, "10/20", path.Match(`/path/to/10/20`).Values["foo"])
				assert.Equal(t, "bla", path.Match(`/path/to/bla`).Values["foo"])
				assert.Equal(t, "", path.Match(`/path/to/`).Values["foo"])
			})
		})

		t.Run("Edge cases", func(t *testing.T) {
			t.Run("Should match /", func(t *testing.T) {
				var path = NewPath(`/`)
				assert.NotNil(t, path)
				assert.NotNil(t, path.Match(`/`))
				assert.Empty(t, path.Match(`/`).Values)
			})
		})

		t.Run("Xml Sitemap Case", func(t *testing.T) {
			var path = NewPath(`/sitemap/$id<product-(\d+).xml>`)
			assert.NotNil(t, path)
			assert.NotNil(t, path.Match(`/sitemap/product-1.xml`))
			assert.NotNil(t, path.Match(`/sitemap/product-2.xml`))
			assert.NotNil(t, path.Match(`/sitemap/product-100.xml`))
			assert.Nil(t, path.Match(`/sitemap/product-test.xml`))
			assert.NotEmpty(t, path.Match(`/sitemap/product-1.xml`).Values)
			values := path.Match(`/sitemap/product-1.xml`).Values
			id := strings.TrimPrefix(strings.TrimSuffix(values["id"], ".xml"), "product-")
			assert.Equal(t, id, "1")
		})
	})
}
