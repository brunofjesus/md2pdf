// Package theme defines the Theme struct, which captures all the styling
// information for the PDF output.
//
// It also includes a function to read a JSON file and create a Theme instance
// from it.
//
// Two themes are included by default: LightTheme and DarkTheme, inspired by
// GitHub's light and dark modes, respectively.
package theme

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/brunofjesus/md2pdf/v3/internal/colors"
)

// Theme captures all the styling information for the PDF output.
type Theme struct {
	BackgroundColor colors.Color `json:"backgroundColor"`

	// Normal
	Normal Styler `json:"normal"`

	// Link
	Link Styler `json:"link"`

	// backticked text
	Backtick Styler `json:"backtick"`

	// blockquote text
	Blockquote  Styler  `json:"blockquote"`
	IndentValue float64 `json:"indentValue"`

	// Headings
	Heading Heading `json:"heading"`

	// Table styling
	Table Table `json:"table"`

	// code styling
	Code Code `json:"code"`

	// other
	HorizontalRule HorizontalRule `json:"horizontalRule"`
}

// Table captures the styling for tables, including header and body.
type Table struct {
	Header Styler `json:"header"`
	Body   Styler `json:"body"`
}

// Heading captures the styling for headings (h1 to h6) and the horizontal line
// below them.
type Heading struct {
	H1   Styler         `json:"h1"`
	H2   Styler         `json:"h2"`
	H3   Styler         `json:"h3"`
	H4   Styler         `json:"h4"`
	H5   Styler         `json:"h5"`
	H6   Styler         `json:"h6"`
	Line HorizontalRule `json:"line"`
}

// HorizontalRule captures the styling for horizontal rules, including height
// and color.
type HorizontalRule struct {
	Height float64      `json:"height"`
	Color  colors.Color `json:"color"`
}

// Code captures the styling for code blocks, including text styling, tab width,
// and syntax highlighting colors.
type Code struct {
	Text     Styler                  `json:"text"`
	TabWidth int                     `json:"tabWidth"`
	Colors   map[string]colors.Color `json:"colors"`
}

// Styler is the struct to capture the styling features for text
// Size and Spacing are specified in points.
// The sum of Size and Spacing is used as line height value
// in the fpdf API.
type Styler struct {
	Font      string       `json:"font"`
	Style     string       `json:"style"`
	Size      float64      `json:"size"`
	Spacing   float64      `json:"spacing"`
	TextColor colors.Color `json:"textColor"`
	FillColor colors.Color `json:"fillColor"`
}

// CustomTheme reads a JSON file and returns a Theme instance filled with the
// data from the file. The JSON file should have the same structure as the Theme
// struct.
func CustomTheme(themeJSONFile string) *Theme {
	r := new(Theme)

	config, err := os.ReadFile(filepath.Clean(themeJSONFile))
	if err != nil {
		log.Fatal(err)
	}
	// Fill the instance from the JSON file content
	err = json.Unmarshal(config, r)
	// Check if is there any error while filling the instance
	if err != nil {
		log.Fatal("Error parsing ", themeJSONFile, ":\n", err)
	}

	return r
}
