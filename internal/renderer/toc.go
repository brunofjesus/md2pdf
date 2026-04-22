package renderer

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Processor is satisfied by PdfRenderer and by any decorator that wraps it.
type Processor interface {
	Process(content []byte) error
}

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

// TOCDecorator wraps a PdfRenderer and prepends a table of contents page
// before delegating rendering to the underlying renderer.
type TOCDecorator struct {
	renderer *PdfRenderer
}

// NewTOCDecorator creates a TOCDecorator wrapping r.
func NewTOCDecorator(r *PdfRenderer) *TOCDecorator {
	return &TOCDecorator{renderer: r}
}

// Process generates the TOC page then delegates document rendering to the
// wrapped PdfRenderer.
func (d *TOCDecorator) Process(content []byte) error {
	entries, err := getTOCEntries(content)
	if err != nil {
		return fmt.Errorf("TOCDecorator: failed to collect TOC entries: %w", err)
	}

	headerLinks := make(map[string]*int, len(entries))
	for _, entry := range entries {
		linkID := d.renderer.Pdf.AddLink()
		id := linkID // copy so each pointer is distinct
		headerLinks[entry.Title] = &id
	}

	d.renderer.SetTOCLinks(headerLinks)

	// Render the TOC page.
	r := d.renderer
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

	return d.renderer.Process(content)
}
