package view

import (
	"serverigh/internal/config"
	"serverigh/internal/filesystem"
)

type BrowsePage struct {
	Root            string
	Path            string
	Mode            string
	Theme           string
	MaxPreviewBytes int64
	Listing         FilesView
	PathSummary     PathSummaryView
	Breadcrumbs     BreadcrumbsView
	Search          SearchBoxView
	EmptyState      EmptyPreviewView
	Status          StatusView
	Columns         config.Columns
	Preview         *filesystem.Preview
}
