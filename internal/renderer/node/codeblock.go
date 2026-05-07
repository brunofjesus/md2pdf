package node

import (
	"fmt"
	"math"
	"strings"

	syntaxhighlight "github.com/brunofjesus/md2pdf/internal/highlight"
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
	if ctx.GetTheme().CodeTabWidth > 0 {
		linesWrapped = strings.ReplaceAll(linesWrapped, "\t", strings.Repeat(" ", ctx.GetTheme().CodeTabWidth))
	}
	matches := h.HighlightString(linesWrapped)

	ctx.SetStyler(ctx.GetTheme().Code)
	ctx.Cr()

	lines := strings.Split(linesWrapped, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	codeTheme := ctx.GetTheme().Code
	lineH := codeTheme.Size + codeTheme.Spacing
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
				switch group {
				case highlight.Groups["default"],
					highlight.Groups[""]:
					ctx.SetStyler(codeTheme)
				case highlight.Groups["statement"],
					highlight.Groups["green"]:
					ctx.GetPdf().SetTextColor(42, 170, 138)
				case highlight.Groups["identifier"],
					highlight.Groups["blue"]:
					ctx.GetPdf().SetTextColor(137, 207, 240)
				case highlight.Groups["preproc"]:
					ctx.GetPdf().SetTextColor(255, 80, 80)
				case highlight.Groups["special"],
					highlight.Groups["type.keyword"],
					highlight.Groups["red"]:
					ctx.GetPdf().SetTextColor(255, 80, 80)
				case highlight.Groups["constant"],
					highlight.Groups["constant.number"],
					highlight.Groups["constant.bool"],
					highlight.Groups["symbol.brackets"],
					highlight.Groups["identifier.var"],
					highlight.Groups["cyan"]:
					ctx.GetPdf().SetTextColor(0, 136, 163)
				case highlight.Groups["constant.specialChar"],
					highlight.Groups["constant.string.url"],
					highlight.Groups["constant.string"],
					highlight.Groups["magenta"]:
					ctx.GetPdf().SetTextColor(255, 0, 255)
				case highlight.Groups["type"],
					highlight.Groups["symbol"],
					highlight.Groups["symbol.operator"],
					highlight.Groups["symbol.tag.extended"],
					highlight.Groups["yellow"]:
					ctx.GetPdf().SetTextColor(255, 165, 0)
				case highlight.Groups["comment"],
					highlight.Groups["high.green"]:
					ctx.GetPdf().SetTextColor(82, 204, 0)
				default:
					fmt.Printf("Unknown group: %s\n", group)
					ctx.SetStyler(codeTheme)
				}
			}
			ctx.GetPdf().Write(lineH, string(c))
			colN++
		}
	})

	ctx.SetStyler(codeTheme)
}

func outputUnhighlightedCodeBlock(ctx PdfContext, codeBlock string) {
	ctx.Cr()
	ctx.SetStyler(ctx.GetTheme().Backtick)
	if ctx.GetTheme().CodeTabWidth > 0 {
		codeBlock = strings.ReplaceAll(codeBlock, "\t", strings.Repeat(" ", ctx.GetTheme().CodeTabWidth))
	}
	ctx.MultiCell(ctx.GetTheme().Code, codeBlock)
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
		ctx.SetStyler(codeTheme)
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
	if codeTheme.FillColor != bgColor {
		if h := rectHeightFrom(0, startY); h > 0 {
			drawBg(startY, h)
			pdf.SetXY(startX, startY)
		}
	}

	for lineN, l := range lines {
		if pdf.GetY()+lineHeights[lineN] > usableH {
			pdf.AddPage()
			newY := pdf.GetY()
			if codeTheme.FillColor != bgColor {
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
	ctx.Tracer("processCode", fmt.Sprintf("%s", string(n.AsLeaf().Literal)))
	ctx.Write(ctx.GetTheme().Normal, " ")
	ctx.Tracer("Code (entering)", "")
	codeTheme := ctx.GetTheme().Code
	ctx.SetStyler(codeTheme)
	s := string(n.AsLeaf().Literal)
	hw := ctx.GetPdf().GetStringWidth(s)
	h := codeTheme.Size
	ctx.GetPdf().CellFormat(hw, h, s, "", 0, "C", true, 0, "")
}
