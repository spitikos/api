package docs

import (
	"path/filepath"
	"strings"
)

func trimExtension(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}

func slugToKey(slug string) string {
	return "docs:" + strings.ReplaceAll(slug, "/", ":")
}

func keyToSlug(key string) string {
	return strings.ReplaceAll(strings.TrimPrefix(key, "docs:"), ":", "/")
}

func getFirstH1(content string) string {
	for line := range strings.SplitSeq(content, "\n") {
		h1, found := strings.CutPrefix(line, "# ")
		if found {
			return h1
		}
	}
	return "Untitled"
}
