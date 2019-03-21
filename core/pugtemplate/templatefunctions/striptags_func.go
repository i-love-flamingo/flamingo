package templatefunctions

import (
	"flamingo.me/flamingo/framework/config"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type (
	StriptagsFunc     struct{}
	AllowedAttributes map[string]bool
	AllowedTags       []AllowedTag
	AllowedTag        struct {
		name       string
		attributes AllowedAttributes
	}
)

var nameRe = regexp.MustCompile(`[a-z0-9\-]+`)
var attributesRe = regexp.MustCompile(`[a-z0-9]+([a-z]+)`)

func createTag(definition string) AllowedTag {
	definition = strings.ToLower(definition)
	attributes := make(AllowedAttributes)
	for _, attr := range attributesRe.FindAllString(definition, -1) {
		attributes[attr] = true
	}

	return AllowedTag{
		nameRe.FindString(definition),
		attributes,
	}
}

func (at AllowedTags) Find(tagName string) *AllowedTag {
	for _, tag := range at {
		if tag.name == tagName {
			return &tag
		}
	}
	return nil
}

func (at AllowedTags) Contains(tagName string) bool {
	if tag := at.Find(tagName); tag != nil {
		return true
	}
	return false
}

// Func as implementation of debug method
func (df StriptagsFunc) Func() interface{} {
	return func(htmlString string, allowedTagsConfig ...config.Slice) string {
		doc, err := html.ParseFragment(strings.NewReader(htmlString), nil)
		if err != nil {
			return ""
		}

		var allowedTags AllowedTags
		if len(allowedTagsConfig) == 1 {
			for _, item := range allowedTagsConfig[0] {
				if definition, ok := item.(string); ok {
					allowedTags = append(allowedTags, createTag(definition))
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

func cleanTags(n *html.Node, allowedTags AllowedTags) string {
	res := ""

	if n.Type == html.ElementNode && allowedTags.Contains(n.Data) {
		res += "<"
		res += n.Data
		res += getAllowedAttributes(n.Attr, allowedTags.Find(n.Data).attributes)
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

	if n.Type == html.ElementNode && allowedTags.Contains(n.Data) {
		res += "</" + n.Data + ">"
	}

	return res
}

func getAllowedAttributes(attributes []html.Attribute, allowedAttributes AllowedAttributes) string {
	res := ""
	for _, attr := range attributes {
		if allowedAttributes[attr.Key] {
			res += " " + attr.Key + "=\"" + html.EscapeString(attr.Val) + "\""
		}
	}
	return res
}
