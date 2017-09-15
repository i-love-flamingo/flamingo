package web

import (
	"net/url"
	"strings"
)

// UrlTitle normalizes a title for nice usage in URLs
func UrlTitle(title string) string {
	newTitle := strings.ToLower(strings.Replace(title, " ", "-", -1))
	return url.QueryEscape(newTitle)
}
