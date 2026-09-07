package view

type FilesUpdate struct {
	Listing     FilesView
	PathSummary PathSummaryView
	Breadcrumbs BreadcrumbsView
	Search      SearchBoxView
	EmptyState  EmptyPreviewView
	Status      StatusView
}
