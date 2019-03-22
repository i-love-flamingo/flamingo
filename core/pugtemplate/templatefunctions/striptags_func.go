package templatefunctions

import (
	"flamingo.me/flamingo/framework/config"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

type (
	StriptagsFunc     struct{}
	allowedAttributes map[string]struct{}
	allowedTags       map[string]allowedTag
	allowedTag        struct {
		name       string
		attributes allowedAttributes
	}
)

var (
	nameRe       = regexp.MustCompile(`[a-z0-9\-]+`)
	attributesRe = regexp.MustCompile(`[a-z0-9]+([a-z]+)`)
)

func createTag(definition string) allowedTag {
	definition = strings.ToLower(definition)
	attributes := make(allowedAttributes)
	for _, attr := range attributesRe.FindAllString(definition, -1) {
		attributes[attr] = struct{}{}
	}

	return allowedTag{
		nameRe.FindString(definition),
		attributes,
	}
}

// Func as implementation of debug method
func (df StriptagsFunc) Func() interface{} {
	return func(htmlString string, allowedTagsConfig ...config.Slice) string {
		doc, err := html.ParseFragment(strings.NewReader(htmlString), nil)
		if err != nil {
			return ""
		}

		allowedTags := make(allowedTags)
		if len(allowedTagsConfig) == 1 {
			for _, item := range allowedTagsConfig[0] {
				if definition, ok := item.(string); ok {
					tag := createTag(definition)
					allowedTags[tag.name] = tag
				}
			}
		}

		res := ""
		for _, n := range doc {
			res += cleanTags(n, allowedTags)
		}
		return res
	}
}

func cleanTags(n *html.Node, allowedTags allowedTags) string {
	var allowedTag allowedTag
	res := ""

	if n.Type == html.ElementNode {
		if tag, ok := allowedTags[n.Data]; ok {
			allowedTag = tag
		}
	}

	if allowedTag.name != "" {
		res += "<"
		res += n.Data
		res += getAllowedAttributes(n.Attr, allowedTag.attributes)
		res += ">"
	}

	if n.Type == html.TextNode {
		res += n.Data
	}

	if n.FirstChild != nil {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			res += cleanTags(c, allowedTags)
		}
	}

	if allowedTag.name != "" {
		res += "</" + n.Data + ">"
	}

	return res
}

func getAllowedAttributes(attributes []html.Attribute, allowedAttributes allowedAttributes) string {
	res := ""
	for _, attr := range attributes {
		if _, ok := allowedAttributes[attr.Key]; ok {
			res += " " + attr.Key + "=\"" + html.EscapeString(attr.Val) + "\""
		}
	}
	return res
}
