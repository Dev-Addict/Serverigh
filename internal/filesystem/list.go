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
	IsTrash    bool
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
	CreatedValue string
	CreatedKnown bool
	Mode         string
	ModTime      time.Time
	ModTimeLabel string
	ModTimeValue string
	IsDir        bool
	IsHidden     bool
	IsTrashRoot  bool
	IsTrashEntry bool
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
	if isTrashPath(requestPath) {
		return s.listTrash(entryLimit, options)
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
			if resolved.Path == "/" && name == "trash" {
				continue
			}
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

			listing.Entries = append(
				listing.Entries,
				entryFromInfo(name, displayPath(resolved.Path, name), entryPath, info),
			)
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return sortedListing(listing), nil
}

func entryFromInfo(
	name string,
	requestPath string,
	absolutePath string,
	info os.FileInfo,
) Entry {
	createdTime, createdKnown := fileinfo.CreationTime(absolutePath)

	return Entry{
		Name:         name,
		Path:         requestPath,
		RelativePath: strings.TrimPrefix(requestPath, "/"),
		AbsolutePath: absolutePath,
		Kind:         fileinfo.Kind(info),
		Size:         info.Size(),
		SizeLabel:    fileinfo.FormatSize(info.Size()),
		CreatedTime:  createdTime,
		CreatedLabel: fileinfo.FormatOptionalTime(createdTime, createdKnown),
		CreatedValue: fileinfo.FormatOptionalTimeValue(
			createdTime,
			createdKnown,
		),
		CreatedKnown: createdKnown,
		Mode:         info.Mode().String(),
		ModTime:      info.ModTime(),
		ModTimeLabel: fileinfo.FormatTime(info.ModTime()),
		ModTimeValue: fileinfo.FormatTimeValue(info.ModTime()),
		IsDir:        info.IsDir(),
		IsHidden:     fileinfo.IsHidden(absolutePath, name),
	}
}
