package node

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testImgName = "img.png"

func TestResolveImagePath(t *testing.T) {
	t.Parallel()

	t.Run("existing absolute local file is returned unchanged", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		file := filepath.Join(dir, testImgName)
		writeFile(t, file)

		got, err := ResolveImagePath("", file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != file {
			t.Errorf("ResolveImagePath = %q, want %q", got, file)
		}
	})

	t.Run("local base resolves an existing relative file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, testImgName))

		got, err := ResolveImagePath(dir, testImgName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := filepath.Join(dir, testImgName)
		if got != want {
			t.Errorf("ResolveImagePath = %q, want %q", got, want)
		}
	})

	t.Run("local base with missing file returns not found error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()

		_, err := ResolveImagePath(dir, "missing.png")
		if err == nil {
			t.Fatal("expected an error for a missing local file, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected a not-found error, got %q", err.Error())
		}
	})

	t.Run("empty base with missing non-http file returns not found error", func(t *testing.T) {
		t.Parallel()

		_, err := ResolveImagePath("", "does/not/exist.png")
		if err == nil {
			t.Fatal("expected an error for a missing file, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected a not-found error, got %q", err.Error())
		}
	})
}

// writeFile creates a small non-SVG file so mimetype detection does not trigger
// the SVG conversion branch.
func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.WriteFile(path, []byte("not an svg, just some bytes"), 0o600); err != nil {
		t.Fatalf("failed to write test file %q: %v", path, err)
	}
}
