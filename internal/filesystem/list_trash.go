package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const trashMetaDirectoryName = "trash-meta"

func (s Service) listTrash(
	entryLimit int,
	options ListOptions,
) (DirectoryListing, error) {
	listing := DirectoryListing{
		Path:       TrashPath,
		Root:       s.root,
		Options:    options,
		EntryLimit: entryLimit,
		IsTrash:    true,
	}

	directory, err := os.Open(trashStoragePath(s.root))
	if errors.Is(err, os.ErrNotExist) {
		return listing, nil
	}
	if err != nil {
		return DirectoryListing{}, wrapFileError("open trash directory", err)
	}
	defer directory.Close()

	for {
		entries, err := directory.ReadDir(directoryReadEntries)
		if err != nil && !errors.Is(err, io.EOF) {
			return DirectoryListing{}, wrapFileError("read trash directory", err)
		}

		for _, entry := range entries {
			if entry.Name() == trashMetaDirectoryName {
				continue
			}
			if len(listing.Entries) >= entryLimit {
				listing.Truncated = true

				return sortedListing(listing), nil
			}

			info, err := entry.Info()
			if err != nil {
				return DirectoryListing{}, wrapFileError(
					"inspect trash entry",
					err,
				)
			}

			filename := filepath.Join(trashStoragePath(s.root), entry.Name())
			entryPath := displayPath(TrashPath, entry.Name())
			viewEntry := entryFromInfo(
				s.trashEntryName(entry.Name()),
				entryPath,
				filename,
				info,
			)
			viewEntry.IsTrashEntry = true
			listing.Entries = append(listing.Entries, viewEntry)
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return sortedListing(listing), nil
}
