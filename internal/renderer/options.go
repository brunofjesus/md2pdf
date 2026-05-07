package renderer

import (
	"log"

	"github.com/brunofjesus/md2pdf/internal/renderer/node"
	"github.com/gomarkdown/markdown/parser"
)

// WithHorizontalRuleAsNewPage configures the renderer to interpret horizontal rules (---) as page breaks,
// starting a new page whenever an --- is encountered in the markdown input.
func WithHorizontalRuleAsNewPage() RenderOption {
	return func(r *PdfRenderer) {
		err := r.SetNodeProcessor("HorizontalRule", node.HorizontalRulePageBreakProcessor)
		if err != nil {
			log.Fatalf("failed to set node processor for HorizontalRule: %v", err)
		}
	}
}

// WithBaseURL sets the base URL for resolving relative links and images in the
// input markdown.
// This is useful when the input is a URL, so that relative links and images
// can be resolved correctly.
func WithBaseURL(baseURL string) RenderOption {
	return func(r *PdfRenderer) {
		r.InputBaseURL = baseURL
	}
}

func WithMarkdownParsingExtensions(exts parser.Extensions) RenderOption {
	return func(r *PdfRenderer) {
		r.Extensions = r.Extensions | exts
	}
}

func WithDefaultMarkdownParsingExtensions() RenderOption {
	return func(r *PdfRenderer) {
		r.Extensions = r.Extensions | parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
			parser.Autolink | parser.Strikethrough | parser.SpaceHeadings |
			parser.HeadingIDs | parser.BackslashLineBreak | parser.DefinitionLists
	}
}
