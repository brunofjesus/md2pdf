package theme

import "github.com/brunofjesus/md2pdf/v3/internal/colors"

func LightTheme() *Theme {
	r := Theme{}

	ghText := colors.New(31, 35, 40)        // #1f2328
	ghMuted := colors.New(101, 109, 118)    // #656d76
	ghLink := colors.New(9, 105, 218)       // #0969da
	ghCodeBg := colors.New(246, 248, 250)   // #f6f8fa
	ghInlineBg := colors.New(239, 241, 243) // #eff1f3
	ghWhite := colors.New(255, 255, 255)

	r.BackgroundColor = ghWhite

	// Normal Text
	r.Normal = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghWhite,
	}

	// Link text
	r.Link = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: ghLink,
	}

	// Backticked text
	r.Backtick = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghInlineBg,
	}

	// Code text
	// Codeblock syntax highlight colors (GitHub Light)
	ghPurple := colors.New(130, 80, 223)   // #8250df
	ghDkBlue := colors.New(5, 80, 174)     // #0550ae
	ghRed := colors.New(207, 34, 46)       // #cf222e
	ghStrBlue := colors.New(10, 48, 105)   // #0a3069
	ghBrown := colors.New(149, 56, 0)      // #953800
	ghComment := colors.New(110, 119, 129) // #6e7781

	r.Code = Code{
		Text: Styler{
			Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
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
	}

	// Headings
	r.Heading = Heading{
		H1: Styler{
			Font: "LiberationSans", Style: "b", Size: 32, Spacing: 5,
			TextColor: ghText, FillColor: ghWhite,
		},
		H2: Styler{
			Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
			TextColor: ghText, FillColor: ghWhite,
		},
		H3: Styler{
			Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
			TextColor: ghText, FillColor: ghWhite,
		},
		H4: Styler{
			Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
			TextColor: ghText, FillColor: ghWhite,
		},
		H5: Styler{
			Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
			TextColor: ghText, FillColor: ghWhite,
		},
		H6: Styler{
			Font: "LiberationSans", Style: "b", Size: 13.6, Spacing: 5,
			TextColor: ghMuted, FillColor: ghWhite,
		},
		Line: HorizontalRule{
			Height: 0.75,
			Color:  colors.New(209, 217, 224),
		},
	}

	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "i", Size: 12, Spacing: 2,
		TextColor: ghMuted, FillColor: ghWhite,
	}

	// Table Header Text
	r.THeader = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghCodeBg,
	}

	// Table Body Text
	r.TBody = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghWhite,
	}

	// Horizontal Rule
	r.HorizontalRule = HorizontalRule{
		Height: 3,
		Color:  colors.New(209, 217, 224),
	}

	return &r
}
