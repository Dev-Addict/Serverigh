package filesystem

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"
	"github.com/sahilm/fuzzy"
)

const (
	defaultSearchLimit = 12
)

type SearchOptions struct {
	Query      string
	ActivePath string
	Limit      int
}

type SearchResults struct {
	Query     string
	Results   []SearchResult
	Truncated bool
}

type SearchResult struct {
	Name         string
	Path         string
	ParentPath   string
	Kind         string
	Size         int64
	SizeLabel    string
	Mode         string
	ModTimeLabel string
	CreatedLabel string
	CreatedKnown bool
	IsDir        bool
	Score        int
}

func (s Service) Search(options SearchOptions) (SearchResults, error) {
	query := strings.TrimSpace(options.Query)
	activePath, err := cleanRequestPath(options.ActivePath)
	if err != nil {
		return SearchResults{}, err
	}

	if options.Limit < 1 {
		options.Limit = defaultSearchLimit
	}

	results := SearchResults{
		Query: query,
	}
	if query == "" {
		return results, nil
	}

	entries, err := s.searchEntries()
	if err != nil {
		return SearchResults{}, err
	}

	matches := fuzzy.FindFromNoSort(query, searchEntrySource(entries))
	results.Truncated = len(matches) > options.Limit
	results.Results = topSearchResults(
		searchResultsFromMatches(entries, matches, activePath),
		options.Limit,
	)

	return results, nil
}

func (s Service) searchEntries() ([]searchEntry, error) {
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
			if err != nil {
				return wrapFileError("search files", err)
			}

			if filename == s.root {
				return nil
			}

			name := entry.Name()
			hidden := isHidden(filename, name)
			if hidden && !s.showHidden {
				if entry.IsDir() {
					return filepath.SkipDir
				}

				return nil
			}

			searchEntry, ok, err := searchEntryForDirEntry(
				filename,
				entry,
				s.root,
			)
			if err != nil {
				return wrapFileError("inspect search entry", err)
			}

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

type searchEntry struct {
	Name         string
	Path         string
	ParentPath   string
	Kind         string
	Size         int64
	SizeLabel    string
	Mode         string
	ModTimeLabel string
	CreatedLabel string
	CreatedKnown bool
	IsDir        bool
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
) (searchEntry, bool, error) {
	relativePath, err := filepath.Rel(root, filename)
	if err != nil {
		return searchEntry{}, false, nil
	}

	info, err := fastwalk.StatDirEntry(filename, entry)
	if err != nil {
		return searchEntry{}, false, err
	}

	logicalPath := "/" + filepath.ToSlash(relativePath)
	parentPath := path.Dir(logicalPath)
	if parentPath == "." {
		parentPath = "/"
	}

	createdTime, createdKnown := creationTime(filename)

	return searchEntry{
		Name:         entry.Name(),
		Path:         logicalPath,
		ParentPath:   parentPath,
		Kind:         entryKind(info),
		Size:         info.Size(),
		SizeLabel:    formatSize(info.Size()),
		Mode:         info.Mode().String(),
		ModTimeLabel: info.ModTime().Format("2006-01-02 15:04"),
		CreatedLabel: formatOptionalTime(createdTime, createdKnown),
		CreatedKnown: createdKnown,
		IsDir:        info.IsDir(),
	}, true, nil
}

func searchResultsFromMatches(
	entries []searchEntry,
	matches fuzzy.Matches,
	activePath string,
) []SearchResult {
	results := make([]SearchResult, 0, len(matches))
	for _, match := range matches {
		entry := entries[match.Index]
		score := match.Score + scopeScore(entry.ParentPath, activePath)
		if !entry.IsDir {
			score += 30
		}

		results = append(results, SearchResult{
			Name:         entry.Name,
			Path:         entry.Path,
			ParentPath:   entry.ParentPath,
			Kind:         entry.Kind,
			Size:         entry.Size,
			SizeLabel:    entry.SizeLabel,
			Mode:         entry.Mode,
			ModTimeLabel: entry.ModTimeLabel,
			CreatedLabel: entry.CreatedLabel,
			CreatedKnown: entry.CreatedKnown,
			IsDir:        entry.IsDir,
			Score:        score,
		})
	}

	return results
}

func scopeScore(parentPath string, activePath string) int {
	activePath = cleanSearchScope(activePath)
	parentPath = cleanSearchScope(parentPath)
	if parentPath == activePath {
		return 500
	}

	if activePath != "/" && strings.HasPrefix(parentPath, activePath+"/") {
		return 300 - pathDepth(strings.TrimPrefix(parentPath, activePath))
	}

	if activePath == "/" {
		return 250 - pathDepth(parentPath)
	}

	return 0
}

func topSearchResults(results []SearchResult, limit int) []SearchResult {
	sort.SliceStable(results, func(i int, j int) bool {
		left := results[i]
		right := results[j]
		if left.Score != right.Score {
			return left.Score > right.Score
		}

		if left.IsDir != right.IsDir {
			return !left.IsDir
		}

		return compareNames(left.Path, right.Path) < 0
	})

	if len(results) > limit {
		return results[:limit]
	}

	return results
}

func cleanSearchScope(activePath string) string {
	cleanPath := path.Clean("/" + strings.TrimPrefix(activePath, "/"))
	if cleanPath == "." {
		return "/"
	}

	return cleanPath
}

func pathDepth(activePath string) int {
	activePath = strings.Trim(activePath, "/")
	if activePath == "" {
		return 0
	}

	return strings.Count(activePath, "/") + 1
}
