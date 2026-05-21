package renderer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

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

// MarginalSection defines a content block within a header or footer, with positioning,
// dimensions, and optional text or background image content.
type MarginalSection struct {
	Width     float64 `json:"width,omitempty"`
	Height    float64 `json:"height"`
	RelativeX float64 `json:"relativeX,omitempty"`
	RelativeY float64 `json:"relativeY,omitempty"`

	BackgroundImage string                      `json:"backgroundImage,omitempty"`
	Text            *MarginalSectionTextContent `json:"text,omitempty"`
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
	if section.BackgroundImage != "" {
		fmt.Println("TODO: add background image to marginal section")
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
			true,
			0,
			"",
		)
	}
}
