package web

import (
	"net/url"
	"strings"
)

// URLTitle normalizes a title for nice usage in URLs
func URLTitle(title string) string {
	newTitle := strings.ToLower(strings.Replace(title, " ", "-", -1))
	return url.QueryEscape(newTitle)
}
