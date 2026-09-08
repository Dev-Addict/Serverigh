package previewer

import (
	"bytes"
	"html/template"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	lightSyntaxHighlightStyle = "xcode"
	darkSyntaxHighlightStyle  = "monokai"
)

func highlightCode(
	name string,
	mimeType string,
	content string,
) (template.HTML, template.HTML, string, bool) {
	lexer := lexers.Match(name)
	if lexer == nil {
		lexer = lexers.MatchMimeType(mimeType)
	}

	if !isHighlightLexer(lexer) {
		return "", "", "", false
	}

	lexer = chroma.Coalesce(lexer)
	lightHTML, ok := renderHighlightedCode(lexer, lightSyntaxHighlightStyle, content)
	if !ok {
		return "", "", "", false
	}

	darkHTML, ok := renderHighlightedCode(lexer, darkSyntaxHighlightStyle, content)
	if !ok {
		return "", "", "", false
	}

	return lightHTML, darkHTML, lexer.Config().Name, true
}

func renderHighlightedCode(
	lexer chroma.Lexer,
	styleName string,
	content string,
) (template.HTML, bool) {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(
		html.TabWidth(4),
		html.WrapLongLines(false),
	)

	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return "", false
	}

	var output bytes.Buffer
	if err := formatter.Format(&output, style, iterator); err != nil {
		return "", false
	}

	return template.HTML(output.String()), true
}

func isHighlightLexer(lexer chroma.Lexer) bool {
	if lexer == nil {
		return false
	}

	switch lexer.Config().Name {
	case "fallback", "plaintext":
		return false
	default:
		return true
	}
}
