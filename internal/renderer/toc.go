package renderer

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// tocEntry represents a table of contents entry.
type tocEntry struct {
	Level int
	Title string
	ID    string
}

// tocVisitor implements ast.NodeVisitor to collect headings.
type tocVisitor struct {
	Entries []tocEntry
}

// Visit implements the ast.NodeVisitor interface.
func (v *tocVisitor) Visit(node ast.Node, entering bool) ast.WalkStatus {
	if !entering {
		return ast.GoToNext
	}

	if heading, ok := node.(*ast.Heading); ok {
		title := extractTextFromNode(heading)
		if title != "" {
			id := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(title), " ", "-"))
			id = strings.ReplaceAll(id, ".", "")
			id = strings.ReplaceAll(id, ",", "")
			id = strings.ReplaceAll(id, "!", "")
			id = strings.ReplaceAll(id, "?", "")

			v.Entries = append(v.Entries, tocEntry{
				Level: heading.Level,
				Title: title,
				ID:    id,
			})
		}
	}

	return ast.GoToNext
}

// extractTextFromNode recursively extracts text content from AST nodes.
func extractTextFromNode(node ast.Node) string {
	var text strings.Builder

	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n := node.(type) {
			case *ast.Text:
				text.Write(n.Literal)
			case *ast.Code:
				text.Write(n.Literal)
			}
		}
		return ast.GoToNext
	})

	return text.String()
}

// getTOCEntries parses content and returns all heading entries for the TOC.
func getTOCEntries(content []byte) ([]tocEntry, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := markdown.Parse(content, p)

	visitor := &tocVisitor{}
	ast.Walk(doc, visitor)

	return visitor.Entries, nil
}

// WithTableOfContents registers a pre-processor that generates a table of
// contents page before the main document content is rendered.
func WithTableOfContents() RenderOption {
	return func(r *PdfRenderer) {
		r.preProcessors = append(r.preProcessors, func(content []byte) error {
			entries, err := getTOCEntries(content)
			if err != nil {
				return fmt.Errorf("failed to collect TOC entries: %w", err)
			}

			headerLinks := make(map[string]*int, len(entries))
			for _, entry := range entries {
				linkID := r.Pdf.AddLink()
				id := linkID // copy so each pointer is distinct
				headerLinks[entry.Title] = &id
			}

			r.SetTOCLinks(headerLinks)

			// Render the TOC page.
			r.Pdf.SetFont(r.Theme.Normal.Font, "B", 24)
			r.Pdf.Cell(40, 10, "Table of Contents")
			r.Pdf.Ln(30)

			for _, entry := range entries {
				if linkPtr, exists := headerLinks[entry.Title]; exists {
					r.Pdf.SetFont(r.Theme.Normal.Font, "", 12)
					r.Pdf.SetTextColor(100, 149, 237)
					indent := strings.Repeat("  ", entry.Level-1)
					r.Pdf.WriteLinkID(8, fmt.Sprintf("%s • %s", indent, entry.Title), *linkPtr)
					r.Pdf.Ln(15)
				}
			}
			r.Pdf.AddPage()

			return nil
		})
	}
}
