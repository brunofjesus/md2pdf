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

func New(red, green, blue int) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

// MarshalJSON outputs the color as a hex string like "#rrggbb".
func (c Color) MarshalJSON() ([]byte, error) {
	hex := fmt.Sprintf("\"#%02x%02x%02x\"", c.Red, c.Green, c.Blue)
	return []byte(hex), nil
}

// UnmarshalJSON accepts both a hex string ("#rrggbb") and the legacy
// object format ({"Red":0,"Green":0,"Blue":0}).
func (c *Color) UnmarshalJSON(data []byte) error {
	// Try string first
	var hex string
	if err := json.Unmarshal(data, &hex); err == nil {
		hex = strings.TrimPrefix(hex, "#")
		if len(hex) != 6 {
			return fmt.Errorf("invalid hex color length: %q", hex)
		}
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &c.Red, &c.Green, &c.Blue)
		return err
	}

	return nil
}
