package renderer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"codeberg.org/go-pdf/fpdf"
	"github.com/brunofjesus/md2pdf/v3/internal/colors"
	"github.com/brunofjesus/md2pdf/v3/internal/renderer/node"
)

// Marginal is the struct the defines either an Header or a Footer.
// Can by passed to a function like [WithHeader] or [WithFooter] in order
// for it to be included on the document.
type Marginal struct {
	BackgroundColor *colors.Color     `json:"backgroundColor,omitempty"`
	Height          float64           `json:"height"`
	Left            []MarginalSection `json:"left,omitempty"`
	Center          []MarginalSection `json:"center,omitempty"`
	Right           []MarginalSection `json:"right,omitempty"`
}

// FromJSONFile reads a JSON file and unmarshals its content into the Marginal struct.
func (m *Marginal) FromJSONFile(file string) error {
	f, err := os.Open(file) //nolint:gosec // G304: file path is an intentional user-supplied CLI flag
	if err != nil {
		return fmt.Errorf("opening marginal file %q: %w", file, err)
	}
	defer func() {
		_ = f.Close()
	}()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var parsed Marginal
	if err := dec.Decode(&parsed); err != nil {
		return fmt.Errorf("parsing marginal file %q: %w", file, err)
	}

	*m = parsed

	return nil
}

// MarginalHorizontalAlignment represents the horizontal text alignment within a marginal section.
type MarginalHorizontalAlignment rune

const (
	// MarginalHorizontalAlignmentLeft aligns text to the left edge of the section.
	MarginalHorizontalAlignmentLeft MarginalHorizontalAlignment = 'L'
	// MarginalHorizontalAlignmentCenter centers text horizontally within the section.
	MarginalHorizontalAlignmentCenter MarginalHorizontalAlignment = 'C'
	// MarginalHorizontalAlignmentRight aligns text to the right edge of the section.
	MarginalHorizontalAlignmentRight MarginalHorizontalAlignment = 'R'
)

// MarshalJSON encodes the alignment as a single-character string (e.g. "L").
func (a MarginalHorizontalAlignment) MarshalJSON() ([]byte, error) {
	return marshalRuneChar(rune(a))
}

// UnmarshalJSON decodes a single-character string (e.g. "L") into the alignment.
func (a *MarginalHorizontalAlignment) UnmarshalJSON(data []byte) error {
	r, err := unmarshalRuneChar(data)
	if err != nil {
		return err
	}

	*a = MarginalHorizontalAlignment(r)

	return nil
}

// MarginalVerticalAlignment represents the vertical text alignment within a marginal section.
type MarginalVerticalAlignment rune

const (
	// MarginalVerticalAlignmentTop aligns text to the top edge of the section.
	MarginalVerticalAlignmentTop MarginalVerticalAlignment = 'T'
	// MarginalVerticalAlignmentMiddle centers text vertically within the section.
	MarginalVerticalAlignmentMiddle MarginalVerticalAlignment = 'M'
	// MarginalVerticalAlignmentBottom aligns text to the bottom edge of the section.
	MarginalVerticalAlignmentBottom MarginalVerticalAlignment = 'B'
	// MarginalVerticalAlignmentBaseline aligns text to the baseline of the section.
	MarginalVerticalAlignmentBaseline MarginalVerticalAlignment = 'A'
)

// MarshalJSON encodes the alignment as a single-character string (e.g. "M").
func (a MarginalVerticalAlignment) MarshalJSON() ([]byte, error) {
	return marshalRuneChar(rune(a))
}

// UnmarshalJSON decodes a single-character string (e.g. "M") into the alignment.
func (a *MarginalVerticalAlignment) UnmarshalJSON(data []byte) error {
	r, err := unmarshalRuneChar(data)
	if err != nil {
		return err
	}

	*a = MarginalVerticalAlignment(r)

	return nil
}

// MarginalFontStyle represents a font style modifier (bold, italic, underline, strikethrough).
type MarginalFontStyle rune

const (
	// MarginalFontStyleBold applies bold styling to the text.
	MarginalFontStyleBold MarginalFontStyle = 'B'
	// MarginalFontStyleItalic applies italic styling to the text.
	MarginalFontStyleItalic MarginalFontStyle = 'I'
	// MarginalFontStyleUnderline applies underline styling to the text.
	MarginalFontStyleUnderline MarginalFontStyle = 'U'
	// MarginalFontStyleStrikethrough applies strikethrough styling to the text.
	MarginalFontStyleStrikethrough MarginalFontStyle = 'S'
)

// MarshalJSON encodes the font style as a single-character string (e.g. "B").
func (s MarginalFontStyle) MarshalJSON() ([]byte, error) {
	return marshalRuneChar(rune(s))
}

// UnmarshalJSON decodes a single-character string (e.g. "B") into the font style.
func (s *MarginalFontStyle) UnmarshalJSON(data []byte) error {
	r, err := unmarshalRuneChar(data)
	if err != nil {
		return err
	}

	*s = MarginalFontStyle(r)

	return nil
}

// MarginalSection defines a content block within a header or footer, with positioning,
// dimensions, and optional text or background image content.
type MarginalSection struct {
	Width     float64 `json:"width,omitempty"`
	Height    float64 `json:"height"`
	RelativeX float64 `json:"relativeX,omitempty"`
	RelativeY float64 `json:"relativeY,omitempty"`

	BackgroundImage string                      `json:"backgroundImage,omitempty"`
	Text            *MarginalSectionTextContent `json:"text,omitempty"`

	resolvedBackgroundImage string `json:"-"`
}

// ResolvedBackgroundImage resolves the background image path for the marginal section,
// downloading it if necessary, and caches the result for future use.
// It returns the resolved local file path or an error if resolution fails.
func (c *MarginalSection) ResolvedBackgroundImage(r *PdfRenderer) (string, error) {
	if c.resolvedBackgroundImage != "" {
		return c.resolvedBackgroundImage, nil
	} else if c.BackgroundImage != "" {
		path, err := node.ResolveImagePath(r, c.BackgroundImage)
		if err != nil {
			return "", err
		}

		c.resolvedBackgroundImage = path
	}

	return c.resolvedBackgroundImage, nil
}

// MarginalSectionTextContent defines the text content and styling for a marginal section,
// including font size, style, color, and alignment.
type MarginalSectionTextContent struct {
	Text      string              `json:"text"`
	FontSize  float64             `json:"fontSize,omitempty"`
	FontStyle []MarginalFontStyle `json:"fontStyle,omitempty"`
	Color     *colors.Color       `json:"color,omitempty"`

	HorizontalAlignment []MarginalHorizontalAlignment `json:"horizontalAlignment,omitempty"`
	VerticalAlignment   []MarginalVerticalAlignment   `json:"verticalAlignment,omitempty"`
}

func (c MarginalSectionTextContent) getFontStyleString() string {
	var result strings.Builder
	for _, style := range c.FontStyle {
		result.WriteString(string(style))
	}

	return result.String()
}

func (c MarginalSectionTextContent) getAlignmentString() string {
	var result strings.Builder
	for _, h := range c.HorizontalAlignment {
		result.WriteRune(rune(h))
	}

	for _, v := range c.VerticalAlignment {
		result.WriteRune(rune(v))
	}

	return result.String()
}

func resolveTextPlaceholders(r *PdfRenderer, text string) string {
	replacer := strings.NewReplacer(
		"%AUTHOR%", r.metadata[MetadataKeyAuthor],
		"%TITLE%", r.metadata[MetadataKeyTitle],
		"%PAGE_NUMBER%", strconv.Itoa(r.Pdf.PageNo()),
	)

	return replacer.Replace(text)
}

func marginalSectionWidth(r *PdfRenderer, section MarginalSection) float64 {
	if section.Width != 0 {
		return section.Width
	}

	if section.Text != nil {
		fontSize := section.Text.FontSize
		if fontSize == 0 {
			fontSize = r.Theme.Normal.Size
		}

		r.Pdf.SetFont(
			r.Theme.Normal.Font,
			section.Text.getFontStyleString(),
			fontSize,
		)

		return r.Pdf.GetStringWidth(resolveTextPlaceholders(r, section.Text.Text))
	}

	return 0
}

func drawMarginalSection(r *PdfRenderer, section MarginalSection) {
	backgroundImage, err := section.ResolvedBackgroundImage(r)
	if err == nil && backgroundImage != "" {
		x, y := r.Pdf.GetXY()
		r.Pdf.ImageOptions(
			backgroundImage, x, y, section.Width, section.Height, false,
			fpdf.ImageOptions{
				ImageType: "", ReadDpi: true, AllowNegativePosition: false,
			},
			0, "",
		)
	} else if err != nil {
		log.Printf("Error resolving background image for marginal section: %v", err)
	}

	if section.Text != nil {
		fontSize := section.Text.FontSize
		if fontSize == 0 {
			fontSize = r.Theme.Normal.Size
		}

		r.Pdf.SetFont(
			r.Theme.Normal.Font,
			section.Text.getFontStyleString(),
			fontSize,
		)

		if section.Text.Color != nil {
			r.Pdf.SetTextColor(
				section.Text.Color.Red,
				section.Text.Color.Green,
				section.Text.Color.Blue,
			)
		} else {
			r.Pdf.SetTextColor(
				r.Theme.Normal.TextColor.Red,
				r.Theme.Normal.TextColor.Green,
				r.Theme.Normal.TextColor.Blue,
			)
		}

		txtContent := resolveTextPlaceholders(r, section.Text.Text)

		width := section.Width
		if section.Width == 0 {
			width = r.Pdf.GetStringWidth(txtContent)
		}

		r.Pdf.CellFormat(
			width,
			section.Height,
			txtContent,
			"",
			1,
			section.Text.getAlignmentString(),
			backgroundImage == "",
			0,
			"",
		)
	}
}
