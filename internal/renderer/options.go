package renderer

import (
	"fmt"
	"log"

	"github.com/brunofjesus/md2pdf/v3/internal/renderer/node"
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

func WithDefaultFooter(orientation, author, title string) RenderOption {
	return func(r *PdfRenderer) {
		r.Pdf.SetFooterFunc(func() {
			r.Pdf.SetFillColor(r.Theme.BackgroundColor.Red, r.Theme.BackgroundColor.Green, r.Theme.BackgroundColor.Blue)
			// Position at 1.5 cm from bottom
			r.Pdf.SetY(-15)
			// Arial italic 8
			r.Pdf.SetFont(r.Theme.Normal.Font, "I", 8)
			// Text color in gray
			r.Pdf.SetTextColor(128, 128, 128)
			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())
			r.Pdf.SetX(4)
			r.Pdf.CellFormat(0, 10, author, "", 0, "", true, 0, "")
			middle := w / 2
			if orientation == "landscape" {
				middle = h / 2
			}
			r.Pdf.SetX(middle - float64(len(title)))
			r.Pdf.CellFormat(0, 10, title, "", 0, "", true, 0, "")
			r.Pdf.SetX(-40)
			r.Pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", r.Pdf.PageNo()), "", 0, "", true, 0, "")
		})
	}
}
