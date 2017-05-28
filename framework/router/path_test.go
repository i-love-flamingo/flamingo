package router

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Path Test", func() {
	Context("Path Handling", func() {
		Context("Incoming Request", func() {
			var path = NewPath(`/path/to/:something/$id<[0-9]+>/*foo`)

			It("Should properly parse request path's", func() {
				Expect(path).ToNot(BeNil())
				var match = path.Match(`/path/to/something123/445566/foo/bar`)
				Expect(match).ToNot(BeNil())
				Expect(match.Values).To(HaveKeyWithValue("something", "something123"))
				Expect(match.Values).To(HaveKeyWithValue("id", "445566"))
				Expect(match.Values).To(HaveKeyWithValue("foo", "foo/bar"))
			})

			It("Should properly render the path", func() {
				Expect(path.Render(map[string]string{
					"something": "something123",
					"id":        "445566",
					"foo":       "foo/bar",
				})).To(Equal(`/path/to/something123/445566/foo/bar`))

				_, err := path.Render(map[string]string{
					"something": "something123",
					"id":        "aaa",
					"foo":       "foo/bar",
				})
				Expect(err).To(MatchError(`param id in wrong format`))
			})
		})

		Context("Fixed part matching", func() {
			It("Should match the whole part", func() {
				var path = NewPath(`/path/to/something`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/something`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/something`).Values).To(BeEmpty())
				Expect(path.Match(`/pat/to/something`)).To(BeNil())
			})
		})

		Context("Param part matching", func() {
			It("Should match the different parameters", func() {
				var path = NewPath(`/path/to/:something/:else/end`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/something/else/end`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/something/else/end`).Values).To(HaveKeyWithValue("something", "something"))
				Expect(path.Match(`/path/to/something/else/end`).Values).To(HaveKeyWithValue("else", "else"))
				Expect(path.Match(`/path/to/foo/bar/end`).Values).To(HaveKeyWithValue("something", "foo"))
				Expect(path.Match(`/path/to/foo/bar/end`).Values).To(HaveKeyWithValue("else", "bar"))
				Expect(path.Match(`/path/to/foo/bar`)).To(BeNil())
				Expect(path.Match(`/path/to///end`).Values).To(HaveKeyWithValue("something", ""))
				Expect(path.Match(`/path/to///end`).Values).To(HaveKeyWithValue("else", ""))
			})

			It("Edge-case at the end", func() {
				var path = NewPath(`/path/to/:something`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/something`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/something`).Values).To(HaveKeyWithValue("something", "something"))
				Expect(path.Match(`/path/to/foo`).Values).To(HaveKeyWithValue("something", "foo"))
				Expect(path.Match(`/path/to`)).To(BeNil())
				Expect(path.Match(`/path/to/`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/`).Values).To(HaveKeyWithValue("something", ""))
				Expect(path.Match(`/path/to//`).Values).To(HaveKeyWithValue("something", ""))
				Expect(path.Match(`/path/to///`)).To(BeNil())
			})
		})

		Context("Wildcard part matching", func() {
			It("Should match the wildcard", func() {
				var path = NewPath(`/path/to/*something`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/foo/bar`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/foo/bar`).Values).To(HaveKeyWithValue("something", "foo/bar"))
				Expect(path.Match(`/path/to/`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/`).Values).To(HaveKeyWithValue("something", ""))
				Expect(path.Match(`/path/to/foo/bar/`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/foo/bar/`).Values).To(HaveKeyWithValue("something", "foo/bar/"))
			})
		})

		Context("Regex part matching", func() {
			It("Should find the regexp's", func() {
				var path = NewPath(`/path/to/$<[0-9]+>/$id<[0-9]+>`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveKeyWithValue("id", "20"))
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveLen(1))
				Expect(path.Match(`/path/to/10/`)).To(BeNil())
				Expect(path.Match(`/path/to//`)).To(BeNil())
				Expect(path.Match(`/path/to//20`)).To(BeNil())

				path = NewPath(`/path/to/$id<[0-9]+>/$id2<[0-9]+>`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveKeyWithValue("id", "10"))
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveKeyWithValue("id2", "20"))
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveLen(2))
			})

			It("Should match fullpath regex", func() {
				var path = NewPath(`/path/to/$foo<.*>`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`)).ToNot(BeNil())
				Expect(path.Match(`/path/to/10/20`).Values).To(HaveKeyWithValue("foo", "10/20"))
				Expect(path.Match(`/path/to/bla`).Values).To(HaveKeyWithValue("foo", "bla"))
				Expect(path.Match(`/path/to/`).Values).To(HaveKeyWithValue("foo", ""))
			})
		})

		Context("Edge cases", func() {
			It("Should match /", func() {
				var path = NewPath(`/`)
				Expect(path).ToNot(BeNil())
				Expect(path.Match(`/`)).ToNot(BeNil())
				Expect(path.Match(`/`).Values).To(BeEmpty())
			})
		})
	})
})
