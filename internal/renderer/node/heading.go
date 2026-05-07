package node

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessHeading handles *ast.Heading entering/leaving.
func ProcessHeading(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.Heading)
	if entering {
		ctx.Cr()
		var style = ctx.GetTheme().Normal
		switch node.Level {
		case 1:
			style = ctx.GetTheme().H1
		case 2:
			style = ctx.GetTheme().H2
		case 3:
			style = ctx.GetTheme().H3
		case 4:
			style = ctx.GetTheme().H4
		case 5:
			style = ctx.GetTheme().H5
		case 6:
			style = ctx.GetTheme().H6
		}
		ctx.Tracer(fmt.Sprintf("Heading (%d, entering)", node.Level),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		x := &ContainerState{
			TextStyle:  style,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		ctx.PushState(x)
	} else {
		ctx.Tracer("Heading (leaving)", "")
		ctx.Cr()
		ctx.PopState()
	}
}
