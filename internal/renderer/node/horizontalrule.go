package node

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

// HorizontalRulePageBreakProcessor starts a new page on HR.
func HorizontalRulePageBreakProcessor(ctx PdfContext, _ ast.Node, _ bool) {
	ctx.Tracer("HorizontalRulePageBreakProcessor", "")
	ctx.GetPdf().AddPage()
}

// HorizontalRuleLineProcessor draws a horizontal line on HR.
func HorizontalRuleLineProcessor(ctx PdfContext, _ ast.Node, _ bool) {
	ctx.Tracer("HorizontalRuleLineProcessor", "")
	ctx.Cr()
	pdf := ctx.GetPdf()
	x, y := pdf.GetXY()
	lm, _, _, _ := pdf.GetMargins()
	w, _ := pdf.GetPageSize()
	newx := w - lm
	ctx.Tracer("... From X,Y", fmt.Sprintf("%v,%v", x, y))
	pdf.MoveTo(x, y)
	ctx.Tracer("...   To X,Y", fmt.Sprintf("%v,%v", newx, y))
	pdf.LineTo(newx, y)
	pdf.SetLineWidth(3)
	pdf.SetFillColor(200, 200, 200)
	pdf.DrawPath("F")
	ctx.Cr()
}
