package colors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Color is a RGB set of ints; for a nice picker
// see https://www.w3schools.com/colors/colors_picker.asp
type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

// FromRGB creates a new Color with the given RGB values.
func FromRGB(red, green, blue int) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

// FromHex creates a new Color from a hex string like "#rrggbb".
func FromHex(hex string) (*Color, error) {
	c := new(Color)

	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return c, fmt.Errorf("invalid hex color length: %q", hex)
	}

	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &c.Red, &c.Green, &c.Blue)

	return c, err
}

// MarshalJSON outputs the color as a hex string like "#rrggbb".
func (c Color) MarshalJSON() ([]byte, error) {
	hex := fmt.Sprintf("\"#%02x%02x%02x\"", c.Red, c.Green, c.Blue)
	return []byte(hex), nil
}

// UnmarshalJSON accepts a hex string ("#rrggbb").
func (c *Color) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err == nil {
		newColor, err := FromHex(hex)
		if err != nil {
			return err
		}

		c.Red = newColor.Red
		c.Green = newColor.Green
		c.Blue = newColor.Blue
	}

	return nil
}
