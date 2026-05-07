package node

import (
	"codeberg.org/go-pdf/fpdf"
	"github.com/brunofjesus/md2pdf/internal/theme"
	"github.com/gomarkdown/markdown/ast"
)

// ListType describes the kind of list being rendered.
type ListType int

const (
	NotList ListType = iota
	Unordered
	Ordered
	Definition
)

func (n ListType) String() string {
	switch n {
	case NotList:
		return "Not a List"
	case Unordered:
		return "Unordered"
	case Ordered:
		return "Ordered"
	case Definition:
		return "Definition"
	}
	return ""
}

// ContainerState holds the rendering state for a single level of the AST
// container nesting (list, blockquote, table cell, link, heading, etc.).
type ContainerState struct {
	TextStyle      theme.Styler
	LeftMargin     float64
	FirstParagraph bool

	// List fields
	ListKind   ListType
	ItemNumber int // only meaningful for ordered lists

	// Link fields
	Destination string

	// Table cell fields
	IsHeader             bool
	CellInnerString      string
	CellInnerStringStyle *theme.Styler
}

// NodeProcessor is a function that handles rendering of a particular AST node
// type. entering is true when the walker enters the node, false when leaving.
type NodeProcessor func(ctx PdfContext, node ast.Node, entering bool)

// PdfContext is the interface that node processors use to interact with the
// PDF renderer. It abstracts PdfRenderer so the node/ package has no import
// dependency on the renderer package.
type PdfContext interface {
	// Tracing / logging
	Tracer(source, msg string)

	// Writing primitives
	Cr()
	Write(s theme.Styler, text string)
	MultiCell(s theme.Styler, text string)
	WriteLink(s theme.Styler, display, url string)

	// Style application
	SetStyler(s theme.Styler)

	// State stack
	PushState(s *ContainerState)
	PopState() *ContainerState
	PeekState() *ContainerState
	ParentState() *ContainerState
	StackDepth() int

	// Layout
	SetLeftMargin(margin float64)

	// Theme access
	GetTheme() *theme.Theme

	// Computed values
	GetIndentValue() float64
	GetNormalEm() float64

	// Input context
	GetInputBaseURL() string

	// TOC links
	GetTOCLinks() map[string]*int

	// Table column widths
	GetColumnWidths(node ast.Node) []float64

	// Direct PDF access — pragmatic escape hatch for node processors
	// that need low-level fpdf operations (images, table cells, etc.).
	GetPdf() *fpdf.Fpdf
}
