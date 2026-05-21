package renderer

import (
	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type Footer struct {
	BackgroundColor *colors.Color
	Height          float64
	Left            []MarginalSection
	Center          []MarginalSection
	Right           []MarginalSection
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
				drawMarginalSection(r, section)
			}

			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())

			centerX := w / 2
			if r.orientation == "landscape" {
				centerX = h / 2
			}

			for _, section := range footer.Center {
				width := marginalSectionWidth(r, section)
				r.Pdf.SetXY(centerX-width/2+section.RelativeX, -footer.Height+section.RelativeY)
				drawMarginalSection(r, section)
			}

			for _, section := range footer.Right {
				width := marginalSectionWidth(r, section)
				r.Pdf.SetXY(w-section.RelativeX-width, -footer.Height+section.RelativeY)
				drawMarginalSection(r, section)
			}
		})
	}
}

// WithDefaultFooter configures the PDF renderer to include a default footer on each page.
func WithDefaultFooter() RenderOption {
	gray := colors.FromRGB(128, 128, 128)

	return WithFooter(Footer{
		Height: 12,
		Center: []MarginalSection{
			{
				Height: 8,
				Text: &MarginalSectionTextContent{
					Text:                "%PAGE_NUMBER% / %PAGE_TOTAL%",
					FontSize:            9,
					FontStyle:           []MarginalFontStyle{},
					Color:               &gray,
					HorizontalAlignment: []MarginalHorizontalAlignment{MarginalHorizontalAlignmentCenter},
					VerticalAlignment:   []MarginalVerticalAlignment{MarginalVerticalAlignmentMiddle},
				},
			},
		},
	})
}
