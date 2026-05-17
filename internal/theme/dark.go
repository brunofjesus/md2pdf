package theme

import "github.com/brunofjesus/md2pdf/v3/internal/colors"

const (
	darkThemeFont                 = "LiberationSans"
	darkThemeMonoFont             = "LiberationMono"
	darkThemeFontSize             = 12
	darkThemeHorizontalRuleHeight = 3
)

// DarkTheme returns a Theme with colors and styles inspired by GitHub's dark mode.
func DarkTheme() *Theme {
	// Regular colors
	ghText := colors.FromRGB(230, 237, 243)  // #e6edf3
	ghMuted := colors.FromRGB(125, 133, 144) // #7d8590
	ghLink := colors.FromRGB(74, 158, 255)   // #4a9eff
	ghCodeBg := colors.FromRGB(22, 27, 34)   // #161b22
	ghInlineBg := colors.FromRGB(52, 59, 66) // #343b42
	ghBg := colors.FromRGB(13, 17, 23)       // #0d1117

	// Codeblock syntax highlight colors (GitHub Dark)
	ghRedSyn := colors.FromRGB(255, 123, 114)  // #ff7b72
	ghBlueSyn := colors.FromRGB(121, 192, 255) // #79c0ff
	ghStrCyan := colors.FromRGB(165, 214, 255) // #a5d6ff
	ghOrange := colors.FromRGB(255, 166, 87)   // #ffa657
	ghGray := colors.FromRGB(139, 148, 158)    // #8b949e

	//nolint:dupl
	r := Theme{
		BackgroundColor: ghBg,
		Normal: Styler{
			Font: darkThemeFont, Style: "", Size: darkThemeFontSize, Spacing: 2,
			TextColor: ghText, FillColor: ghBg,
		},
		Link: Styler{
			Font: darkThemeFont, Style: "", Size: darkThemeFontSize, Spacing: 2,
			TextColor: ghLink, FillColor: ghBg,
		},
		Backtick: Styler{
			Font: darkThemeMonoFont, Style: "", Size: darkThemeFontSize, Spacing: 2,
			TextColor: ghText, FillColor: ghInlineBg,
		},
		Code: Code{
			Text: Styler{
				Font: darkThemeMonoFont, Style: "", Size: darkThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghCodeBg,
			},
			TabWidth: 4,
			Colors: map[string]colors.Color{
				"default":              ghText,
				"statement":            ghRedSyn,
				"green":                ghRedSyn,
				"identifier":           ghBlueSyn,
				"blue":                 ghBlueSyn,
				"preproc":              ghRedSyn,
				"special":              ghRedSyn,
				"type.keyword":         ghRedSyn,
				"red":                  ghRedSyn,
				"constant":             ghBlueSyn,
				"constant.number":      ghBlueSyn,
				"constant.bool":        ghBlueSyn,
				"symbol.brackets":      ghText,
				"identifier.var":       ghBlueSyn,
				"cyan":                 ghBlueSyn,
				"constant.specialChar": ghStrCyan,
				"constant.string.url":  ghStrCyan,
				"constant.string":      ghStrCyan,
				"magenta":              ghStrCyan,
				"type":                 ghOrange,
				"symbol":               ghOrange,
				"symbol.operator":      ghText,
				"symbol.tag.extended":  ghOrange,
				"yellow":               ghOrange,
				"comment":              ghGray,
				"high.green":           ghGray,
			},
		},
		Heading: Heading{
			H1: Styler{
				Font: darkThemeFont, Style: "b", Size: 32, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H2: Styler{
				Font: darkThemeFont, Style: "b", Size: 24, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H3: Styler{
				Font: darkThemeFont, Style: "b", Size: 20, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H4: Styler{
				Font: darkThemeFont, Style: "b", Size: 16, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H5: Styler{
				Font: darkThemeFont, Style: "b", Size: 14, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H6: Styler{
				Font: darkThemeFont, Style: "b", Size: 13.6, Spacing: 5,
				TextColor: ghMuted, FillColor: ghBg,
			},
			Line: HorizontalRule{
				Height: 0.75,
				Color:  colors.FromRGB(61, 68, 77),
			},
		},
		Blockquote: Styler{
			Font: darkThemeFont, Style: "i", Size: darkThemeFontSize, Spacing: 2,
			TextColor: ghMuted, FillColor: ghBg,
		},
		Table: Table{
			Header: Styler{
				Font: darkThemeFont, Style: "b", Size: darkThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghCodeBg,
			},
			Body: Styler{
				Font: darkThemeFont, Style: "", Size: darkThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghBg,
			},
		},
		HorizontalRule: HorizontalRule{
			Height: darkThemeHorizontalRuleHeight,
			Color:  colors.FromRGB(61, 68, 77),
		},
		IndentValue: 2,
	}

	return &r
}
