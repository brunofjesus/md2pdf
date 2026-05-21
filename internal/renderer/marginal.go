package renderer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type MarginalHorizontalAlignment rune

const (
	MarginalHorizontalAlignmentLeft   MarginalHorizontalAlignment = 'L'
	MarginalHorizontalAlignmentCenter MarginalHorizontalAlignment = 'C'
	MarginalHorizontalAlignmentRight  MarginalHorizontalAlignment = 'R'
)

type MarginalVerticalAlignment rune

const (
	MarginalVerticalAlignmentTop      MarginalVerticalAlignment = 'T'
	MarginalVerticalAlignmentMiddle   MarginalVerticalAlignment = 'M'
	MarginalVerticalAlignmentBottom   MarginalVerticalAlignment = 'B'
	MarginalVerticalAlignmentBaseline MarginalVerticalAlignment = 'A'
)

type MarginalFontStyle rune

const (
	MarginalFontStyleBold          MarginalFontStyle = 'B'
	MarginalFontStyleItalic        MarginalFontStyle = 'I'
	MarginalFontStyleUnderline     MarginalFontStyle = 'U'
	MarginalFontStyleStrikethrough MarginalFontStyle = 'S'
)

type MarginalSection struct {
	Width     float64
	Height    float64
	RelativeX float64
	RelativeY float64

	BackgroundImage string
	Text            *MarginalSectionTextContent
}

type MarginalSectionTextContent struct {
	Text      string
	FontSize  float64
	FontStyle []MarginalFontStyle
	Color     *colors.Color

	HorizontalAlignment []MarginalHorizontalAlignment
	VerticalAlignment   []MarginalVerticalAlignment
}

func (c MarginalSectionTextContent) GetFontStyleString() string {
	var result strings.Builder
	for _, style := range c.FontStyle {
		result.WriteString(string(style))
	}

	return result.String()
}

func (c MarginalSectionTextContent) GetAlignmentString() string {
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
		"%PAGE_TOTAL%", strconv.Itoa(r.Pdf.PageCount()),
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
			section.Text.GetFontStyleString(),
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
			section.Text.GetFontStyleString(),
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
			section.Text.GetAlignmentString(),
			true,
			0,
			"",
		)
	}
}
