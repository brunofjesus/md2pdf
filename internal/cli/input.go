// Package cli provides helper functions for the command line interface of the application.
package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/brunofjesus/md2pdf/v3/internal/renderer"
)

// InputProcessor gets an input string and provides a reader for the content and
// any options needed to render it.
type InputProcessor func(input string) ([]renderer.RenderOption, io.ReadCloser, error)

// GetInputProcessor determines the type of input
// (file, directory, URL, or stdin) and returns the appropriate InputProcessor
// function.
func GetInputProcessor(input string) (InputProcessor, error) {
	if input == "" {
		return processStdinInput, nil
	} else if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return processHTTPInput, nil
	}

	fileInfo, err := os.Stat(input)

	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to stat input: %w", err)
	case fileInfo.IsDir():
		return processDirInput, nil
	default:
		return processFileInput, nil
	}
}

func processDirInput(input string) ([]renderer.RenderOption, io.ReadCloser, error) {
	var content []byte

	files, err := glob(input, []string{".md", ".markdown"})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to glob directory: %w", err)
	}

	for i, path := range files {
		fileContents, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		content = append(content, fileContents...)
		if i < len(files)-1 {
			content = append(content, []byte("---\n")...)
		}
	}

	abs, err := filepath.Abs(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	reader := io.NopCloser(bytes.NewReader(content))

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, reader, nil
}

func processFileInput(input string) ([]renderer.RenderOption, io.ReadCloser, error) {
	file, err := os.Open(filepath.Clean(input))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	abs, err := filepath.Abs(input)
	if err != nil {
		_ = file.Close()

		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, file, nil
}

func processStdinInput(_ string) ([]renderer.RenderOption, io.ReadCloser, error) {
	abs, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, io.NopCloser(os.Stdin), nil
}

func processHTTPInput(input string) ([]renderer.RenderOption, io.ReadCloser, error) {
	//nolint:exhaustruct
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, input, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()

		return nil, nil, errors.New("received non 200 response code: " + fmt.Sprintf("HTTP %d", resp.StatusCode))
	}

	// get the base URL so we can adjust relative links and images
	opts := []renderer.RenderOption{renderer.WithBaseURL(
		strings.Replace(filepath.Dir(input), ":/", "://", 1),
	)}

	return opts, resp.Body, nil
}

// glob recursively walks the given directory and returns a
// list of file paths that have extensions in validExts.
func glob(dir string, validExts []string) ([]string, error) {
	files := []string{}

	err := filepath.WalkDir(dir, func(path string, f fs.DirEntry, _ error) error {
		if slices.Contains(validExts, filepath.Ext(path)) && !f.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
