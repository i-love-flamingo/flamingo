package templatefunctions_test

import (
	"flamingo.me/flamingo/core/pugtemplate/templatefunctions"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/template"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStriptagsFunc(t *testing.T) {

	var striptagsFunc template.Func = new(templatefunctions.StriptagsFunc)
	striptags := striptagsFunc.Func().(func(htmlString string, allowedTagsConfig ...config.Slice) string)

	// basic test without parameter
	// should remove all tags and keep only text
	assert.Equal(t, "should be ok", striptags("should be ok"))
	assert.Equal(t, "Headline\nHello world", striptags("<h1>Headline</h1>\n<p>Hello world</p>"))

	// advanced test with parameter
	// should only keep wanted tags
	var res string
	res = striptags("<h1>Hello</h1>", config.Slice{"h1", "h2", "h3", "p"})
	assert.Equal(t, "<h1>Hello</h1>", res)

	res = striptags("<h1>Remove Headline</h1>", config.Slice{"p"})
	assert.Equal(t, "Remove Headline", res)

	res = striptags("<h1>Remove Headline</h1><p>Keep paragraphs</p>", config.Slice{"p"})
	assert.Equal(t, "Remove Headline<p>Keep paragraphs</p>", res)

	// remove invalid attributes and tags
	res = striptags("<h1 style=\"font-size: 500px\">Remove Scripts</h1><script src=\"http://miner.tld/x.js\">", config.Slice{"h1"})
	assert.Equal(t, "<h1>Remove Scripts</h1>", res, "should remove invalid attributes")

	//keep allowed attributes
	res = striptags("<p>I'm a paragraph containing a <a href=\"http://tld.com\">link</a></p>", config.Slice{"p", "a(href link rel)"})
	assert.Equal(t, "<p>I'm a paragraph containing a <a href=\"http://tld.com\">link</a></p>", res, "should remove invalid attributes")
}
