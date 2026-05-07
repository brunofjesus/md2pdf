package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/brunofjesus/md2pdf/internal/renderer"
)

type InputProcessor func(input string) ([]renderer.RenderOption, []byte, error)

func GetInputProcessor(input string) (InputProcessor, error) {
	if input == "" {
		return processStdinInput, nil
	} else if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return processHTTPInput, nil
	}

	fileInfo, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("failed to stat input: %w", err)
	} else if fileInfo.IsDir() {
		return processDirInput, nil
	} else {
		return processFileInput, nil
	}
}

func processDirInput(input string) ([]renderer.RenderOption, []byte, error) {
	var content []byte

	files, err := glob(input, []string{".md", ".markdown"})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to glob directory: %w", err)
	}

	for i, filePath := range files {
		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
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

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, content, nil
}

func processFileInput(input string) ([]renderer.RenderOption, []byte, error) {
	content, err := os.ReadFile(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	abs, err := filepath.Abs(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, content, nil
}

func processStdinInput(input string) ([]renderer.RenderOption, []byte, error) {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read from stdin: %w", err)
	}

	abs, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	return []renderer.RenderOption{renderer.WithBaseURL(abs)}, content, err
}

func processHTTPInput(input string) ([]renderer.RenderOption, []byte, error) {
	resp, err := http.Get(input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, nil, errors.New("Received non 200 response code: " + fmt.Sprintf("HTTP %d", resp.StatusCode))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// get the base URL so we can adjust relative links and images
	opts := []renderer.RenderOption{renderer.WithBaseURL(
		strings.Replace(filepath.Dir(input), ":/", "://", 1),
	)}

	return opts, content, err
}

// glob recursively walks the given directory and returns a
// list of file paths that have extensions in validExts
func glob(dir string, validExts []string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(dir, func(path string, f fs.DirEntry, err error) error {
		if slices.Contains(validExts, filepath.Ext(path)) && !f.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
