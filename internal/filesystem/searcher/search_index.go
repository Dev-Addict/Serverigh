package searcher

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"github.com/charlievieth/fastwalk"

	"serverigh/internal/filesystem/fileinfo"
)

const IndexCacheTTL = 2 * time.Second

type Index struct {
	mutex   sync.Mutex
	entries []searchEntry
	expires time.Time
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) Clear() {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.entries = nil
	i.expires = time.Time{}
}

func (s Service) cachedSearchEntries(ctx context.Context) ([]searchEntry, error) {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return nil, ctxErr
	}

	now := time.Now()
	s.Index.mutex.Lock()
	if now.Before(s.Index.expires) {
		entries := append([]searchEntry(nil), s.Index.entries...)
		s.Index.mutex.Unlock()

		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}

		return entries, nil
	}
	s.Index.mutex.Unlock()

	entries, err := s.searchEntries(ctx)
	if err != nil {
		return nil, err
	}

	s.Index.mutex.Lock()
	s.Index.entries = append(s.Index.entries[:0], entries...)
	s.Index.expires = now.Add(IndexCacheTTL)
	cachedEntries := append([]searchEntry(nil), s.Index.entries...)
	s.Index.mutex.Unlock()

	return cachedEntries, nil
}

func (s Service) searchEntries(ctx context.Context) ([]searchEntry, error) {
	var (
		entries []searchEntry
		mutex   sync.Mutex
	)

	walkErr := fastwalk.Walk(
		&fastwalk.Config{
			Follow: false,
			Sort:   fastwalk.SortLexical,
		},
		s.root,
		func(filename string, entry fs.DirEntry, err error) error {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}

			if err != nil {
				if errors.Is(err, fs.ErrPermission) {
					return nil
				}

				return wrapFileError("search files", err)
			}

			if filename == s.root {
				return nil
			}

			name := entry.Name()
			hidden := fileinfo.IsHidden(filename, name)
			if hidden && !s.showHidden {
				if entry.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}

			searchEntry, ok := searchEntryForDirEntry(
				filename,
				entry,
				s.root,
			)

			if ok {
				mutex.Lock()
				entries = append(entries, searchEntry)
				mutex.Unlock()
			}

			return nil
		},
	)
	if walkErr != nil {
		return nil, walkErr
	}

	return entries, nil
}
