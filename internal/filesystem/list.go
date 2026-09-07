package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"serverigh/internal/filesystem/fileinfo"
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
			hidden := fileinfo.IsHidden(entryPath, name)
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
			createdTime, createdKnown := fileinfo.CreationTime(entryPath)
			listing.Entries = append(listing.Entries, Entry{
				Name:         name,
				Path:         path,
				RelativePath: strings.TrimPrefix(path, "/"),
				AbsolutePath: entryPath,
				Kind:         fileinfo.Kind(info),
				Size:         info.Size(),
				SizeLabel:    fileinfo.FormatSize(info.Size()),
				CreatedTime:  createdTime,
				CreatedLabel: fileinfo.FormatOptionalTime(createdTime, createdKnown),
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
