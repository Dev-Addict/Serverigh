package filesystem

import "serverigh/internal/filesystem/searcher"

type SearchOptions = searcher.Options
type SearchResults = searcher.Results
type SearchResult = searcher.Result

func (s Service) Search(options SearchOptions) (SearchResults, error) {
	searchService := searcher.NewService(
		s.root,
		s.showHidden,
		s.searchIndex,
	)

	return searchService.Search(options)
}
