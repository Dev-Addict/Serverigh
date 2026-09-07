package searcher

import (
	"context"
	"strings"

	"github.com/sahilm/fuzzy"
)

const (
	defaultSearchLimit   = 12
	minSearchQueryLength = 2
)

type Options struct {
	Context    context.Context
	Query      string
	ActivePath string
	Limit      int
}

type Results struct {
	Query     string
	Results   []Result
	Truncated bool
}

type Result struct {
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

func (s Service) Search(options Options) (Results, error) {
	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}

	query := strings.TrimSpace(options.Query)
	activePath, err := cleanRequestPath(options.ActivePath)
	if err != nil {
		return Results{}, err
	}

	if options.Limit < 1 {
		options.Limit = defaultSearchLimit
	}

	results := Results{
		Query: query,
	}
	if query == "" {
		return results, nil
	}

	if len([]rune(query)) < minSearchQueryLength {
		return results, nil
	}

	entries, err := s.cachedSearchEntries(ctx)
	if err != nil {
		return Results{}, err
	}

	matches := fuzzy.FindFromNoSort(query, searchEntrySource(entries))
	scoredResults := sortedResults(
		searchResultsFromMatches(entries, matches, activePath),
	)
	results.Truncated = len(scoredResults) > options.Limit
	results.Results, err = s.enrichResults(
		ctx,
		scoredResults,
		options.Limit,
	)
	if err != nil {
		return Results{}, err
	}

	return results, nil
}
