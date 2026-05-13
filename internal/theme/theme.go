package theme

import (
	"encoding/json"
	"log"
	"os"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type Theme struct {
	BackgroundColor colors.Color
	// Normal
	Normal   Styler
	NormalEm float64

	// Link
	Link Styler

	// backticked text
	Backtick Styler

	// blockquote text
	Blockquote  Styler
	IndentValue float64

	// Headings
	H1 Styler
	H2 Styler
	H3 Styler
	H4 Styler
	H5 Styler
	H6 Styler

	// Table styling
	THeader Styler
	TBody   Styler

	// code styling
	Code         Styler
	CodeTabWidth int

	// Codeblock syntax highlight colors keyed by highlight group name
	Codeblock map[string]colors.Color
}

// Styler is the struct to capture the styling features for text
// Size and Spacing are specified in points.
// The sum of Size and Spacing is used as line height value
// in the fpdf API
type Styler struct {
	Font      string
	Style     string
	Size      float64
	Spacing   float64
	TextColor colors.Color
	FillColor colors.Color
}

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
	r.Code = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghCodeBg,
	}
	r.CodeTabWidth = 4

	// Headings
	r.H1 = Styler{
		Font: "LiberationSans", Style: "b", Size: 32, Spacing: 5,
		TextColor: ghText, FillColor: ghWhite,
	}
	r.H2 = Styler{
		Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
		TextColor: ghText, FillColor: ghWhite,
	}
	r.H3 = Styler{
		Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
		TextColor: ghText, FillColor: ghWhite,
	}
	r.H4 = Styler{
		Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
		TextColor: ghText, FillColor: ghWhite,
	}
	r.H5 = Styler{
		Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
		TextColor: ghText, FillColor: ghWhite,
	}
	r.H6 = Styler{
		Font: "LiberationSans", Style: "b", Size: 13.6, Spacing: 5,
		TextColor: ghMuted, FillColor: ghWhite,
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

	// Codeblock syntax highlight colors (GitHub Light)
	ghPurple := colors.New(130, 80, 223)   // #8250df
	ghDkBlue := colors.New(5, 80, 174)     // #0550ae
	ghRed := colors.New(207, 34, 46)       // #cf222e
	ghStrBlue := colors.New(10, 48, 105)   // #0a3069
	ghBrown := colors.New(149, 56, 0)      // #953800
	ghComment := colors.New(110, 119, 129) // #6e7781

	r.Codeblock = map[string]colors.Color{
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
	}

	return &r
}

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
	r.Code = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: ghText, FillColor: ghCodeBg,
	}
	r.CodeTabWidth = 4

	// Headings
	r.H1 = Styler{
		Font: "LiberationSans", Style: "b", Size: 32, Spacing: 5,
		TextColor: ghText, FillColor: ghBg,
	}
	r.H2 = Styler{
		Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
		TextColor: ghText, FillColor: ghBg,
	}
	r.H3 = Styler{
		Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
		TextColor: ghText, FillColor: ghBg,
	}
	r.H4 = Styler{
		Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
		TextColor: ghText, FillColor: ghBg,
	}
	r.H5 = Styler{
		Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
		TextColor: ghText, FillColor: ghBg,
	}
	r.H6 = Styler{
		Font: "LiberationSans", Style: "b", Size: 13.6, Spacing: 5,
		TextColor: ghMuted, FillColor: ghBg,
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

	// Codeblock syntax highlight colors (GitHub Dark)
	ghRedSyn := colors.New(255, 123, 114)  // #ff7b72
	ghBlueSyn := colors.New(121, 192, 255) // #79c0ff
	ghStrCyan := colors.New(165, 214, 255) // #a5d6ff
	ghOrange := colors.New(255, 166, 87)   // #ffa657
	ghGray := colors.New(139, 148, 158)    // #8b949e

	r.Codeblock = map[string]colors.Color{
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
	}

	return &r
}

func CustomTheme(themeJSONFile string) *Theme {
	r := Theme{}
	config, err := os.ReadFile(themeJSONFile)
	if err != nil {
		log.Fatal(err)
	}
	// Fill the instance from the JSON file content
	err = json.Unmarshal(config, &r)
	// Check if is there any error while filling the instance
	if err != nil {
		log.Fatal("Error parsing ", themeJSONFile, ":\n", err)
	}

	return &r
}
