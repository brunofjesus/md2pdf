package renderer

import (
	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type Header struct {
	BackgroundColor *colors.Color
	Height          float64
	Left            []MarginalSection
	Center          []MarginalSection
	Right           []MarginalSection
}

func WithHeader(header Header) RenderOption {
	return func(r *PdfRenderer) {
		r.Pdf.SetTopMargin(header.Height)
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
				drawMarginalSection(r, section)
			}

			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())

			centerX := w / 2
			if r.orientation == "landscape" {
				centerX = h / 2
			}

			for _, section := range header.Center {
				width := marginalSectionWidth(r, section)
				r.Pdf.SetXY(centerX-width/2+section.RelativeX, section.RelativeY)
				drawMarginalSection(r, section)
			}

			for _, section := range header.Right {
				width := marginalSectionWidth(r, section)
				r.Pdf.SetXY(w-section.RelativeX-width, section.RelativeY)
				drawMarginalSection(r, section)
			}
		}, true)
	}
}

// WithDefaultHeader configures the PDF renderer to include a default header on each page.
func WithDefaultHeader() RenderOption {
	darkGray := colors.FromRGB(51, 51, 51)
	medGray := colors.FromRGB(102, 102, 102)

	return WithHeader(Header{
		Height: 20,
		Left: []MarginalSection{
			{
				Height:    10,
				RelativeX: 10,
				RelativeY: 5,
				Text: &MarginalSectionTextContent{
					Text:                "%TITLE%",
					FontSize:            10,
					FontStyle:           []MarginalFontStyle{MarginalFontStyleBold},
					Color:               &darkGray,
					HorizontalAlignment: []MarginalHorizontalAlignment{MarginalHorizontalAlignmentLeft},
					VerticalAlignment:   []MarginalVerticalAlignment{MarginalVerticalAlignmentMiddle},
				},
			},
		},
		Right: []MarginalSection{
			{
				Height:    10,
				RelativeX: 10,
				RelativeY: 5,
				Text: &MarginalSectionTextContent{
					Text:                "%AUTHOR%",
					FontSize:            9,
					FontStyle:           []MarginalFontStyle{MarginalFontStyleItalic},
					Color:               &medGray,
					HorizontalAlignment: []MarginalHorizontalAlignment{MarginalHorizontalAlignmentRight},
					VerticalAlignment:   []MarginalVerticalAlignment{MarginalVerticalAlignmentMiddle},
				},
			},
		},
	})
}
