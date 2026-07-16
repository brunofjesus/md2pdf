package renderer

import (
	"encoding/json"
	"fmt"
)

// marshalRuneChar encodes a rune as a single-character JSON string.
func marshalRuneChar(r rune) ([]byte, error) {
	return json.Marshal(string(r))
}

// unmarshalRuneChar decodes a single-character JSON string into a rune.
func unmarshalRuneChar(data []byte) (rune, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return 0, err
	}

	runes := []rune(s)
	if len(runes) != 1 {
		return 0, fmt.Errorf("expected a single character, got %q", s)
	}

	return runes[0], nil
}
