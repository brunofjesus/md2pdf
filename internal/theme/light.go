package theme

import "github.com/brunofjesus/md2pdf/v3/internal/colors"

const (
	lightThemeFont                 = "LiberationSans"
	lightThemeMonoFont             = "LiberationMono"
	lightThemeFontSize             = 12
	lightThemeHorizontalRuleHeight = 3
)

// LightTheme returns a Theme with colors and styles inspired by GitHub's light mode.
func LightTheme() *Theme {
	// Regular colors
	ghText := colors.FromRGB(31, 35, 40)        // #1f2328
	ghMuted := colors.FromRGB(101, 109, 118)    // #656d76
	ghLink := colors.FromRGB(9, 105, 218)       // #0969da
	ghCodeBg := colors.FromRGB(246, 248, 250)   // #f6f8fa
	ghInlineBg := colors.FromRGB(239, 241, 243) // #eff1f3
	ghBg := colors.FromRGB(255, 255, 255)

	// Codeblock syntax highlight colors (GitHub Light)
	ghPurple := colors.FromRGB(130, 80, 223)   // #8250df
	ghDkBlue := colors.FromRGB(5, 80, 174)     // #0550ae
	ghRed := colors.FromRGB(207, 34, 46)       // #cf222e
	ghStrBlue := colors.FromRGB(10, 48, 105)   // #0a3069
	ghBrown := colors.FromRGB(149, 56, 0)      // #953800
	ghComment := colors.FromRGB(110, 119, 129) // #6e7781

	//nolint:dupl
	r := Theme{
		BackgroundColor: ghBg,
		Normal: Styler{
			Font: lightThemeFont, Style: "", Size: lightThemeFontSize, Spacing: 2,
			TextColor: ghText, FillColor: ghBg,
		},
		Link: Styler{
			Font: lightThemeFont, Style: "", Size: lightThemeFontSize, Spacing: 2,
			TextColor: ghLink, FillColor: ghBg,
		},
		Backtick: Styler{
			Font: lightThemeMonoFont, Style: "", Size: lightThemeFontSize, Spacing: 2,
			TextColor: ghText, FillColor: ghInlineBg,
		},
		Code: Code{
			Text: Styler{
				Font: lightThemeMonoFont, Style: "", Size: lightThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghCodeBg,
			},
			TabWidth: 4,
			Colors: map[string]colors.Color{
				"default":              ghText,
				"statement":            ghPurple,
				"green":                ghPurple,
				"identifier":           ghDkBlue,
				"blue":                 ghDkBlue,
				"preproc":              ghRed,
				"special":              ghRed,
				"type.keyword":         ghRed,
				"red":                  ghRed,
				"constant":             ghDkBlue,
				"constant.number":      ghDkBlue,
				"constant.bool":        ghDkBlue,
				"symbol.brackets":      ghText,
				"identifier.var":       ghDkBlue,
				"cyan":                 ghDkBlue,
				"constant.specialChar": ghStrBlue,
				"constant.string.url":  ghStrBlue,
				"constant.string":      ghStrBlue,
				"magenta":              ghStrBlue,
				"type":                 ghBrown,
				"symbol":               ghBrown,
				"symbol.operator":      ghText,
				"symbol.tag.extended":  ghDkBlue,
				"yellow":               ghBrown,
				"comment":              ghComment,
				"high.green":           ghComment,
			},
		},
		Heading: Heading{
			H1: Styler{
				Font: lightThemeFont, Style: "b", Size: 32, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H2: Styler{
				Font: lightThemeFont, Style: "b", Size: 24, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H3: Styler{
				Font: lightThemeFont, Style: "b", Size: 20, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H4: Styler{
				Font: lightThemeFont, Style: "b", Size: 16, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H5: Styler{
				Font: lightThemeFont, Style: "b", Size: 14, Spacing: 5,
				TextColor: ghText, FillColor: ghBg,
			},
			H6: Styler{
				Font: lightThemeFont, Style: "b", Size: 13.6, Spacing: 5,
				TextColor: ghMuted, FillColor: ghBg,
			},
			Line: HorizontalRule{
				Height: 0.75,
				Color:  colors.FromRGB(209, 217, 224),
			},
		},
		Blockquote: Styler{
			Font: lightThemeFont, Style: "i", Size: lightThemeFontSize, Spacing: 2,
			TextColor: ghMuted, FillColor: ghBg,
		},
		Table: Table{
			Header: Styler{
				Font: lightThemeFont, Style: "b", Size: lightThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghCodeBg,
			},
			Body: Styler{
				Font: lightThemeFont, Style: "", Size: lightThemeFontSize, Spacing: 2,
				TextColor: ghText, FillColor: ghBg,
			},
		},
		HorizontalRule: HorizontalRule{
			Height: lightThemeHorizontalRuleHeight,
			Color:  colors.FromRGB(209, 217, 224),
		},
		IndentValue: 2,
	}

	return &r
}
