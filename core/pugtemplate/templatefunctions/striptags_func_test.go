package templatefunctions_test

import (
	"testing"

	"flamingo.me/flamingo/core/pugtemplate/templatefunctions"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/template"
	"github.com/stretchr/testify/assert"
)

func TestStriptagsFunc(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		out         string
		allowedTags config.Slice
	}{
		{"should keep plain text", "do not modify me", "do not modify me", config.Slice{}},
		{"should keep linebreaks", "Hello\nWorld", "Hello\nWorld", config.Slice{}},
		{"should remove tags by default", "<h1>Headline<h1> <p>Paragraph</p>", "Headline Paragraph", config.Slice{}},
		{"should keep defined tags", "<h1>Headline</h1>", "<h1>Headline</h1>", config.Slice{"h1", "h2"}},
		{
			"should remove non whitelisted attributes",
			"<h1 style=\"font-size: 500px\">Keep me</h1><script src=\"http://miner.tld/x.js\">",
			"<h1>Keep me</h1>",
			config.Slice{"h1"},
		},
		{
			"should keep whitelisted attributes",
			"<p>I'm a paragraph containing a <a href=\"http://tld.com\" style=\"font-size:100px\">link</a></p>",
			"<p>I'm a paragraph containing a <a href=\"http://tld.com\">link</a></p>",
			config.Slice{"p", "a(href)"},
		},
		{
			"should keep multiple whitelisted attributes",
			"<a href=\"http://domain.tld\" target=\"_blank\" rel=\"nofollow\">Link with target</a>",
			"<a href=\"http://domain.tld\" target=\"_blank\">Link with target</a>",
			config.Slice{"a(href target)"},
		},
	}

	var stripTagsFunc template.Func = new(templatefunctions.StriptagsFunc)
	stripTags := stripTagsFunc.Func().(func(htmlString string, allowedTagsConfig ...config.Slice) string)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, stripTags(tt.in, tt.allowedTags), tt.name)
		})
	}
}
