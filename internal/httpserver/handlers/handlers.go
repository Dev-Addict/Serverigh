package handlers

import (
	"html/template"

	"serverigh/internal/apperror"
	"serverigh/internal/config"
	"serverigh/internal/filesystem"
	"serverigh/web"
)

type Handlers struct {
	config    ConfigView
	files     filesystem.Service
	templates *template.Template
}

type ConfigView struct {
	Root            string
	Write           bool
	ShowHidden      bool
	MaxPreviewBytes int64
}

func New(cfg config.Config) (Handlers, error) {
	files, err := filesystem.New(
		cfg.Root,
		cfg.ShowHidden,
		cfg.MaxPreviewBytes,
	)
	if err != nil {
		return Handlers{}, err
	}

	templates, err := template.ParseFS(
		web.Templates,
		"templates/*.html",
		"templates/partials/*.html",
	)
	if err != nil {
		return Handlers{}, apperror.WrapOperation(
			apperror.CodeServer,
			"server error",
			"parse templates",
			err,
		)
	}

	return Handlers{
		templates: templates,
		files:     files,
		config: ConfigView{
			Root:            cfg.Root,
			Write:           cfg.Write,
			ShowHidden:      cfg.ShowHidden,
			MaxPreviewBytes: cfg.MaxPreviewBytes,
		},
	}, nil
}
