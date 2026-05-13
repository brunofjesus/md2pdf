package node

import (
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

	"codeberg.org/go-pdf/fpdf"
	"github.com/canhlinh/svg2png"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gomarkdown/markdown/ast"
)

// ProcessImage handles *ast.Image entering/leaving.
func ProcessImage(ctx PdfContext, n ast.Node, entering bool) {
	node := n.(*ast.Image)
	if entering {
		ctx.Cr()
		destination := string(node.Destination)
		tempDir := os.TempDir() + "/" + filepath.Base(os.Args[0])
		_, err := os.Stat(destination)
		if errors.Is(err, os.ErrNotExist) &&
			!strings.HasPrefix(destination, "http") &&
			ctx.GetInputBaseURL() != "" &&
			!strings.HasPrefix(ctx.GetInputBaseURL(), "http") {
			localPath := filepath.Join(ctx.GetInputBaseURL(), destination)
			if _, lerr := os.Stat(localPath); lerr == nil {
				destination = localPath
				err = nil
			}
		}
		if errors.Is(err, os.ErrNotExist) {
			var source string = destination
			if !strings.HasPrefix(destination, "http") {
				if ctx.GetInputBaseURL() != "" {
					source = ctx.GetInputBaseURL() + "/" + destination
				}
			}
			os.MkdirAll(tempDir, 755)
			err := downloadFile(source, tempDir+"/"+filepath.Base(destination))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				destination = tempDir + "/" + filepath.Base(destination)
				fmt.Println("Downloaded image to: " + destination)
			}
		}
		mtype, err := mimetype.DetectFile(destination)
		if mtype.Is("image/svg+xml") {
			re := regexp.MustCompile(`<svg\s*.*\s*width="([0-9\.]+)"\sheight="([0-9\.]+)".*>`)
			contents, _ := os.ReadFile(destination)
			matches := re.FindStringSubmatch(string(contents))
			tf, err := os.CreateTemp(tempDir, "*.svg")
			if err != nil {
				log.Println(err)
				return
			}

			if _, err := tf.Write(contents); err != nil {
				tf.Close()
				log.Println(err)
				return
			}
			if err := tf.Close(); err != nil {
				log.Println(err)
				return
			}
			os.Rename(destination, tf.Name())
			destination = tf.Name()
			width, _ := strconv.ParseFloat(matches[1], 64)
			height, _ := strconv.ParseFloat(matches[2], 64)
			chrome := svg2png.NewChrome().SetHeight(int(height)).SetWith(int(width))
			outputFileName := destination + ".png"
			if err := chrome.Screenshoot(destination, outputFileName); err != nil {
				log.Println(err)
				return
			}
			destination = outputFileName
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
				fpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
		} else {
			ctx.Tracer("Image (file error)", err.Error())
		}
	} else {
		ctx.Tracer("Image (leaving)", "")
	}
}

func downloadFile(url, fileName string) error {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("Redirected to:", req.URL)
			return nil
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "curl/7.84.0")
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code: " + fmt.Sprintf("HTTP %d", response.StatusCode))
	}
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
