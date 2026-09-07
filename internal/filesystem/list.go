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
	Options    ListOptions
	Truncated  bool
	EntryLimit int
}

type Entry struct {
	Name         string
	Path         string
	RelativePath string
	AbsolutePath string
	Kind         string
	Size         int64
	SizeLabel    string
	CreatedTime  time.Time
	CreatedLabel string
	CreatedKnown bool
	Mode         string
	ModTime      time.Time
	ModTimeLabel string
	IsDir        bool
	IsHidden     bool
}

func (s Service) List(requestPath string) (DirectoryListing, error) {
	return s.ListWithOptions(requestPath, DefaultListOptions())
}

func (s Service) ListWithOptions(
	requestPath string,
	options ListOptions,
) (DirectoryListing, error) {
	return s.list(requestPath, maxDirectoryEntries, NormalizeListOptions(options))
}

func (s Service) list(
	requestPath string,
	entryLimit int,
	options ListOptions,
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
		Options:    options,
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
			createdTime, createdKnown := creationTime(entryPath)
			listing.Entries = append(listing.Entries, Entry{
				Name:         name,
				Path:         path,
				RelativePath: strings.TrimPrefix(path, "/"),
				AbsolutePath: entryPath,
				Kind:         entryKind(info),
				Size:         info.Size(),
				SizeLabel:    formatSize(info.Size()),
				CreatedTime:  createdTime,
				CreatedLabel: formatOptionalTime(createdTime, createdKnown),
				CreatedKnown: createdKnown,
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

		return compareEntries(left, right, listing.Options) < 0
	})

	return listing
}

func compareEntries(left Entry, right Entry, options ListOptions) int {
	result := 0
	switch options.Sort {
	case ListSortSize:
		result = compareInt64(left.Size, right.Size)
	case ListSortModified:
		result = compareTimes(left.ModTime, right.ModTime)
	case ListSortCreated:
		result = compareOptionalTimes(
			left.CreatedTime,
			left.CreatedKnown,
			right.CreatedTime,
			right.CreatedKnown,
		)
	default:
		result = compareNames(left.Name, right.Name)
	}

	if result == 0 {
		result = compareNames(left.Name, right.Name)
	}

	if options.Direction == ListDirectionDesc {
		return -result
	}

	return result
}

func compareNames(left string, right string) int {
	return strings.Compare(strings.ToLower(left), strings.ToLower(right))
}

func compareInt64(left int64, right int64) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func compareTimes(left time.Time, right time.Time) int {
	switch {
	case left.Before(right):
		return -1
	case left.After(right):
		return 1
	default:
		return 0
	}
}

func compareOptionalTimes(
	left time.Time,
	leftKnown bool,
	right time.Time,
	rightKnown bool,
) int {
	if leftKnown != rightKnown {
		if leftKnown {
			return -1
		}

		return 1
	}

	return compareTimes(left, right)
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
