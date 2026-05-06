package renderer

// IsHorizontalRuleNewPage if true, will start a new page when encountering a HR (---). Useful for presentations.
func IsHorizontalRuleNewPage(value bool) RenderOption {
	return func(r *PdfRenderer) {
		r.HorizontalRuleNewPage = value
	}
}
