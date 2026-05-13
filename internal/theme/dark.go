package theme

import "github.com/brunofjesus/md2pdf/v3/internal/colors"

func DarkTheme() *Theme {
	r := Theme{}

	ghText := colors.New(230, 237, 243)  // #e6edf3
	ghMuted := colors.New(125, 133, 144) // #7d8590
	ghLink := colors.New(74, 158, 255)   // #4a9eff
	ghCodeBg := colors.New(22, 27, 34)   // #161b22
	ghInlineBg := colors.New(52, 59, 66) // #343b42
	ghBg := colors.New(13, 17, 23)       // #0d1117

	r.BackgroundColor = ghBg

	// Normal Text
	r.Normal = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghBg,
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
	// Codeblock syntax highlight colors (GitHub Dark)
	ghRedSyn := colors.New(255, 123, 114)  // #ff7b72
	ghBlueSyn := colors.New(121, 192, 255) // #79c0ff
	ghStrCyan := colors.New(165, 214, 255) // #a5d6ff
	ghOrange := colors.New(255, 166, 87)   // #ffa657
	ghGray := colors.New(139, 148, 158)    // #8b949e

	r.Code = Code{
		Text: Styler{
			Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
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
	}

	// Headings
	r.Heading = Heading{
		H1: Styler{
			Font: "LiberationSans", Style: "b", Size: 32, Spacing: 5,
			TextColor: ghText, FillColor: ghBg,
		},
		H2: Styler{
			Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
			TextColor: ghText, FillColor: ghBg,
		},
		H3: Styler{
			Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
			TextColor: ghText, FillColor: ghBg,
		},
		H4: Styler{
			Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
			TextColor: ghText, FillColor: ghBg,
		},
		H5: Styler{
			Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
			TextColor: ghText, FillColor: ghBg,
		},
		H6: Styler{
			Font: "LiberationSans", Style: "b", Size: 13.6, Spacing: 5,
			TextColor: ghMuted, FillColor: ghBg,
		},
		Line: HorizontalRule{
			Height: 0.75,
			Color:  colors.New(61, 68, 77),
		},
	}

	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "i", Size: 12, Spacing: 2,
		TextColor: ghMuted, FillColor: ghBg,
	}

	// Table Header Text
	r.THeader = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghCodeBg,
	}

	// Table Body Text
	r.TBody = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghBg,
	}

	// Horizontal Rule
	r.HorizontalRule = HorizontalRule{
		Height: 3,
		Color:  colors.New(61, 68, 77),
	}

	return &r
}
