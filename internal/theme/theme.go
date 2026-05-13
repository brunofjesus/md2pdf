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

	r.BackgroundColor = colors.Lookup("white")
	// Normal Text
	r.Normal = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("black"),
		FillColor: colors.Lookup("white"),
	}

	// Link text
	r.Link = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("cornflowerblue"),
	}

	// Backticked text
	r.Backtick = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.New(37, 27, 14), FillColor: colors.New(200, 200, 200),
	}

	// Quoted Text

	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.New(37, 27, 14), FillColor: colors.New(200, 200, 200),
	}

	// Code text
	r.Code = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.New(37, 27, 14), FillColor: colors.New(200, 200, 200),
	}
	r.CodeTabWidth = 4

	// Headings
	r.H1 = Styler{
		Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}
	r.H2 = Styler{
		Font: "LiberationSans", Style: "b", Size: 22, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}
	r.H3 = Styler{
		Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}
	r.H4 = Styler{
		Font: "LiberationSans", Style: "b", Size: 18, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}
	r.H5 = Styler{
		Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}
	r.H6 = Styler{
		Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}

	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "i", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("black"), FillColor: colors.Lookup("white"),
	}

	// Table Header Text
	r.THeader = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("black"), FillColor: colors.New(180, 180, 180),
	}

	// Table Body Text
	r.TBody = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("black"), FillColor: colors.New(240, 240, 240),
	}

	return &r
}

func DarkTheme() *Theme {
	r := Theme{}

	r.BackgroundColor = colors.Lookup("black")

	// Normal Text
	r.Normal = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("white"),
	}

	// Quoted Text
	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("white"),
	}

	// Link text
	r.Link = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("cornflowerblue"),
	}

	// Backticked text
	r.Backtick = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("lightgrey"), FillColor: colors.New(32, 35, 37),
	}

	// Code text
	r.Code = Styler{
		Font: "LiberationMono", Style: "", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("lightgrey"), FillColor: colors.New(32, 35, 37),
	}
	r.CodeTabWidth = 4

	// Headings
	r.H1 = Styler{
		Font: "LiberationSans", Style: "b", Size: 24, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}
	r.H2 = Styler{
		Font: "LiberationSans", Style: "b", Size: 22, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}
	r.H3 = Styler{
		Font: "LiberationSans", Style: "b", Size: 20, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}
	r.H4 = Styler{
		Font: "LiberationSans", Style: "b", Size: 18, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}
	r.H5 = Styler{
		Font: "LiberationSans", Style: "b", Size: 16, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}
	r.H6 = Styler{
		Font: "LiberationSans", Style: "b", Size: 14, Spacing: 5,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}

	r.Blockquote = Styler{
		Font: "LiberationSans", Style: "i", Size: 12, Spacing: 2,
		FillColor: colors.Lookup("black"), TextColor: colors.Lookup("darkgray"),
	}

	// Table Header Text
	r.THeader = Styler{
		Font: "LiberationSans", Style: "b", Size: 12, Spacing: 2,
		TextColor: colors.Lookup("darkgray"), FillColor: colors.New(27, 27, 27),
	}

	// Table Body Text
	r.TBody = Styler{
		Font: "LiberationSans", Style: "", Size: 12, Spacing: 2,
		FillColor: colors.New(200, 200, 200), TextColor: colors.New(128, 128, 128),
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
