package handlers

import (
	"html/template"

	"serverigh/internal/config"
	"serverigh/web"
)

type Handlers struct {
	config    ConfigView
	templates *template.Template
}

type ConfigView struct {
	Root            string
	Write           bool
	ShowHidden      bool
	MaxPreviewBytes int64
}

func New(cfg config.Config) Handlers {
	return Handlers{
		templates: template.Must(
			template.ParseFS(web.Templates, "templates/*.html"),
		),
		config: ConfigView{
			Root:            cfg.Root,
			Write:           cfg.Write,
			ShowHidden:      cfg.ShowHidden,
			MaxPreviewBytes: cfg.MaxPreviewBytes,
		},
	}
}
