package previewer

import (
	"html/template"

	"github.com/russross/blackfriday/v2"
)

func renderMarkdown(content []byte) template.HTML {
	flags := blackfriday.SkipHTML |
		blackfriday.Safelink |
		blackfriday.NoreferrerLinks |
		blackfriday.NoopenerLinks

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: flags,
	})

	html := blackfriday.Run(
		content,
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
		blackfriday.WithRenderer(renderer),
	)

	return template.HTML(html)
}
