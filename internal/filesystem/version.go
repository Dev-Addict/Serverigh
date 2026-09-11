package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"serverigh/internal/filesystem/fileinfo"
)

type DirectoryVersion struct {
	Path    string
	Version string
}

func (s Service) DirectoryVersion(
	requestPath string,
	options ListOptions,
) (DirectoryVersion, error) {
	options = NormalizeListOptions(options)
	if isTrashPath(requestPath) {
		return s.directoryVersionForAbsolute(
			TrashPath,
			trashStoragePath(s.root),
			options,
			true,
		)
	}

	resolved, err := s.Resolve(requestPath)
	if err != nil {
		return DirectoryVersion{}, err
	}
	if !resolved.Info.IsDir() {
		return DirectoryVersion{}, ErrNotDirectory
	}

	return s.directoryVersionForAbsolute(
		resolved.Path,
		resolved.Absolute,
		options,
		false,
	)
}

func (s Service) directoryVersionForAbsolute(
	logicalPath string,
	absolutePath string,
	options ListOptions,
	isTrash bool,
) (DirectoryVersion, error) {
	listing := DirectoryListing{
		Path:       logicalPath,
		Root:       s.root,
		Options:    options,
		EntryLimit: maxDirectoryEntries,
		IsTrash:    isTrash,
	}

	directory, err := os.Open(absolutePath)
	if errors.Is(err, os.ErrNotExist) && isTrash {
		return DirectoryVersion{
			Path:    listing.Path,
			Version: directoryListingVersion(listing),
		}, nil
	}
	if err != nil {
		return DirectoryVersion{}, wrapFileError("open directory", err)
	}
	defer directory.Close()

	for {
		entries, err := directory.ReadDir(directoryReadEntries)
		if err != nil && !errors.Is(err, io.EOF) {
			return DirectoryVersion{}, wrapFileError("read directory", err)
		}

		for _, entry := range entries {
			entryPath := filepath.Join(absolutePath, entry.Name())
			if s.skipVersionEntry(logicalPath, entryPath, entry.Name(), isTrash) {
				continue
			}
			if len(listing.Entries) >= maxDirectoryEntries {
				listing.Truncated = true
				listing = sortedListing(listing)

				return DirectoryVersion{
					Path:    listing.Path,
					Version: directoryListingVersion(listing),
				}, nil
			}

			info, err := entry.Info()
			if err != nil {
				return DirectoryVersion{}, wrapFileError(
					"inspect directory entry",
					err,
				)
			}

			name := entry.Name()
			displayEntryPath := displayPath(logicalPath, name)
			if isTrash {
				name = s.trashEntryName(name)
			}
			listing.Entries = append(listing.Entries, Entry{
				Name:         name,
				Path:         displayEntryPath,
				Kind:         fileinfo.Kind(info),
				Size:         info.Size(),
				Mode:         info.Mode().String(),
				ModTime:      info.ModTime(),
				IsDir:        info.IsDir(),
				IsHidden:     fileinfo.IsHidden(entryPath, entry.Name()),
				IsTrashEntry: isTrash,
			})
		}

		if errors.Is(err, io.EOF) {
			break
		}
	}

	listing = sortedListing(listing)

	return DirectoryVersion{
		Path:    listing.Path,
		Version: directoryListingVersion(listing),
	}, nil
}

func (s Service) skipVersionEntry(
	logicalPath string,
	absolutePath string,
	name string,
	isTrash bool,
) bool {
	if isTrash {
		return name == trashMetaDirectoryName
	}
	if logicalPath == "/" && name == "trash" {
		return true
	}

	return !s.showHidden && fileinfo.IsHidden(absolutePath, name)
}

func directoryListingVersion(listing DirectoryListing) string {
	hash := sha256.New()
	writeVersionValue(hash, listing.Path)
	writeVersionValue(hash, string(listing.Options.Sort))
	writeVersionValue(hash, string(listing.Options.Direction))
	writeVersionValue(hash, strconv.FormatBool(listing.Truncated))
	writeVersionValue(hash, strconv.Itoa(listing.EntryLimit))
	writeVersionValue(hash, strconv.FormatBool(listing.IsTrash))

	for _, entry := range listing.Entries {
		writeVersionValue(hash, entry.Name)
		writeVersionValue(hash, entry.Path)
		writeVersionValue(hash, entry.Kind)
		writeVersionValue(hash, strconv.FormatInt(entry.Size, 10))
		writeVersionValue(hash, entry.Mode)
		writeVersionValue(hash, timeVersion(entry.ModTime))
		writeVersionValue(hash, strconv.FormatBool(entry.IsDir))
		writeVersionValue(hash, strconv.FormatBool(entry.IsHidden))
		writeVersionValue(hash, strconv.FormatBool(entry.IsTrashEntry))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func writeVersionValue(hash interface{ Write([]byte) (int, error) }, value string) {
	_, _ = hash.Write([]byte(value))
	_, _ = hash.Write([]byte{0})
}

func timeVersion(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return strconv.FormatInt(value.UnixNano(), 10)
}
