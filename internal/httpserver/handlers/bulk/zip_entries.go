package bulk

import (
	"os"
	"path"
	"path/filepath"

	"serverigh/internal/filesystem"
)

type bulkZipEntry struct {
	Path string
	Size int64
}

func collectBulkZipEntries(
	targets []filesystem.ResolvedPath,
) ([]bulkZipEntry, error) {
	entries := make([]bulkZipEntry, 0, len(targets))
	for _, target := range targets {
		nextEntries, err := zipEntriesForTarget(target)
		if err != nil {
			return nil, err
		}
		entries = append(entries, nextEntries...)
		if err := validateBulkZipLimits(entries); err != nil {
			return nil, err
		}
	}

	return entries, nil
}

func zipEntriesForTarget(
	target filesystem.ResolvedPath,
) ([]bulkZipEntry, error) {
	if !target.Info.IsDir() {
		return []bulkZipEntry{{
			Path: target.Path,
			Size: target.Info.Size(),
		}}, nil
	}

	return zipEntriesForDirectory(target)
}

func zipEntriesForDirectory(
	target filesystem.ResolvedPath,
) ([]bulkZipEntry, error) {
	entries := []bulkZipEntry{}
	err := filepath.WalkDir(
		target.Absolute,
		func(filename string, entry os.DirEntry, err error) error {
			return appendWalkedZipEntry(&entries, target, filename, entry, err)
		},
	)

	return entries, err
}

func appendWalkedZipEntry(
	entries *[]bulkZipEntry,
	target filesystem.ResolvedPath,
	filename string,
	entry os.DirEntry,
	err error,
) error {
	if err != nil || entry.IsDir() {
		return err
	}

	info, err := entry.Info()
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(target.Absolute, filename)
	if err != nil {
		return err
	}

	*entries = append(*entries, bulkZipEntry{
		Path: path.Join(target.Path, filepath.ToSlash(relative)),
		Size: info.Size(),
	})

	return validateBulkZipLimits(*entries)
}
