package renderer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/brunofjesus/md2pdf/v3/internal/colors"
	"github.com/brunofjesus/md2pdf/v3/internal/fonts"
	"github.com/brunofjesus/md2pdf/v3/internal/renderer/node"
	"github.com/brunofjesus/md2pdf/v3/internal/theme"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Ensure PdfRenderer satisfies node.PdfContext.
var _ node.PdfContext = (*PdfRenderer)(nil)

// RenderOption allows to define functions to configure the renderer.
type RenderOption func(r *PdfRenderer)

// NodeProcessor is a function type that takes a PdfRenderer, an AST node and a
// boolean indicating whether the node is being entered or exited.
//
// It is used to define processing functions for different types of AST nodes
// during PDF generation.
type NodeProcessor func(r *PdfRenderer, node ast.Node, entering bool)

// Theme [light|dark].
type Theme int

const (
	// DARK theme const.
	DARK Theme = 1
	// LIGHT theme const.
	LIGHT Theme = 2
	// CUSTOM theme const.
	CUSTOM Theme = 3
)

// identCharLen is the number of characters to per identation level.
const identCharLen float64 = 2

// PdfRenderer is the struct to manage conversion of a markdown object
// to PDF format.
type PdfRenderer struct {
	// Pdf can be used to access the underlying created fpdf object
	// prior to processing the markdown source
	Pdf                *fpdf.Fpdf
	orientation, units string
	papersize, fontdir string

	// trace/log file if present
	tracerFile string
	w          *bufio.Writer

	// default margins for safe keeping
	mleft, mtop, mright, mbottom float64

	Theme *theme.Theme
	// normal
	NormalEm float64
	// blockquote
	IndentValue float64

	cs states

	InputBaseURL string
	Extensions   parser.Extensions
	ColumnWidths map[ast.Node][]float64

	tocLinks map[string]*int

	// nodeProcessors maps AST node type names to their processor functions.
	nodeProcessors map[string]node.Processor

	// preProcessors are functions that run before the main rendering pass.
	preProcessors []func(content []byte) error
}

// PdfRendererParams struct to hold params passed to NewPdfRenderer.
type PdfRendererParams struct {
	Title                             string
	Orientation, PageSize, TracerFile string
	Opts                              []RenderOption
	Theme                             Theme
	CustomThemeFile                   string
}

// NewPdfRenderer creates and configures an PdfRenderer object,
// which satisfies the Renderer interface.
func NewPdfRenderer(params PdfRendererParams) *PdfRenderer {
	r := new(PdfRenderer)

	// set filenames
	r.tracerFile = params.TracerFile

	// Global things
	r.orientation = "portrait"
	if params.Orientation != "" {
		r.orientation = params.Orientation
	}

	r.units = "pt"
	r.papersize = "A4"

	if params.PageSize != "" {
		r.papersize = params.PageSize
	}

	r.fontdir = "."

	r.Pdf = fpdf.New(r.orientation, r.units, r.papersize, r.fontdir)

	// Register Liberation Sans (SIL Open Font License) for all styles.
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
	r.SetStyler(r.Theme.Normal)
	r.mleft, r.mtop, r.mright, r.mbottom = r.Pdf.GetMargins()
	r.NormalEm = r.Pdf.GetStringWidth("m")
	r.IndentValue = identCharLen * r.NormalEm

	r.cs = states{stack: make([]*node.ContainerState, 0)}
	initcurrent := &node.ContainerState{
		ListKind:  node.NotList,
		TextStyle: r.Theme.Normal, LeftMargin: r.mleft,
	}
	r.cs.push(initcurrent)

	r.Pdf.SetSubject(params.Title, true)
	r.Pdf.SetTitle(params.Title, true)

	// Register default node processors.
	r.nodeProcessors = map[string]node.Processor{
		"Text":           node.ProcessText,
		"Emph":           node.ProcessEmph,
		"Strong":         node.ProcessStrong,
		"Link":           node.ProcessLink,
		"Image":          node.ProcessImage,
		"Code":           node.ProcessCode,
		"Paragraph":      node.ProcessParagraph,
		"BlockQuote":     node.ProcessBlockQuote,
		"Heading":        node.ProcessHeading,
		"HTMLBlock":      node.ProcessHTMLBlock,
		"List":           node.ProcessList,
		"ListItem":       node.ProcessItem,
		"CodeBlock":      node.ProcessCodeBlock,
		"Table":          node.ProcessTable,
		"TableHeader":    node.ProcessTableHead,
		"TableBody":      node.ProcessTableBody,
		"TableRow":       node.ProcessTableRow,
		"TableCell":      node.ProcessTableCell,
		"HorizontalRule": node.HorizontalRuleLineProcessor,
	}

	for _, o := range params.Opts {
		o(r)
	}

	if r.Extensions == 0 {
		WithDefaultMarkdownParsingExtensions()(r)
	}

	return r
}

// ---------------------------------------------------------------------------
// node.PdfContext implementation
// ---------------------------------------------------------------------------

// Tracer writes trace messages to the tracer file if it is set. The message is
// indented according to the current depth of the container state stack.
func (r *PdfRenderer) Tracer(source, msg string) {
	if r.tracerFile != "" {
		indent := strings.Repeat("-", len(r.cs.stack)-1)
		_, _ = fmt.Fprintf(r.w, "%v[%v] %v\n", indent, source, msg)
	}
}

// Cr writes a newline to the PDF document. The line height is determined by the
// current text style's size and spacing.
func (r *PdfRenderer) Cr() {
	LH := r.cs.peek().TextStyle.Size + r.cs.peek().TextStyle.Spacing
	r.Tracer("cr()", fmt.Sprintf("LH=%v", LH))
	r.Write(r.cs.peek().TextStyle, "\n")
}

// Write writes text to the PDF document using the provided styler for font size
// and spacing.
func (r *PdfRenderer) Write(s theme.Styler, t string) {
	r.Pdf.Write(s.Size+s.Spacing, t)
}

// MultiCell writes text to the PDF document in a multi-cell format, allowing for
// line breaks. The cell width is set to 0 (full width), and the line height is
// determined by the provided styler's size and spacing.
func (r *PdfRenderer) MultiCell(s theme.Styler, t string) {
	r.Pdf.MultiCell(0, s.Size+s.Spacing, t, "", "", true)
}

// WriteLink writes a hyperlink to the PDF document using the provided styler for
// font size and spacing. The display text and URL are specified as parameters.
func (r *PdfRenderer) WriteLink(s theme.Styler, display, url string) {
	r.Pdf.WriteLinkString(s.Size+s.Spacing, display, url)
}

// SetStyler sets the font and colors in the PDF document according to the provided
// styler. If the styler's style is "bb", it is treated as "b" (bold) for compatibility.
func (r *PdfRenderer) SetStyler(styler theme.Styler) {
	if styler.Style == "bb" {
		styler.Style = "b"
	}

	r.Pdf.SetFont(styler.Font, styler.Style, styler.Size)
	r.Pdf.SetTextColor(styler.TextColor.Red, styler.TextColor.Green, styler.TextColor.Blue)
	r.Pdf.SetFillColor(styler.FillColor.Red, styler.FillColor.Green, styler.FillColor.Blue)
}

// PushState pushes a new container state onto the stack.
func (r *PdfRenderer) PushState(s *node.ContainerState) { r.cs.push(s) }

// PopState pops the current container state from the stack and returns it.
func (r *PdfRenderer) PopState() *node.ContainerState { return r.cs.pop() }

// PeekState returns the current container state without modifying the stack.
func (r *PdfRenderer) PeekState() *node.ContainerState { return r.cs.peek() }

// ParentState returns the parent container state, which is the second-to-last.
func (r *PdfRenderer) ParentState() *node.ContainerState { return r.cs.parent() }

// StackDepth returns the current depth of the container state stack.
func (r *PdfRenderer) StackDepth() int { return len(r.cs.stack) }

// SetLeftMargin sets the left margin of the PDF document to the specified value.
func (r *PdfRenderer) SetLeftMargin(margin float64) {
	r.Pdf.SetLeftMargin(margin)
}

// GetTheme returns the current theme being used by the renderer.
func (r *PdfRenderer) GetTheme() *theme.Theme { return r.Theme }

// GetIndentValue returns the current indent value used for lists and
// blockquotes.
func (r *PdfRenderer) GetIndentValue() float64 { return r.IndentValue }

// GetNormalEm returns the width of the character "m" in the current font.
func (r *PdfRenderer) GetNormalEm() float64 { return r.NormalEm }

// GetInputBaseURL returns the base URL for resolving relative links in the
// markdown document.
func (r *PdfRenderer) GetInputBaseURL() string { return r.InputBaseURL }

// GetTOCLinks returns the map of heading text to link IDs used for generating
// the Table of Contents.
func (r *PdfRenderer) GetTOCLinks() map[string]*int { return r.tocLinks }

// GetPdf returns the underlying fpdf.Fpdf object, allowing direct access to
// its methods and properties.
func (r *PdfRenderer) GetPdf() *fpdf.Fpdf { return r.Pdf }

// GetColumnWidths returns the map of AST nodes to their corresponding column
// widths, which is used for table rendering.
func (r *PdfRenderer) GetColumnWidths(n ast.Node) []float64 {
	return r.ColumnWidths[n]
}

// SetTOCLinks stores the heading→linkID map used by processText to place
// link anchors when each heading is rendered. Called by TOCDecorator.
func (r *PdfRenderer) SetTOCLinks(tocHeaders map[string]*int) {
	r.tocLinks = tocHeaders
}

// SetNodeProcessor allows users to set a custom node processor for a given node type.
func (r *PdfRenderer) SetNodeProcessor(nodeType string, processor node.Processor) error {
	if _, ok := r.nodeProcessors[nodeType]; ok {
		r.nodeProcessors[nodeType] = processor

		return nil
	}

	return fmt.Errorf("node type %s not found in nodeProcessors", nodeType)
}

// ---------------------------------------------------------------------------
// Processing
// ---------------------------------------------------------------------------

// Process takes the markdown content, parses it to generate the PDF.
func (r *PdfRenderer) Process(reader io.Reader) error {
	// try to open tracer
	var f *os.File
	var err error

	if r.tracerFile != "" {
		f, err = os.Create(r.tracerFile)
		if err != nil {
			return fmt.Errorf("os.Create() on tracefile error: %w", err)
		}
		defer func() { _ = f.Close() }()

		r.w = bufio.NewWriter(f)
		defer func() { _ = r.w.Flush() }()
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll() error: %w", err)
	}

	for _, pp := range r.preProcessors {
		if err := pp(content); err != nil {
			return fmt.Errorf("pre-processor error: %w", err)
		}
	}

	err = r.Run(content)
	if err != nil {
		return err
	}

	return nil
}

// Output writes the generated PDF to the provided io.Writer.
func (r *PdfRenderer) Output(w io.Writer) error {
	return r.Pdf.Output(w)
}

// OutputAndClose writes the generated PDF to the provided io.WriteCloser
// and closes it.
func (r *PdfRenderer) OutputAndClose(w io.WriteCloser) error {
	return r.Pdf.OutputAndClose(w)
}

// OutputFileAndClose writes the generated PDF to a file and closes it.
func (r *PdfRenderer) OutputFileAndClose(filename string) error {
	return r.Pdf.OutputFileAndClose(filename)
}

// Run takes the markdown content as a byte slice, preprocesses it, parses it into an AST,
// and renders it to PDF using the registered node processors.
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

// Parses all tables and sets the column width to the longest string in that column.
func setColumnWidths(doc ast.Node, r *PdfRenderer) {
	const cellPaddingPct = 0.2

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
				textlength += textlength * cellPaddingPct

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

// UpdateParagraphStyler - update with default styler.
func (r *PdfRenderer) UpdateParagraphStyler(defaultStyler theme.Styler) {
	initcurrent := &node.ContainerState{
		ListKind:  node.NotList,
		TextStyle: defaultStyler, LeftMargin: r.mleft,
	}
	r.cs.push(initcurrent)
}

// ---------------------------------------------------------------------------
// AST walker callbacks (gomarkdown Renderer interface)
// ---------------------------------------------------------------------------

// RenderNode dispatches each AST node to the registered NodeProcessor.
func (r *PdfRenderer) RenderNode(_ io.Writer, n ast.Node, entering bool) ast.WalkStatus {
	var key string

	switch n.(type) {
	case *ast.Text:
		key = "Text"
	case *ast.Softbreak:
		r.Tracer("Softbreak", "Output newline")
		r.Cr()

		return ast.GoToNext
	case *ast.Hardbreak:
		r.Tracer("Hardbreak", "Output newline")
		r.Cr()

		return ast.GoToNext
	case *ast.Emph:
		key = "Emph"
	case *ast.Strong:
		key = "Strong"
	case *ast.Del:
		if entering {
			r.Tracer("DEL (entering)", "Not handled")
		} else {
			r.Tracer("DEL (leaving)", "Not handled")
		}

		return ast.GoToNext
	case *ast.HTMLSpan:
		r.Tracer("HTMLSpan", "Not handled")

		return ast.GoToNext
	case *ast.Link:
		key = "Link"
	case *ast.Image:
		key = "Image"
	case *ast.Code:
		key = "Code"
	case *ast.Document:
		r.Tracer("Document", "Not Handled")

		return ast.GoToNext
	case *ast.Paragraph:
		key = "Paragraph"
	case *ast.BlockQuote:
		key = "BlockQuote"
	case *ast.HTMLBlock:
		key = "HTMLBlock"
	case *ast.Heading:
		key = "Heading"
	case *ast.HorizontalRule:
		key = "HorizontalRule"
	case *ast.List:
		key = "List"
	case *ast.ListItem:
		key = "ListItem"
	case *ast.CodeBlock:
		key = "CodeBlock"
	case *ast.Table:
		key = "Table"
	case *ast.TableHeader:
		key = "TableHeader"
	case *ast.TableBody:
		key = "TableBody"
	case *ast.TableRow:
		key = "TableRow"
	case *ast.TableCell:
		key = "TableCell"
	default:
		fmt.Printf("Unknown node type: %T. Skipping\n", n)

		return ast.GoToNext
	}

	if proc, ok := r.nodeProcessors[key]; ok {
		proc(r, n, entering)
	}

	return ast.GoToNext
}

// RenderHeader is not supported.
func (r *PdfRenderer) RenderHeader(_ io.Writer, _ ast.Node) {
	r.Tracer("RenderHeader", "Not handled")
}

// RenderFooter is not supported.
func (r *PdfRenderer) RenderFooter(_ io.Writer, _ ast.Node) {
	r.Tracer("RenderFooter", "Not handled")
}

func dorect(doc *fpdf.Fpdf, x, y, w, h float64, color colors.Color) {
	doc.SetFillColor(color.Red, color.Green, color.Blue)
	doc.Rect(x, y, w, h, "F")
}
