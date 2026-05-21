package renderer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type FooterHorizontalAlignment rune

const (
	FooterHorizontalAlignmentLeft   FooterHorizontalAlignment = 'L'
	FooterHorizontalAlignmentCenter FooterHorizontalAlignment = 'C'
	FooterHorizontalAlignmentRight  FooterHorizontalAlignment = 'R'
)

type FooterVerticalAlignment rune

const (
	FooterVerticalAlignmentTop      FooterVerticalAlignment = 'T'
	FooterVerticalAlignmentMiddle   FooterVerticalAlignment = 'M'
	FooterVerticalAlignmentBottom   FooterVerticalAlignment = 'B'
	FooterVerticalAlignmentBaseline FooterVerticalAlignment = 'A'
)

type FooterFontStyle rune

const (
	FooterFontStyleBold          FooterFontStyle = 'B'
	FooterFontStyleItalic        FooterFontStyle = 'I'
	FooterFontStyleUnderline     FooterFontStyle = 'U'
	FooterFontStyleStrikethrough FooterFontStyle = 'S'
)

type Footer struct {
	BackgroundColor *colors.Color
	Height          float64
	Left            []FooterSection
	Center          []FooterSection
	Right           []FooterSection
}

type FooterSection struct {
	Width     float64
	Height    float64
	RelativeX float64
	RelativeY float64

	BackgroundImage string
	Text            *FooterSectionTextContent
}

type FooterSectionTextContent struct {
	Text      string
	FontSize  float64
	FontStyle []FooterFontStyle
	Color     colors.Color

	HorizontalAlignment []FooterHorizontalAlignment
	VerticalAlignment   []FooterVerticalAlignment
}

func (c FooterSectionTextContent) GetFontStyleString() string {
	var result strings.Builder
	for _, style := range c.FontStyle {
		result.WriteString(string(style))
	}

	return result.String()
}

func (c FooterSectionTextContent) GetAlignmentString() string {
	var result strings.Builder
	for _, h := range c.HorizontalAlignment {
		result.WriteRune(rune(h))
	}

	for _, v := range c.VerticalAlignment {
		result.WriteRune(rune(v))
	}

	return result.String()
}

func WithFooter(footer Footer) RenderOption {
	return func(r *PdfRenderer) {
		r.Pdf.SetFooterFunc(func() {
			if footer.BackgroundColor != nil {
				r.Pdf.SetFillColor(
					footer.BackgroundColor.Red,
					footer.BackgroundColor.Green,
					footer.BackgroundColor.Blue,
				)
			} else {
				r.Pdf.SetFillColor(
					r.Theme.BackgroundColor.Red,
					r.Theme.BackgroundColor.Green,
					r.Theme.BackgroundColor.Blue,
				)
			}

			for _, section := range footer.Left {
				r.Pdf.SetXY(section.RelativeX, -footer.Height+section.RelativeY)
				drawFooterSection(r, section)
			}

			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())

			centerX := w / 2
			if r.orientation == "landscape" {
				centerX = h / 2
			}

			for _, section := range footer.Center {
				width := footerSectionWidth(r, section)
				r.Pdf.SetXY(centerX-width/2+section.RelativeX, -footer.Height+section.RelativeY)
				drawFooterSection(r, section)
			}

			for _, section := range footer.Right {
				width := footerSectionWidth(r, section)
				r.Pdf.SetXY(w-section.RelativeX-width, -footer.Height+section.RelativeY)
				drawFooterSection(r, section)
			}
		})
	}
}

func WithHeader(header Footer) RenderOption {
	return func(r *PdfRenderer) {
		r.Pdf.SetHeaderFuncMode(func() {
			if header.BackgroundColor != nil {
				r.Pdf.SetFillColor(
					header.BackgroundColor.Red,
					header.BackgroundColor.Green,
					header.BackgroundColor.Blue,
				)
			} else {
				r.Pdf.SetFillColor(
					r.Theme.BackgroundColor.Red,
					r.Theme.BackgroundColor.Green,
					r.Theme.BackgroundColor.Blue,
				)
			}

			for _, section := range header.Left {
				r.Pdf.SetXY(section.RelativeX, section.RelativeY)
				drawFooterSection(r, section)
			}

			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())

			centerX := w / 2
			if r.orientation == "landscape" {
				centerX = h / 2
			}

			for _, section := range header.Center {
				width := footerSectionWidth(r, section)
				r.Pdf.SetXY(centerX-width/2+section.RelativeX, section.RelativeY)
				drawFooterSection(r, section)
			}

			for _, section := range header.Right {
				width := footerSectionWidth(r, section)
				r.Pdf.SetXY(w-section.RelativeX-width, section.RelativeY)
				drawFooterSection(r, section)
			}
		}, false)
	}
}

func footerSectionWidth(r *PdfRenderer, section FooterSection) float64 {
	if section.Width != 0 {
		return section.Width
	}

	if section.Text != nil {
		r.Pdf.SetFont(
			r.Theme.Normal.Font,
			section.Text.GetFontStyleString(),
			section.Text.FontSize,
		)

		return r.Pdf.GetStringWidth(resolveTextPlaceholders(r, section.Text.Text))
	}

	return 0
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

func drawFooterSection(r *PdfRenderer, section FooterSection) {
	if section.BackgroundImage != "" {
		fmt.Println("TODO: add background image to footer section")
	}

	if section.Text != nil {
		r.Pdf.SetFont(
			r.Theme.Normal.Font,
			section.Text.GetFontStyleString(),
			section.Text.FontSize,
		)

		r.Pdf.SetTextColor(
			section.Text.Color.Red,
			section.Text.Color.Green,
			section.Text.Color.Blue,
		)

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

// WithDefaultFooter configures the PDF renderer to include a default footer on each page.
func WithDefaultFooter() RenderOption {
	return WithFooter(Footer{
		Height: 15,
		Left: []FooterSection{
			{
				Width:     0,
				Height:    10,
				RelativeX: 10,
				Text: &FooterSectionTextContent{
					Text:                "%AUTHOR%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleItalic},
					Color:               colors.FromRGB(255, 50, 50),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentLeft},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
		Center: []FooterSection{
			{
				Width:  10,
				Height: 10,
				Text: &FooterSectionTextContent{
					Text:                "%TITLE%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleBold},
					Color:               colors.FromRGB(50, 50, 255),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentCenter},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
		Right: []FooterSection{
			{
				Width:     0,
				Height:    10,
				RelativeX: 10,
				Text: &FooterSectionTextContent{
					Text:                "Page %PAGE_NUMBER%/%PAGE_TOTAL%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleItalic},
					Color:               colors.FromRGB(50, 255, 50),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentRight},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
	})
}

// WithDefaultHeader configures the PDF renderer to include a default header on each page.
func WithDefaultHeader() RenderOption {
	return WithHeader(Footer{
		Height: 40,
		Left: []FooterSection{
			{
				Width:     0,
				Height:    10,
				RelativeX: 10,
				Text: &FooterSectionTextContent{
					Text:                "%AUTHOR%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleItalic},
					Color:               colors.FromRGB(255, 50, 50),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentLeft},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
		Center: []FooterSection{
			{
				Width:  10,
				Height: 10,
				Text: &FooterSectionTextContent{
					Text:                "%TITLE%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleBold},
					Color:               colors.FromRGB(50, 50, 255),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentCenter},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
		Right: []FooterSection{
			{
				Width:     0,
				Height:    10,
				RelativeX: 10,
				Text: &FooterSectionTextContent{
					Text:                "Page %PAGE_NUMBER%/%PAGE_TOTAL%",
					FontSize:            12,
					FontStyle:           []FooterFontStyle{FooterFontStyleItalic},
					Color:               colors.FromRGB(50, 255, 50),
					HorizontalAlignment: []FooterHorizontalAlignment{FooterHorizontalAlignmentRight},
					VerticalAlignment:   []FooterVerticalAlignment{FooterVerticalAlignmentMiddle},
				},
			},
		},
	})
}
