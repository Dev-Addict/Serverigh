package view

import "serverigh/internal/filesystem"

type BrowsePage struct {
	Root        string
	Path        string
	Mode        string
	Listing     FilesView
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	Search      SearchBoxView
	EmptyState  EmptyPreviewView
	Status      StatusView
	Preview     *filesystem.Preview
}
