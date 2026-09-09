package folder

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"serverigh/internal/filesystem/fileinfo"
)

const directoryReadEntries = 256

func (s Service) Children(parentPath string) ([]Entry, bool, error) {
	parent, err := s.resolve(parentPath)
	if err != nil {
		return nil, false, err
	}
	if !parent.Info.IsDir() {
		return nil, false, ErrNotDirectory
	}

	directory, err := os.Open(parent.Absolute)
	if err != nil {
		return nil, false, wrapFileError("open folder children", err)
	}
	defer directory.Close()

	children, truncated, err := s.children(parent, directory)
	if err != nil {
		return nil, false, err
	}
	sort.SliceStable(children, func(left int, right int) bool {
		return strings.ToLower(children[left].Name) <
			strings.ToLower(children[right].Name)
	})

	return children, truncated, nil
}

func (s Service) children(
	parent resolvedPath,
	directory *os.File,
) ([]Entry, bool, error) {
	depth := depth(parent.Path) + 1
	children := make([]Entry, 0, directoryReadEntries)
	for {
		entries, err := directory.ReadDir(directoryReadEntries)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, false, wrapFileError("read folder children", err)
		}

		for _, entry := range entries {
			child, ok := s.child(parent, entry, depth)
			if ok {
				children = append(children, child)
			}
			if len(children) >= maxTreeEntries {
				return children, true, nil
			}
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return children, false, nil
}

func (s Service) child(
	parent resolvedPath,
	entry os.DirEntry,
	depth int,
) (Entry, bool) {
	childPath := filepath.Join(parent.Absolute, entry.Name())
	if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
		return Entry{}, false
	}
	if fileinfo.IsHidden(childPath, entry.Name()) && !s.showHidden {
		return Entry{}, false
	}

	return Entry{
		Name:        entry.Name(),
		Path:        displayPath(parent.Path, entry.Name()),
		Depth:       depth,
		HasChildren: true,
	}, true
}
