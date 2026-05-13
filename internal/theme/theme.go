package theme

import (
	"encoding/json"
	"log"
	"os"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

type Theme struct {
	BackgroundColor colors.Color `json:"background_color"`

	// Normal
	Normal Styler `json:"normal"`

	// Link
	Link Styler `json:"link"`

	// backticked text
	Backtick Styler `json:"backtick"`

	// blockquote text
	Blockquote  Styler  `json:"blockquote"`
	IndentValue float64 `json:"indent_value"`

	// Headings
	Heading Heading `json:"heading"`

	// Table styling
	THeader Styler `json:"table_header"`
	TBody   Styler `json:"table_body"`

	// code styling
	Code Code `json:"code"`

	// other
	HorizontalRule HorizontalRule `json:"horizontal_rule"`
}

type Heading struct {
	H1   Styler         `json:"h1"`
	H2   Styler         `json:"h2"`
	H3   Styler         `json:"h3"`
	H4   Styler         `json:"h4"`
	H5   Styler         `json:"h5"`
	H6   Styler         `json:"h6"`
	Line HorizontalRule `json:"line"`
}

type HorizontalRule struct {
	Height float64      `json:"height"`
	Color  colors.Color `json:"color"`
}

type Code struct {
	Text     Styler                  `json:"text"`
	TabWidth int                     `json:"tab_width"`
	Colors   map[string]colors.Color `json:"colors"`
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
