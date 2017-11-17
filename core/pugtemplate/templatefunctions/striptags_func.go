package templatefunctions

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

type (
	StriptagsFunc struct{}
)

// Name alias for use in template
func (df StriptagsFunc) Name() string {
	return "striptags"
}

// Func as implementation of debug method
func (df StriptagsFunc) Func() interface{} {
	return func(htmlString string) string {
		doc, err := html.Parse(strings.NewReader(htmlString))
		if err != nil {
			return ""
		}

		removeScript(doc)
		buf := new(bytes.Buffer)
		if err := html.Render(buf, doc); err != nil {
			return ""
		}
		return buf.String()
	}
}

func removeScript(n *html.Node) {
	// if note is script tag
	if n.Type == html.ElementNode && n.Data == "script" {
		n.Parent.RemoveChild(n)
		return // script tag is gone...
	}
	// traverse DOM
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeScript(c)
	}
}
