package node

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessHeading handles *ast.Heading entering/leaving.
func ProcessHeading(ctx PdfContext, n ast.Node, entering bool) {
	node, ok := n.(*ast.Heading)
	if !ok {
		ctx.Tracer("Heading: not a Heading", "")
		return
	}

	if entering {
		ctx.Cr()
		style := ctx.GetTheme().Normal

		switch node.Level {
		case 1:
			style = ctx.GetTheme().Heading.H1
		case 2:
			style = ctx.GetTheme().Heading.H2
		case 3:
			style = ctx.GetTheme().Heading.H3
		case 4:
			style = ctx.GetTheme().Heading.H4
		case 5:
			style = ctx.GetTheme().Heading.H5
		case 6:
			style = ctx.GetTheme().Heading.H6
		}

		ctx.Tracer(
			fmt.Sprintf("Heading (%d, entering)", node.Level),
			ast.ToString(node.AsContainer()),
		)

		x := &ContainerState{
			TextStyle:  style,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		ctx.PushState(x)
	} else {
		ctx.Cr()

		lineHeight := ctx.GetTheme().Heading.Line.Height
		if ctx.GetTheme().Heading.Line.Height > 0 {
			lineColor := ctx.GetTheme().Heading.Line.Color

			pdf := ctx.GetPdf()
			x, y := pdf.GetXY()
			lm, _, _, _ := pdf.GetMargins()
			w, _ := pdf.GetPageSize()
			newx := w - lm

			ctx.Tracer("... Drawing underline from X,Y", fmt.Sprintf("%v,%v", x, y))
			pdf.MoveTo(x, y)
			ctx.Tracer("...   To X,Y", fmt.Sprintf("%v,%v", newx, y))
			pdf.LineTo(newx, y)
			pdf.SetLineWidth(lineHeight)
			pdf.SetDrawColor(lineColor.Red, lineColor.Green, lineColor.Blue)
			pdf.DrawPath("D")
		}

		ctx.PopState()
		ctx.Tracer("Heading (leaving)", "")
	}
}
