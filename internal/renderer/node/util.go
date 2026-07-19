package node

import (
	"path/filepath"
	"strings"
)

// joinBase joins a relative destination onto a base (a local directory or an
// HTTP base URL), stripping a leading "./". HTTP bases use URL-style "/"
// concatenation; local bases use filepath.Join.
func joinBase(base, dest string) string {
	dest = strings.TrimPrefix(dest, "./")
	if strings.HasPrefix(base, "http") {
		return strings.TrimRight(base, "/") + "/" + dest
	}

	return filepath.Join(base, dest)
}
