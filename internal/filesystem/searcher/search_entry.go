package searcher

import (
	"io/fs"
	"path"
	"path/filepath"
)

type searchEntry struct {
	Name       string
	Path       string
	ParentPath string
	Kind       string
	IsDir      bool
}

type searchEntrySource []searchEntry

func (s searchEntrySource) String(index int) string {
	return s[index].Path
}

func (s searchEntrySource) Len() int {
	return len(s)
}

func searchEntryForDirEntry(
	filename string,
	entry fs.DirEntry,
	root string,
) (searchEntry, bool) {
	relativePath, err := filepath.Rel(root, filename)
	if err != nil {
		return searchEntry{}, false
	}

	logicalPath := "/" + filepath.ToSlash(relativePath)
	parentPath := path.Dir(logicalPath)
	if parentPath == "." {
		parentPath = "/"
	}

	return searchEntry{
		Name:       entry.Name(),
		Path:       logicalPath,
		ParentPath: parentPath,
		Kind:       searchEntryKind(entry),
		IsDir:      entry.IsDir(),
	}, true
}

func searchEntryKind(entry fs.DirEntry) string {
	if entry.IsDir() {
		return "folder"
	}

	if entry.Type()&fs.ModeSymlink != 0 {
		return "symlink"
	}

	return "file"
}
