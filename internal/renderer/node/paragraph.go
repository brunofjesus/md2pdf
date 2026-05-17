package node

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessParagraph handles *ast.Paragraph entering/leaving.
func ProcessParagraph(ctx PdfContext, n ast.Node, entering bool) {
	node, ok := n.(*ast.Paragraph)
	if !ok {
		ctx.Tracer("Paragraph: not a Paragraph", "")
		return
	}

	ctx.SetStyler(ctx.GetTheme().Normal)

	if entering { //nolint:nestif
		ctx.Tracer("Paragraph (entering)", "")
		lm, tm, rm, bm := ctx.GetPdf().GetMargins()
		ctx.Tracer("... Margins (left, top, right, bottom:",
			fmt.Sprintf("%v %v %v %v", lm, tm, rm, bm))

		if IsListItem(node.Parent) {
			t := ctx.PeekState().ListKind
			if t == Unordered || t == Ordered || t == Definition {
				if ctx.PeekState().FirstParagraph {
					ctx.Tracer("First Para within a list", "breaking")
				} else {
					ctx.Tracer("Not First Para within a list", "indent etc.")
					ctx.Cr()
				}
			}

			return
		}

		ctx.Cr()
	} else {
		ctx.Tracer("Paragraph (leaving)", "")
		lm, tm, rm, bm := ctx.GetPdf().GetMargins()
		ctx.Tracer("... Margins (left, top, right, bottom:",
			fmt.Sprintf("%v %v %v %v", lm, tm, rm, bm))

		if IsListItem(node.Parent) {
			t := ctx.PeekState().ListKind
			if t == Unordered || t == Ordered || t == Definition {
				if ctx.PeekState().FirstParagraph {
					ctx.PeekState().FirstParagraph = false
				} else {
					ctx.Tracer("Not First Para within a list", "")
					ctx.Cr()
				}
			}

			return
		}

		ctx.Cr()
	}
}

// ProcessBlockQuote handles *ast.BlockQuote entering/leaving.
func ProcessBlockQuote(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.Tracer("BlockQuote (entering)", "")
		curleftmargin, _, _, _ := ctx.GetPdf().GetMargins()
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Blockquote,
			ListKind:   NotList,
			LeftMargin: curleftmargin + ctx.GetIndentValue(),
		}
		ctx.PushState(x)
		ctx.SetLeftMargin(curleftmargin + ctx.GetIndentValue())
	} else {
		ctx.Tracer("BlockQuote (leaving)", "")
		curleftmargin, _, _, _ := ctx.GetPdf().GetMargins()
		ctx.SetLeftMargin(curleftmargin - ctx.GetIndentValue())
		ctx.PopState()
		ctx.Cr()
	}
}
