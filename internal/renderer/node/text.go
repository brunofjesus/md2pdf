package node

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessText handles *ast.Text nodes.
func ProcessText(ctx PdfContext, n ast.Node, _ bool) {
	node := n.(*ast.Text)
	currentStyle := ctx.PeekState().TextStyle
	ctx.SetStyler(currentStyle)
	s := string(node.Literal)
	s = strings.ReplaceAll(s, "\n", " ")
	ctx.Tracer("Text", s)

	if tableState.inCell {
		ctx.PeekState().CellInnerString += s
		ctx.PeekState().CellInnerStringStyle = &currentStyle
		return
	}
	switch node.Parent.(type) {
	case *ast.Link:
		ctx.WriteLink(currentStyle, s, ctx.PeekState().Destination)
	case *ast.Heading:
		tocLinks := ctx.GetTOCLinks()
		if len(tocLinks) > 0 {
			if linkPtr, exists := tocLinks[s]; exists {
				link := *linkPtr
				ctx.GetPdf().SetLink(link, -1, -1)
				ctx.Tracer("Text Heading", fmt.Sprintf("Set link for header '%s' with link ID: %d\n", s, link))
			} else {
				ctx.Tracer("Text Heading", fmt.Sprintf("Header '%s' not found in links map\n", s))
			}
		}
		ctx.Write(currentStyle, s)
	case *ast.BlockQuote:
		ctx.Tracer("Text BlockQuote", s)
		ctx.MultiCell(currentStyle, s)
	default:
		ctx.Write(currentStyle, s)
	}
}
