package util

import (
	"os"
	"path/filepath"
	"strings"
)

func getAbsUploadPath(s string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	path := filepath.Dir(ex)

	return filepath.Join(path, s)
}

func getAbsUploadUrl(s string) string {
	if isUrl(s) {
		return strings.TrimSuffix(s, "/")
	}

	rp := strings.TrimSuffix(strings.TrimSuffix(s, "*"), "/")

	return rp
}
