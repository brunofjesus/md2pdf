package node

import (
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

// ProcessList handles *ast.List entering/leaving.
func ProcessList(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.List)
	kind := Unordered
	if node.ListFlags&ast.ListTypeOrdered != 0 {
		kind = Ordered
	}
	if node.ListFlags&ast.ListTypeDefinition != 0 {
		kind = Definition
	}
	ctx.SetStyler(ctx.GetTheme().Normal)
	if entering {
		ctx.Tracer(fmt.Sprintf("%v List (entering)", kind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		ctx.SetLeftMargin(ctx.PeekState().LeftMargin + ctx.GetIndentValue())
		ctx.Tracer("... List Left Margin",
			fmt.Sprintf("set to %v", ctx.PeekState().LeftMargin+ctx.GetIndentValue()))
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Normal,
			ItemNumber: 0,
			ListKind:   kind,
			LeftMargin: ctx.PeekState().LeftMargin + ctx.GetIndentValue(),
		}
		ctx.PushState(x)
	} else {
		ctx.Tracer(fmt.Sprintf("%v List (leaving)", kind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		ctx.SetLeftMargin(ctx.PeekState().LeftMargin - ctx.GetIndentValue())
		ctx.Tracer("... Reset List Left Margin",
			fmt.Sprintf("re-set to %v", ctx.PeekState().LeftMargin-ctx.GetIndentValue()))
		ctx.PopState()
		if ctx.StackDepth() < 2 {
			ctx.Cr()
		}
	}
}

// ProcessItem handles *ast.ListItem entering/leaving.
func ProcessItem(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.ListItem)
	if entering {
		ctx.Tracer(fmt.Sprintf("%v Item (entering) #%v",
			ctx.PeekState().ListKind, ctx.PeekState().ItemNumber+1),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		ctx.Cr()
		x := &ContainerState{
			TextStyle:      ctx.GetTheme().Normal,
			ItemNumber:     ctx.PeekState().ItemNumber + 1,
			ListKind:       ctx.PeekState().ListKind,
			FirstParagraph: true,
			LeftMargin:     ctx.PeekState().LeftMargin,
		}
		ctx.PushState(x)
		normalEm := ctx.GetNormalEm()
		if ctx.PeekState().ListKind == Unordered {
			bulletChar := "•"
			currFontSize, _ := ctx.GetPdf().GetFontSize()
			if node.BulletChar != 45 {
				bulletChar = "▪"
				ctx.GetPdf().SetFont("", "", 25)
			}
			ctx.GetPdf().CellFormat(4*normalEm, ctx.GetTheme().Normal.Size+ctx.GetTheme().Normal.Spacing,
				bulletChar,
				"", 0, "RB", false, 0, "")
			ctx.GetPdf().SetFont("", "", currFontSize)
		} else if ctx.PeekState().ListKind == Ordered {
			ctx.GetPdf().CellFormat(4*normalEm, ctx.GetTheme().Normal.Size+ctx.GetTheme().Normal.Spacing,
				fmt.Sprintf("%v.", ctx.PeekState().ItemNumber),
				"", 0, "RB", false, 0, "")
		}
		ctx.SetLeftMargin(ctx.PeekState().LeftMargin + (4 * normalEm))
		ctx.GetPdf().SetX(ctx.PeekState().LeftMargin + (4 * normalEm))
	} else {
		ctx.Tracer(fmt.Sprintf("%v Item (leaving)",
			ctx.PeekState().ListKind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		ctx.SetLeftMargin(ctx.PeekState().LeftMargin)
		ctx.ParentState().ItemNumber++
		ctx.PopState()
	}
}

// IsListItem checks whether node is an *ast.ListItem.
func IsListItem(node ast.Node) bool {
	_, ok := node.(*ast.ListItem)
	return ok
}

// ProcessEmph handles *ast.Emph entering/leaving.
func ProcessEmph(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.Tracer("Emph (entering)", "")
		ctx.PeekState().TextStyle.Style += "i"
	} else {
		ctx.Tracer("Emph (leaving)", "")
		ctx.PeekState().TextStyle.Style = strings.ReplaceAll(
			ctx.PeekState().TextStyle.Style, "i", "")
	}
}

// ProcessStrong handles *ast.Strong entering/leaving.
func ProcessStrong(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.PeekState().TextStyle.Style += "b"
		ctx.Tracer("Strong (entering)", "")
	} else {
		ctx.Tracer("Strong (leaving)", "")
		ctx.PeekState().TextStyle.Style = strings.ReplaceAll(
			ctx.PeekState().TextStyle.Style, "b", "")
	}
}
