package filesystem

import (
	"path/filepath"
	"testing"
)

func TestSearchPrioritizesCurrentDirectory(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs"))
	mkdir(t, filepath.Join(root, "archive"))
	writeFile(t, filepath.Join(root, "docs", "report.txt"), "current")
	writeFile(t, filepath.Join(root, "archive", "report.txt"), "archive")

	files := newTestService(t, root, false, 1024)
	results, err := files.Search(SearchOptions{
		Query:      "report",
		ActivePath: "/docs",
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("search files: %v", err)
	}

	if len(results.Results) < 2 {
		t.Fatalf("expected search results, got %#v", results)
	}

	if results.Results[0].Path != "/docs/report.txt" {
		t.Fatalf("expected current directory result first, got %#v", results.Results)
	}
}

func TestSearchPrioritizesActiveSubtree(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "docs", "api"))
	writeFile(t, filepath.Join(root, "docs", "api", "config.yml"), "current")
	writeFile(t, filepath.Join(root, "config.yml"), "root")

	files := newTestService(t, root, false, 1024)
	results, err := files.Search(SearchOptions{
		Query:      "config",
		ActivePath: "/docs",
		Limit:      5,
	})
	if err != nil {
		t.Fatalf("search files: %v", err)
	}

	if len(results.Results) < 2 {
		t.Fatalf("expected search results, got %#v", results)
	}

	if results.Results[0].Path != "/docs/api/config.yml" {
		t.Fatalf("expected active subtree result first, got %#v", results.Results)
	}
}

func TestSearchUsesFuzzyMatching(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "server_config.yaml"), "config")

	files := newTestService(t, root, false, 1024)
	results, err := files.Search(SearchOptions{
		Query:      "svcfg",
		ActivePath: "/",
	})
	if err != nil {
		t.Fatalf("search files: %v", err)
	}

	if len(results.Results) != 1 {
		t.Fatalf("expected fuzzy result, got %#v", results)
	}
}

func TestSearchLimitsReturnedResults(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "one-report.txt"), "one")
	writeFile(t, filepath.Join(root, "two-report.txt"), "two")

	files := newTestService(t, root, false, 1024)
	results, err := files.Search(SearchOptions{
		Query:      "report",
		ActivePath: "/",
		Limit:      1,
	})
	if err != nil {
		t.Fatalf("search files: %v", err)
	}

	if len(results.Results) != 1 {
		t.Fatalf("expected one returned result, got %#v", results)
	}

	if !results.Truncated {
		t.Fatalf("expected truncated result set, got %#v", results)
	}
}

func TestSearchSkipsHiddenFilesUnlessConfigured(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".env"), "secret")

	hiddenDisabled := newTestService(t, root, false, 1024)
	hiddenResults, err := hiddenDisabled.Search(SearchOptions{
		Query:      "env",
		ActivePath: "/",
	})
	if err != nil {
		t.Fatalf("search files: %v", err)
	}

	if len(hiddenResults.Results) != 0 {
		t.Fatalf("expected hidden file to be skipped, got %#v", hiddenResults)
	}

	hiddenEnabled := newTestService(t, root, true, 1024)
	visibleResults, err := hiddenEnabled.Search(SearchOptions{
		Query:      "env",
		ActivePath: "/",
	})
	if err != nil {
		t.Fatalf("search files with hidden enabled: %v", err)
	}

	if len(visibleResults.Results) != 1 || visibleResults.Results[0].Path != "/.env" {
		t.Fatalf("expected hidden result when configured, got %#v", visibleResults)
	}
}
