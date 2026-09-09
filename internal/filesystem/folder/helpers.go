package folder

import (
	"path"
	"strings"
)

func (tree *Tree) append(entry Entry) {
	if len(tree.Entries) >= maxTreeEntries {
		tree.Truncated = true

		return
	}

	tree.Entries = append(tree.Entries, entry)
}

func ancestors(requestPath string) []string {
	cleanPath := path.Clean("/" + strings.TrimPrefix(requestPath, "/"))
	if cleanPath == "/" {
		return []string{"/"}
	}

	parts := strings.Split(strings.TrimPrefix(cleanPath, "/"), "/")
	ancestors := []string{"/"}
	for index := range parts {
		ancestors = append(ancestors, "/"+path.Join(parts[:index+1]...))
	}

	return ancestors
}

func depth(requestPath string) int {
	if requestPath == "/" {
		return 0
	}

	return strings.Count(strings.Trim(requestPath, "/"), "/") + 1
}
