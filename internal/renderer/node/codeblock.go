package node

import (
	"fmt"
	"math"
	"strings"

	syntaxhighlight "github.com/brunofjesus/md2pdf/v3/internal/highlight"
	"github.com/gomarkdown/markdown/ast"
	highlight "github.com/jessp01/gohighlight"
	"github.com/mitchellh/go-wordwrap"
)

// ProcessCodeBlock handles *ast.CodeBlock nodes.
func ProcessCodeBlock(ctx PdfContext, n ast.Node, _ bool) {
	node := n.(*ast.CodeBlock)
	ctx.Tracer("Codeblock", fmt.Sprintf("%v", ast.ToString(node.AsLeaf())))

	currentStyle := ctx.PeekState().TextStyle
	ctx.SetStyler(currentStyle)

	if len(node.Info) < 1 {
		outputUnhighlightedCodeBlock(ctx, string(node.Literal))
		return
	}

	if strings.HasPrefix(string(node.Literal), "<script") && string(node.Info) == "html" {
		node.Info = []byte("javascript")
	}
	syntaxFile, lerr := syntaxhighlight.Files.ReadFile(string(node.Info) + ".yaml")
	if lerr != nil {
		outputUnhighlightedCodeBlock(ctx, string(node.Literal))
		return
	}
	syntaxDef, _ := highlight.ParseDef(syntaxFile)
	h := highlight.NewHighlighter(syntaxDef)
	linesWrapped := wordwrap.WrapString(string(node.Literal), 90)
	if ctx.GetTheme().Code.TabWidth > 0 {
		linesWrapped = strings.ReplaceAll(linesWrapped, "\t", strings.Repeat(" ", ctx.GetTheme().Code.TabWidth))
	}
	matches := h.HighlightString(linesWrapped)

	ctx.SetStyler(ctx.GetTheme().Code.Text)
	ctx.Cr()

	lines := strings.Split(linesWrapped, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Build reverse lookup from group ID to group name
	reverseGroups := make(map[highlight.Group]string, len(highlight.Groups))
	for name, id := range highlight.Groups {
		reverseGroups[id] = name
	}

	codeTheme := ctx.GetTheme().Code
	cbColors := ctx.GetTheme().Code.Colors
	lineH := codeTheme.Text.Size + codeTheme.Text.Spacing
	lm, _, rm, _ := ctx.GetPdf().GetMargins()
	pw, _ := ctx.GetPdf().GetPageSize()
	availW := pw - lm - rm

	lineHeights := make([]float64, len(lines))
	for i, l := range lines {
		w := ctx.GetPdf().GetStringWidth(l)
		if w <= 0 {
			lineHeights[i] = lineH
		} else {
			lineHeights[i] = math.Ceil(w/availW) * lineH
		}
	}

	drawCodeFill(ctx, lines, lineHeights, func(lineN int, l string) {
		colN := 0
		for _, c := range l {
			if group, ok := matches[lineN][colN]; ok {
				groupName := reverseGroups[group]
				if color, found := cbColors[groupName]; found {
					ctx.GetPdf().SetTextColor(color.Red, color.Green, color.Blue)
				} else {
					ctx.SetStyler(codeTheme.Text)
				}
			}
			ctx.GetPdf().Write(lineH, string(c))
			colN++
		}
	})

	ctx.SetStyler(codeTheme.Text)
}

func outputUnhighlightedCodeBlock(ctx PdfContext, codeBlock string) {
	ctx.Cr()
	ctx.SetStyler(ctx.GetTheme().Backtick)
	if ctx.GetTheme().Code.TabWidth > 0 {
		codeBlock = strings.ReplaceAll(codeBlock, "\t", strings.Repeat(" ", ctx.GetTheme().Code.TabWidth))
	}
	ctx.MultiCell(ctx.GetTheme().Code.Text, codeBlock)
}

// drawCodeFill manages background rectangles and page breaks for highlighted code blocks.
func drawCodeFill(ctx PdfContext, lines []string, lineHeights []float64, renderLine func(lineN int, l string)) {
	pdf := ctx.GetPdf()
	codeTheme := ctx.GetTheme().Code
	bgColor := ctx.GetTheme().BackgroundColor

	lm, _, rm, bm := pdf.GetMargins()
	pw, ph := pdf.GetPageSize()
	availW := pw - lm - rm
	usableH := ph - bm

	drawBg := func(y, height float64) {
		ctx.SetStyler(codeTheme.Text)
		pdf.Rect(lm, y, availW, height, "F")
	}

	rectHeightFrom := func(from int, pageTopY float64) float64 {
		h := 0.0
		for i := from; i < len(lines); i++ {
			if pageTopY+h+lineHeights[i] > usableH {
				break
			}
			h += lineHeights[i]
		}
		return h
	}

	autoBreak, pbMargin := pdf.GetAutoPageBreak()
	pdf.SetAutoPageBreak(false, pbMargin)
	defer pdf.SetAutoPageBreak(autoBreak, pbMargin)

	startX, startY := pdf.GetXY()
	if codeTheme.Text.FillColor != bgColor {
		if h := rectHeightFrom(0, startY); h > 0 {
			drawBg(startY, h)
			pdf.SetXY(startX, startY)
		}
	}

	for lineN, l := range lines {
		if pdf.GetY()+lineHeights[lineN] > usableH {
			pdf.AddPage()
			newY := pdf.GetY()
			if codeTheme.Text.FillColor != bgColor {
				if h := rectHeightFrom(lineN, newY); h > 0 {
					drawBg(newY, h)
				}
			}
			pdf.SetX(lm)
		}

		renderLine(lineN, l)
		ctx.Cr()
	}
}

// ProcessCode handles inline *ast.Code nodes.
func ProcessCode(ctx PdfContext, n ast.Node, _ bool) {
	ctx.Tracer("processCode", string(n.AsLeaf().Literal))
	ctx.Write(ctx.GetTheme().Normal, " ")
	ctx.Tracer("Code (entering)", "")
	codeTheme := ctx.GetTheme().Code
	ctx.SetStyler(codeTheme.Text)
	s := string(n.AsLeaf().Literal)
	hw := ctx.GetPdf().GetStringWidth(s)
	h := codeTheme.Text.Size
	ctx.GetPdf().CellFormat(hw, h, s, "", 0, "C", true, 0, "")
}
