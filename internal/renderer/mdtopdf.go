/*
 * Markdown to PDF Converter
 * Available at http://github.com/solworktech/md2pdf
 *
 * Copyright © Cecil New <cecil.new@gmail.com>, Jesse Portnoy <jesse@packman.io>.
 * Distributed under the MIT License.
 * See README.md for details.
 *
 * Dependencies
 * This package depends on two other packages:
 *
 * Go Markdown processor
 *   Available at https://github.com/gomarkdown/markdown
 *
 * fpdf - a PDF document generator with high level support for
 *   text, drawing and images.
 *   Available at https://codeberg.org/go-pdf/fpdf
 */

// Package renderer converts markdown to PDF.
package renderer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/solworktech/md2pdf/v2/internal/colors"
	"github.com/solworktech/md2pdf/v2/internal/fonts"
	"github.com/solworktech/md2pdf/v2/internal/theme"
)

// Ensure PdfRenderer satisfies the Processor interface.
var _ Processor = (*PdfRenderer)(nil)

// RenderOption allows to define functions to configure the renderer
type RenderOption func(r *PdfRenderer)

// Theme [light|dark]
type Theme int

const (
	// DARK theme const
	DARK Theme = 1
	// LIGHT theme const
	LIGHT Theme = 2
	// CUSTOM theme const
	CUSTOM Theme = 3
)

// PdfRenderer is the struct to manage conversion of a markdown object
// to PDF format.
type PdfRenderer struct {
	// Pdf can be used to access the underlying created fpdf object
	// prior to processing the markdown source
	Pdf                *fpdf.Fpdf
	orientation, units string
	papersize, fontdir string

	// trace/log file if present
	pdfFile, tracerFile string
	w                   *bufio.Writer

	// default margins for safe keeping
	mleft, mtop, mright, mbottom float64

	Theme *theme.Theme
	// normal
	NormalEm float64
	// blockquote
	IndentValue float64

	cs states

	// update styling
	NeedCodeStyleUpdate       bool
	NeedBlockquoteStyleUpdate bool
	HorizontalRuleNewPage     bool
	InputBaseURL              string
	Extensions                parser.Extensions
	ColumnWidths              map[ast.Node][]float64

	tocLinks map[string]*int
}

// SetTOCLinks stores the heading→linkID map used by processText to place
// link anchors when each heading is rendered. Called by TOCDecorator.
func (r *PdfRenderer) SetTOCLinks(tocHeaders map[string]*int) {
	r.tocLinks = tocHeaders
}

// PdfRendererParams struct to hold params passed to NewPdfRenderer
type PdfRendererParams struct {
	Orientation, Papersz, PdfFile, TracerFile string
	Opts                                      []RenderOption
	Theme                                     Theme
	CustomThemeFile                           string
}

// NewPdfRenderer creates and configures an PdfRenderer object,
// which satisfies the Renderer interface.
func NewPdfRenderer(params PdfRendererParams) *PdfRenderer {
	r := new(PdfRenderer)

	// set filenames
	r.pdfFile = params.PdfFile
	r.tracerFile = params.TracerFile

	// Global things
	r.orientation = "portrait"
	if params.Orientation != "" {
		r.orientation = params.Orientation
	}

	r.units = "pt"
	r.papersize = "Letter"
	if params.Papersz != "" {
		r.papersize = params.Papersz
	}

	r.fontdir = "."

	r.Pdf = fpdf.New(r.orientation, r.units, r.papersize, r.fontdir)

	// Register Liberation Sans (SIL Open Font License) for all styles.
	// This provides full UTF-8 support including all Latin characters.
	r.Pdf.AddUTF8FontFromBytes("LiberationSans", "", fonts.LiberationSansRegular)
	r.Pdf.AddUTF8FontFromBytes("LiberationSans", "B", fonts.LiberationSansBold)
	r.Pdf.AddUTF8FontFromBytes("LiberationSans", "I", fonts.LiberationSansItalic)
	r.Pdf.AddUTF8FontFromBytes("LiberationSans", "BI", fonts.LiberationSansBoldItalic)

	// Register Liberation Mono (SIL Open Font License) for code and code blocks.
	r.Pdf.AddUTF8FontFromBytes("LiberationMono", "", fonts.LiberationMonoRegular)
	r.Pdf.AddUTF8FontFromBytes("LiberationMono", "B", fonts.LiberationMonoBold)
	r.Pdf.AddUTF8FontFromBytes("LiberationMono", "I", fonts.LiberationMonoItalic)
	r.Pdf.AddUTF8FontFromBytes("LiberationMono", "BI", fonts.LiberationMonoBoldItalic)

	switch params.Theme {
	case DARK:
		r.Theme = theme.DarkTheme()
	case LIGHT:
		r.Theme = theme.LightTheme()
	case CUSTOM:
		if params.CustomThemeFile != "" {
			r.Theme = theme.CustomTheme(params.CustomThemeFile)
		}
	default:
		r.Theme = theme.LightTheme()
	}

	r.Pdf.SetHeaderFunc(func() {
		w, h := r.Pdf.GetPageSize()
		dorect(r.Pdf, 0, 0, w, h, r.Theme.BackgroundColor)
	})

	r.Pdf.AddPage()
	// set default font
	r.setStyler(r.Theme.Normal)
	r.mleft, r.mtop, r.mright, r.mbottom = r.Pdf.GetMargins()
	r.NormalEm = r.Pdf.GetStringWidth("m")
	r.IndentValue = 3 * r.NormalEm

	r.cs = states{stack: make([]*containerState, 0)}
	initcurrent := &containerState{
		listkind:  notlist,
		textStyle: r.Theme.Normal, leftMargin: r.mleft,
	}
	r.cs.push(initcurrent)

	for _, o := range params.Opts {
		o(r)
	}

	return r
}

// NewPdfRendererWithDefaultStyler creates and configures an PdfRenderer object,
// which satisfies the Renderer interface.
// update default styler for normal
func NewPdfRendererWithDefaultStyler(orient, papersz, pdfFile, tracerFile string, defaultStyler theme.Styler, opts []RenderOption, theme Theme) *PdfRenderer {
	opts = append(opts, func(r *PdfRenderer) {
		r.Theme.Normal = defaultStyler
	})
	params := PdfRendererParams{
		Orientation: orient,
		Papersz:     papersz,
		PdfFile:     pdfFile,
		TracerFile:  tracerFile,
		Opts:        opts,
		Theme:       theme,
	}

	return NewPdfRenderer(params)
}

// Process takes the markdown content, parses it to generate the PDF
func (r *PdfRenderer) Process(content []byte) error {
	// try to open tracer
	var f *os.File
	var err error
	if r.tracerFile != "" {
		f, err = os.Create(r.tracerFile)
		if err != nil {
			return fmt.Errorf("os.Create() on tracefile error:%v", err)
		}
		defer f.Close()
		r.w = bufio.NewWriter(f)
		defer r.w.Flush()
	}

	err = r.Run(content)
	if err != nil {
		return fmt.Errorf("error on %v:%v", r.pdfFile, err)
	}

	err = r.Pdf.OutputFileAndClose(r.pdfFile)
	if err != nil {
		return fmt.Errorf("error on %v:%v", r.pdfFile, err)
	}

	return nil
}

// Run takes the markdown content, parses it but don't generate the PDF. you can access the PDF with youRenderer.Pdf
func (r *PdfRenderer) Run(content []byte) error {
	// Preprocess content by changing all CRLF to LF
	s := content
	s = markdown.NormalizeNewlines(s)

	p := parser.NewWithExtensions(r.Extensions)
	doc := markdown.Parse(s, p)

	setColumnWidths(doc, r)
	_ = markdown.Render(doc, r)

	return nil
}

// Parses all tables and sets the column width to the longest string in that column
func setColumnWidths(doc ast.Node, r *PdfRenderer) {
	columnWidths := map[ast.Node][]float64{}
	intable := false
	inheader := true
	cellnum := 0
	lengths := []float64{}
	textlength := float64(0)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch n := node.(type) {
		case *ast.Table:
			if entering {
				intable = true
			} else {
				intable = false
				columnWidths[node] = lengths
			}

		case *ast.TableHeader:
			inheader = entering
			if entering {
				lengths = []float64{}
			}
		case *ast.TableRow:
			if entering {
				cellnum = 0
			}
		case *ast.TableCell:
			if entering {
				if inheader {
					lengths = append(lengths, 0)
				}
			} else {
				textlength += textlength * 0.2

				currentMax := lengths[cellnum]
				if textlength > currentMax {
					lengths[cellnum] = textlength
				}
				textlength = 0
				cellnum++
			}
		case *ast.Text:
			if entering && intable {
				l := r.Pdf.GetStringWidth(string(n.Literal))
				textlength += l
			}
		}
		return ast.GoToNext
	})
	r.ColumnWidths = columnWidths
}

// UpdateParagraphStyler - update with default styler
func (r *PdfRenderer) UpdateParagraphStyler(defaultStyler theme.Styler) {
	initcurrent := &containerState{
		listkind:  notlist,
		textStyle: defaultStyler, leftMargin: r.mleft,
	}
	r.cs.push(initcurrent)
}

// UpdateCodeStyler - update code fill styler
func (r *PdfRenderer) UpdateCodeStyler() {
	r.NeedCodeStyleUpdate = true
}

// UpdateBlockquoteStyler - update Blockquote fill styler
func (r *PdfRenderer) UpdateBlockquoteStyler() {
	r.NeedBlockquoteStyleUpdate = true
}

func (r *PdfRenderer) setStyler(s theme.Styler) {
	// see https://github.com/solworktech/md2pdf/issues/18#issuecomment-2179694815
	// This does not address the root cause
	// (https://github.com/solworktech/md2pdf/issues/18#issuecomment-2179694815)
	// but it will correct all cases and is safer.
	if s.Style == "bb" {
		s.Style = "b"
	}
	r.Pdf.SetFont(s.Font, s.Style, s.Size)
	r.Pdf.SetTextColor(s.TextColor.Red, s.TextColor.Green, s.TextColor.Blue)
	r.Pdf.SetFillColor(s.FillColor.Red, s.FillColor.Green, s.FillColor.Blue)
}

func (r *PdfRenderer) write(s theme.Styler, t string) {
	// fmt.Printf("%s, %#v\n",t, s)
	r.Pdf.Write(s.Size+s.Spacing, t)
}

func (r *PdfRenderer) multiCell(s theme.Styler, t string) {
	r.Pdf.MultiCell(0, s.Size+s.Spacing, t, "", "", true)
}

func (r *PdfRenderer) writeLink(s theme.Styler, display, url string) {
	r.Pdf.WriteLinkString(s.Size+s.Spacing, display, url)
}

// RenderNode is a default renderer of a single node of a syntax tree. For
// block nodes it will be called twice: first time with entering=true, second
// time with entering=false, so that it could know when it's working on an open
// tag and when on close. It writes the result to w.
//
// The return value is a way to tell the calling walker to adjust its walk
// pattern: e.g. it can terminate the traversal by returning Terminate. Or it
// can ask the walker to skip a subtree of this node by returning SkipChildren.
// The typical behavior is to return GoToNext, which asks for the usual
// traversal to the next node.
// (above taken verbatim from the blackfriday v2 package)
func (r *PdfRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		r.processText(node)
	case *ast.Softbreak:
		r.tracer("Softbreak", "Output newline")
		r.cr()
	case *ast.Hardbreak:
		r.tracer("Hardbreak", "Output newline")
		r.cr()
	case *ast.Emph:
		r.processEmph(node, entering)
	case *ast.Strong:
		r.processStrong(node, entering)
	case *ast.Del:
		if entering {
			r.tracer("DEL (entering)", "Not handled")
		} else {
			r.tracer("DEL (leaving)", "Not handled")
		}
	case *ast.HTMLSpan:
		r.tracer("HTMLSpan", "Not handled")
	case *ast.Link:
		r.processLink(*node, entering)
	case *ast.Image:
		r.processImage(*node, entering)
	case *ast.Code:
		r.processCode(node)
	case *ast.Document:
		r.tracer("Document", "Not Handled")
	case *ast.Paragraph:
		r.processParagraph(node, entering)
	case *ast.BlockQuote:
		r.processBlockQuote(node, entering)
	case *ast.HTMLBlock:
		r.processHTMLBlock(node)
	case *ast.Heading:
		r.processHeading(*node, entering)
	case *ast.HorizontalRule:
		r.processHorizontalRule(node)
	case *ast.List:
		r.processList(*node, entering)
	case *ast.ListItem:
		r.processItem(*node, entering)
	case *ast.CodeBlock:
		r.processCodeblock(*node)
	case *ast.Table:
		r.processTable(node, entering)
	case *ast.TableHeader:
		r.processTableHead(node, entering)
	case *ast.TableBody:
		r.processTableBody(node, entering)
	case *ast.TableRow:
		r.processTableRow(node, entering)
	case *ast.TableCell:
		r.processTableCell(*node, entering)
	default:
		fmt.Printf("Unknown node type: %T. Skipping\n", node)
	}
	return ast.GoToNext
}

// RenderHeader is not supported.
func (r *PdfRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	r.tracer("RenderHeader", "Not handled")
}

// RenderFooter is not supported.
func (r *PdfRenderer) RenderFooter(w io.Writer, _ ast.Node) {
}

func (r *PdfRenderer) cr() {
	LH := r.cs.peek().textStyle.Size + r.cs.peek().textStyle.Spacing
	r.tracer("cr()", fmt.Sprintf("LH=%v", LH))
	r.write(r.cs.peek().textStyle, "\n")
}

// Tracer traces parse and pdf generation activity.
func (r *PdfRenderer) tracer(source, msg string) {
	if r.tracerFile != "" {
		indent := strings.Repeat("-", len(r.cs.stack)-1)
		_, _ = fmt.Fprintf(r.w, "%v[%v] %v\n", indent, source, msg)
	}
}

func dorect(doc *fpdf.Fpdf, x, y, w, h float64, color colors.Color) {
	doc.SetFillColor(color.Red, color.Green, color.Blue)
	doc.Rect(x, y, w, h, "F")
}
