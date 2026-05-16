// Package md2pdf converts Markdown documents to PDF files.
//
// Create a converter with [New], feed it Markdown content via one of the
// Process methods ([Md2Pdf.Process], [Md2Pdf.ProcessAndClose], or
// [Md2Pdf.ProcessFileAndClose]), and write the result with one of the Output
// methods ([Md2Pdf.Output], [Md2Pdf.OutputAndClose], or
// [Md2Pdf.OutputFileAndClose]).
//
// Functional options such as [WithTableOfContents], [WithHorizontalRuleAsNewPage],
// [WithBaseURL], and [WithDefaultFooter] can be passed via [Md2PdfParams.Options]
// to customize the conversion.
package md2pdf

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/brunofjesus/md2pdf/v3/internal/renderer"
)

// Md2Pdf converts Markdown content to PDF. Create one with [New], feed it
// Markdown via one of the Process methods, then write the result with one of
// the Output methods.
type Md2Pdf struct {
	pdfRenderer *renderer.PdfRenderer
	params      Md2PdfParams
}

type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// Md2PdfParams holds the configuration for creating a new [Md2Pdf] instance.
type Md2PdfParams struct {
	// Title is the PDF document title stored in the file metadata.
	Title string
	// Orientation is the page orientation: "portrait" (default) or "landscape".
	Orientation Orientation
	// PageSize is the page size, e.g. "A4", "Letter". Defaults to "A4".
	PageSize string
	// Options is a list of functional options to configure the renderer.
	Options []Option
	// Theme selects the color theme: "light", "dark", or a path to a custom
	// theme JSON file.
	Theme string
}

// Option is a functional option that configures an [Md2Pdf] instance.
type Option func(r *Md2Pdf)

// New creates a new [Md2Pdf] converter with the given parameters.
func New(p Md2PdfParams) (*Md2Pdf, error) {
	switch p.Orientation {
	case "":
		p.Orientation = OrientationPortrait
	case OrientationPortrait, OrientationLandscape:
	default:
		return nil, fmt.Errorf("invalid orientation: %s", p.Orientation)
	}

	var theme renderer.Theme
	var customThemeFile string
	switch p.Theme {
	case "light":
		theme = renderer.LIGHT
	case "dark":
		theme = renderer.DARK
	default:
		theme = renderer.CUSTOM
		customThemeFile = p.Theme
	}

	rend := renderer.NewPdfRenderer(renderer.PdfRendererParams{
		Title:           p.Title,
		Orientation:     string(p.Orientation),
		PageSize:        p.PageSize,
		Theme:           theme,
		CustomThemeFile: customThemeFile,
	})

	md2pdf := &Md2Pdf{
		pdfRenderer: rend,
	}

	for _, opt := range p.Options {
		opt(md2pdf)
	}

	return md2pdf, nil
}

// Process reads Markdown from reader and renders it into an in-memory PDF.
// Call one of the Output methods afterwards to retrieve the result.
func (m *Md2Pdf) Process(reader io.Reader) error {
	return m.pdfRenderer.Process(reader)
}

// ProcessAndClose reads Markdown from reader, renders it into an in-memory PDF,
// and closes the reader. Both processing and close errors are reported via
// [errors.Join]. Call one of the Output methods afterwards to retrieve the result.
func (m *Md2Pdf) ProcessAndClose(reader io.ReadCloser) error {
	errs := make([]error, 0, 2)

	if err := m.pdfRenderer.Process(reader); err != nil {
		errs = append(errs, fmt.Errorf("failed to process markdown input: %w", err))
	}

	if err := reader.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close markdown input reader: %w", err))
	}

	return errors.Join(errs...)
}

// ProcessFileAndClose opens the given Markdown file, renders it into an
// in-memory PDF, and closes the file. Call one of the Output methods afterwards
// to retrieve the result.
func (m *Md2Pdf) ProcessFileAndClose(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	errs := make([]error, 0, 2)
	if err := m.pdfRenderer.Process(file); err != nil {
		errs = append(errs, fmt.Errorf("failed to process markdown input: %w", err))
	}

	if err := file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close markdown input file: %w", err))
	}

	return errors.Join(errs...)
}

// Output writes the generated PDF to the provided io.Writer.
func (m *Md2Pdf) Output(w io.Writer) error {
	return m.pdfRenderer.Output(w)
}

// OutputAndClose writes the generated PDF to the provided io.WriteCloser
// and closes it.
func (m *Md2Pdf) OutputAndClose(w io.WriteCloser) error {
	return m.pdfRenderer.OutputAndClose(w)
}

// OutputFileAndClose writes the generated PDF to a file and closes it.
func (m *Md2Pdf) OutputFileAndClose(filename string) error {
	return m.pdfRenderer.OutputFileAndClose(filename)
}

// WithTableOfContents configures the renderer to generate a table of contents
// based on the headings in the markdown input.
func WithTableOfContents() Option {
	return func(m *Md2Pdf) {
		renderer.WithTableOfContents()(m.pdfRenderer)
	}
}

// WithHorizontalRuleAsNewPage configures the renderer to interpret horizontal
// rules (---) as page breaks,
func WithHorizontalRuleAsNewPage() Option {
	return func(m *Md2Pdf) {
		renderer.WithHorizontalRuleAsNewPage()(m.pdfRenderer)
	}
}

// WithBaseURL sets the base URL for resolving relative links and images in the
// input markdown.
func WithBaseURL(baseURL string) Option {
	return func(m *Md2Pdf) {
		renderer.WithBaseURL(baseURL)(m.pdfRenderer)
	}
}

// WithDefaultFooter configures the renderer to add a default footer to each page of the PDF,
// containing the specified orientation, author, and title.
func WithDefaultFooter(author, title string) Option {
	return func(m *Md2Pdf) {
		renderer.WithDefaultFooter(string(m.params.Orientation), author, title)(m.pdfRenderer)
	}
}
