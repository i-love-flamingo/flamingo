package web

import (
	"fmt"
	"strings"
)

// URLTitle normalizes a title for nice usage in URLs
func URLTitle(title string) string {
	url := strings.ToLower(strings.Replace(strings.Replace(title, "/", "_", -1), " ", "-", -1))
	url = strings.Replace(url, "-_", "-", -1)
	url = strings.Replace(url, "%", "-", -1)
	for strings.Contains(url, "--") {
		url = strings.Replace(url, "--", "-", -1)
	}

	return url
}

// fmtErrorf shim to allow go 1.12 backwards compatibility
var fmtErrorf = fmt.Errorf
