package web

import (
	"strings"
)

// URLTitle normalizes a title for nice usage in URLs
func URLTitle(title string) string {
	url := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(title, "/", "_"), " ", "-"))
	url = strings.ReplaceAll(url, "-_", "-")
	url = strings.ReplaceAll(url, "%", "-")

	for strings.Contains(url, "--") {
		url = strings.ReplaceAll(url, "--", "-")
	}
	return url
}
