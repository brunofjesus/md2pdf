package node

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/canhlinh/svg2png"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gomarkdown/markdown/ast"
)

// ProcessImage handles *ast.Image entering/leaving.
func ProcessImage(ctx PdfContext, n ast.Node, entering bool) {
	node, ok := n.(*ast.Image)
	if !ok {
		ctx.Tracer("Image: not an Image", "")
		return
	}

	if entering {
		ctx.Cr()

		destination, err := ResolveImagePath(ctx.GetInputBaseURL(), string(node.Destination))
		if err != nil {
			ctx.Tracer("Image (resolve error)", err.Error())
			log.Println(err)

			return
		}

		ctx.Tracer("Image (entering)",
			fmt.Sprintf("Destination[%v] Title[%v]",
				destination,
				string(node.Title)))
		imgPath := destination

		_, err = os.Stat(imgPath)
		if err == nil {
			ctx.GetPdf().ImageOptions(destination,
				-1, 0, 0, 0, true,
				fpdf.ImageOptions{ImageType: "", ReadDpi: true, AllowNegativePosition: false}, 0, "")
		} else {
			ctx.Tracer("Image (file error)", err.Error())
		}
	} else {
		ctx.Tracer("Image (leaving)", "")
	}
}

// ResolveImagePath resolves a raw image destination (local path, base-URL,
// relative path, HTTP URL, or SVG) into a local file path ready to draw.
// baseURL is the base against which relative destinations are resolved; it may
// be a local directory or an HTTP base URL (empty when there is none).
func ResolveImagePath(baseURL, destination string) (string, error) {
	tempDir := os.TempDir() + "/" + filepath.Base(os.Args[0])

	source, needsDownload, err := locateImageSource(baseURL, destination)
	if err != nil {
		return "", err
	}

	if needsDownload {
		if mkErr := os.MkdirAll(tempDir, 0o750); mkErr != nil {
			return "", fmt.Errorf("failed to create temp directory %s: %w", tempDir, mkErr)
		}

		localPath := tempDir + "/" + filepath.Base(destination)
		if dlErr := downloadFile(source, localPath); dlErr != nil {
			return "", fmt.Errorf("failed to download image from %s: %w", source, dlErr)
		}

		destination = localPath
		fmt.Println("Downloaded image to: " + destination)
	} else {
		destination = source
	}

	mtype, _ := mimetype.DetectFile(destination)
	//nolint:nestif
	if mtype.Is("image/svg+xml") {
		re := regexp.MustCompile(`<svg\s*.*\s*width="([0-9\.]+)"\sheight="([0-9\.]+)".*>`)
		contents, _ := os.ReadFile(filepath.Clean(destination))
		matches := re.FindStringSubmatch(string(contents))

		tf, err := os.CreateTemp(tempDir, "*.svg")
		if err != nil {
			log.Println(err)
			return "", err
		}

		if _, err := tf.Write(contents); err != nil {
			_ = tf.Close()
			log.Println(err) //nolint:wsl_v5

			return "", err
		}

		if err := tf.Close(); err != nil {
			log.Println(err)
			return "", err
		}

		if renameErr := os.Rename(destination, tf.Name()); renameErr != nil {
			log.Println(renameErr)
			return "", renameErr
		}

		destination = tf.Name()
		width, _ := strconv.ParseFloat(matches[1], 64)
		height, _ := strconv.ParseFloat(matches[2], 64)
		chrome := svg2png.NewChrome().SetHeight(int(height)).SetWith(int(width))
		outputFileName := destination + ".png"

		//nolint:misspell
		if err := chrome.Screenshoot(destination, outputFileName); err != nil {
			log.Println(err)
			return "", err
		}

		destination = outputFileName
	}

	return destination, nil
}

// locateImageSource decides how a raw destination should be obtained. It
// returns the source to use, whether that source must be downloaded, and an
// error when a local image cannot be found.
//
// Resolution order:
//   - absolute HTTP destination -> download as-is
//   - destination existing on disk -> use as-is
//   - local (non-http) base URL -> join and use if it exists, else "not found"
//   - HTTP base URL -> join and download
//   - otherwise -> "not found"
func locateImageSource(baseURL, destination string) (string, bool, error) {
	if strings.HasPrefix(destination, "http") {
		return destination, true, nil
	}

	if _, statErr := os.Stat(destination); statErr == nil {
		return destination, false, nil
	}

	if baseURL != "" && !strings.HasPrefix(baseURL, "http") {
		localPath := joinBase(baseURL, destination)
		if _, statErr := os.Stat(localPath); statErr == nil {
			return localPath, false, nil
		}

		return "", false, fmt.Errorf("image not found: %s", destination)
	}

	if strings.HasPrefix(baseURL, "http") {
		return joinBase(baseURL, destination), true, nil
	}

	return "", false, fmt.Errorf("image not found: %s", destination)
}

func downloadFile(url, fileName string) error {
	//nolint:exhaustruct
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			fmt.Println("Redirected to:", req.URL)
			return nil
		},
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "curl/7.84.0")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return errors.New("Received non 200 response code: " + fmt.Sprintf("HTTP %d", response.StatusCode))
	}

	file, err := os.Create(filepath.Clean(fileName))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
