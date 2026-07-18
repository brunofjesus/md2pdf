package renderer

import (
	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

// WithFooter configures the PDF renderer to include a custom footer on each page.
func WithFooter(footer Marginal) RenderOption {
	return func(r *PdfRenderer) {
		r.Pdf.SetFooterFunc(func() {
			w, h, _ := r.Pdf.PageSize(r.Pdf.PageNo())

			if footer.BackgroundColor != nil {
				r.Pdf.SetFillColor(
					footer.BackgroundColor.Red,
					footer.BackgroundColor.Green,
					footer.BackgroundColor.Blue,
				)

				r.Pdf.SetDrawColor(footer.BackgroundColor.Red, footer.BackgroundColor.Green, footer.BackgroundColor.Blue)
				r.Pdf.Rect(0, h-footer.Height, w, footer.Height, "DF")
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

	return WithFooter(Marginal{
		BackgroundColor: nil,
		Left:            nil,
		Right:           nil,
		Height:          12,
		Center: []MarginalSection{
			{
				Width:                   0,
				Height:                  8,
				RelativeX:               0,
				RelativeY:               0,
				BackgroundImage:         "",
				resolvedBackgroundImage: "",
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
