package util

import (
	"net/url"
	"strings"
)

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func urlDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

func isUrl(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
