package templatefunctions

import (
	"strings"

	"golang.org/x/net/html"
)

type (
	StriptagsFunc struct{}
)

// Func as implementation of debug method
func (df StriptagsFunc) Func() interface{} {
	return func(htmlString string) string {
		doc, err := html.ParseFragment(strings.NewReader(htmlString), nil)
		if err != nil {
			return ""
		}

		res := ""
		for _, n := range doc {
			res += removeTags(n)
		}
		return res
	}
}

func removeTags(n *html.Node) string {
	res := ""

	if n.Type == html.TextNode {
		res += n.Data
	}
	if n.FirstChild != nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			res += removeTags(c)
		}
	}

	return res
}
