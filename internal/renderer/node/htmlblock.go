package node

import "github.com/gomarkdown/markdown/ast"

// ProcessHTMLBlock handles *ast.HTMLBlock nodes.
func ProcessHTMLBlock(ctx PdfContext, n ast.Node, _ bool) {
	ctx.Tracer("HTMLBlock", string(n.AsLeaf().Literal))
	ctx.Cr()
	ctx.SetStyler(ctx.GetTheme().Backtick)
	ctx.GetPdf().CellFormat(0, ctx.GetTheme().Backtick.Size,
		string(n.AsLeaf().Literal), "", 1, "LT", true, 0, "")
	ctx.Cr()
}
