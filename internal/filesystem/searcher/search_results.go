package searcher

import (
	"path"
	"sort"
	"strings"

	"github.com/sahilm/fuzzy"
)

func searchResultsFromMatches(
	entries []searchEntry,
	matches fuzzy.Matches,
	activePath string,
) []Result {
	results := make([]Result, 0, len(matches))
	for _, match := range matches {
		entry := entries[match.Index]
		score := match.Score + scopeScore(entry.ParentPath, activePath)
		if !entry.IsDir {
			score += 30
		}

		results = append(results, Result{
			Name:       entry.Name,
			Path:       entry.Path,
			ParentPath: entry.ParentPath,
			Kind:       entry.Kind,
			IsDir:      entry.IsDir,
			Score:      score,
		})
	}

	return results
}

func sortedResults(results []Result) []Result {
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
