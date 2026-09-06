package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maxDirectoryEntries  = 2_000
	directoryReadEntries = 256
)

type DirectoryListing struct {
	Path       string
	Root       string
	Entries    []Entry
	Truncated  bool
	EntryLimit int
}

type Entry struct {
	Name         string
	Path         string
	Kind         string
	Size         int64
	SizeLabel    string
	Mode         string
	ModTime      time.Time
	ModTimeLabel string
	IsDir        bool
	IsHidden     bool
}

func (s Service) List(requestPath string) (DirectoryListing, error) {
	return s.list(requestPath, maxDirectoryEntries)
}

func (s Service) list(
	requestPath string,
	entryLimit int,
) (DirectoryListing, error) {
	if entryLimit < 1 {
		return DirectoryListing{}, ErrInvalidPath
	}

	resolved, err := s.Resolve(requestPath)
	if err != nil {
		return DirectoryListing{}, err
	}

	if !resolved.Info.IsDir() {
		return DirectoryListing{}, ErrNotDirectory
	}

	directory, err := os.Open(resolved.Absolute)
	if err != nil {
		return DirectoryListing{}, wrapFileError("open directory", err)
	}
	defer directory.Close()

	listing := DirectoryListing{
		Path:       resolved.Path,
		Root:       s.root,
		EntryLimit: entryLimit,
	}

	for {
		entries, err := directory.ReadDir(directoryReadEntries)
		if err != nil && !errors.Is(err, io.EOF) {
			return DirectoryListing{}, wrapFileError("read directory", err)
		}

		for _, entry := range entries {
			name := entry.Name()
			entryPath := filepath.Join(resolved.Absolute, name)
			hidden := isHidden(entryPath, name)
			if hidden && !s.showHidden {
				continue
			}

			if len(listing.Entries) >= entryLimit {
				listing.Truncated = true

				return sortedListing(listing), nil
			}

			info, err := entry.Info()
			if err != nil {
				return DirectoryListing{}, wrapFileError(
					"inspect directory entry",
					err,
				)
			}

			path := displayPath(resolved.Path, name)
			listing.Entries = append(listing.Entries, Entry{
				Name:         name,
				Path:         path,
				Kind:         entryKind(info),
				Size:         info.Size(),
				SizeLabel:    formatSize(info.Size()),
				Mode:         info.Mode().String(),
				ModTime:      info.ModTime(),
				ModTimeLabel: info.ModTime().Format("2006-01-02 15:04"),
				IsDir:        info.IsDir(),
				IsHidden:     hidden,
			})
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return sortedListing(listing), nil
}

func sortedListing(listing DirectoryListing) DirectoryListing {
	sort.SliceStable(listing.Entries, func(i int, j int) bool {
		left := listing.Entries[i]
		right := listing.Entries[j]
		if left.IsDir != right.IsDir {
			return left.IsDir
		}

		return strings.ToLower(left.Name) < strings.ToLower(right.Name)
	})

	return listing
}

func entryKind(info os.FileInfo) string {
	if info.IsDir() {
		return "folder"
	}

	if info.Mode()&os.ModeSymlink != 0 {
		return "symlink"
	}

	return "file"
}
