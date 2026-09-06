package web

import "embed"

//go:embed static/*
var Static embed.FS

//go:embed templates/*.html templates/partials/*.html
var Templates embed.FS
