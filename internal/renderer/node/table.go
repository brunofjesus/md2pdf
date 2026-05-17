package node

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

// tableState holds package-level table rendering state.
// TODO: move these into ContainerState to support concurrent rendering.
var tableState struct {
	cellWidths  []float64
	curDataCell int
	fill        bool
	inCell      bool
}

// ProcessTable handles *ast.Table entering/leaving.
func ProcessTable(ctx PdfContext, n ast.Node, entering bool) {
	if entering {
		ctx.Tracer("Table (entering)", "")
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Table.Header,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		ctx.Cr()
		ctx.PushState(x)
		tableState.fill = false
		tableState.cellWidths = ctx.GetColumnWidths(n)
	} else {
		wSum := 0.0
		for _, w := range tableState.cellWidths {
			wSum += w
		}
		ctx.GetPdf().CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

		ctx.PopState()
		ctx.Tracer("Table (leaving)", "")
		ctx.Cr()
	}
}

// ProcessTableHead handles *ast.TableHeader entering/leaving.
func ProcessTableHead(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.Tracer("TableHead (entering)", "")
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Table.Header,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		ctx.PushState(x)
	} else {
		ctx.PopState()
		ctx.Tracer("TableHead (leaving)", "")
	}
}

// ProcessTableBody handles *ast.TableBody entering/leaving.
func ProcessTableBody(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.Tracer("TableBody (entering)", "")
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Table.Body,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		ctx.PushState(x)
	} else {
		ctx.PopState()
		ctx.Tracer("TableBody (leaving)", "")
		ctx.GetPdf().Ln(-1)
	}
}

// ProcessTableRow handles *ast.TableRow entering/leaving.
func ProcessTableRow(ctx PdfContext, _ ast.Node, entering bool) {
	if entering {
		ctx.Tracer("TableRow (entering)", "")
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Table.Body,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		if ctx.PeekState().IsHeader {
			x.TextStyle = ctx.GetTheme().Table.Header
		}
		ctx.GetPdf().Ln(-1)

		tableState.curDataCell = 0
		ctx.PushState(x)
	} else {
		ctx.PopState()
		ctx.Tracer("TableRow (leaving)", "")
		tableState.fill = !tableState.fill
	}
}

// ProcessTableCell handles *ast.TableCell entering/leaving.
func ProcessTableCell(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.TableCell)
	if entering {
		ctx.Tracer("TableCell (entering)", "")
		x := &ContainerState{
			TextStyle:  ctx.GetTheme().Normal,
			ListKind:   NotList,
			LeftMargin: ctx.PeekState().LeftMargin,
		}
		if node.IsHeader {
			x.IsHeader = true
			x.TextStyle = ctx.GetTheme().Table.Header
			ctx.SetStyler(ctx.GetTheme().Table.Header)
		} else {
			x.TextStyle = ctx.GetTheme().Table.Body
			ctx.SetStyler(ctx.GetTheme().Table.Body)
			x.IsHeader = false
		}
		ctx.PushState(x)
		tableState.inCell = true
	} else {
		tableState.inCell = false
		cs := ctx.PopState()
		currentStyle := cs.TextStyle
		if cs.CellInnerStringStyle != nil {
			currentStyle = *cs.CellInnerStringStyle
		}
		s := cs.CellInnerString
		w := tableState.cellWidths[tableState.curDataCell]
		if cs.IsHeader {
			h, _ := ctx.GetPdf().GetFontSize()
			h += currentStyle.Spacing
			ctx.Tracer("... table header cell",
				fmt.Sprintf("Width=%v, height=%v", w, h))

			ctx.GetPdf().CellFormat(w, h, s, "1", 0, "C", true, 0, "")
		} else {
			h := currentStyle.Size + currentStyle.Spacing
			ctx.GetPdf().CellFormat(w, h, s, "LR", 0, "", tableState.fill, 0, "")
		}
		ctx.Tracer("TableCell (leaving)", "")
		tableState.curDataCell++
	}
}
