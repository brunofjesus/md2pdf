package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/gomarkdown/markdown/parser"
	"github.com/brunofjesus/md2pdf/internal/renderer"
	"golang.org/x/exp/slices"
)

var (
	input             = flag.String("i", "", "Input filename, dir consisting of .md|.markdown files or HTTP(s) URL; default is os.Stdin")
	output            = flag.String("o", "", "Output PDF filename; required")
	title             = flag.String("title", "", "Presentation title")
	author            = flag.String("author", "", "Author's name; used if -footer is passed")
	themeArg          = flag.String("theme", "light", "[light | dark | /path/to/custom/theme.json]")
	hrAsNewPage       = flag.Bool("new-page-on-hr", false, "Interpret HR as a new page; useful for presentations")
	printFooter       = flag.Bool("with-footer", false, "Print doc footer (<author>  <title>  <page number>)")
	generateTOC       = flag.Bool("generate-toc", false, "Auto Generate Table of Contents (TOC)")
	pageSize          = flag.String("page-size", "A4", "[A3 | A4 | A5]")
	orientation       = flag.String("orientation", "portrait", "[portrait | landscape]")
	logFile           = flag.String("log-file", "", "Path to log file")
	help              = flag.Bool("help", false, "Show usage message")
	ver               = flag.Bool("version", false, "Print version and build info")
	version           = "dev"
	commit            = "none"
	date              = "unknown"
	_, fileName, _, _ = runtime.Caller(0)
)

var opts []renderer.RenderOption

func processRemoteInputFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code: " + fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
	content, rerr := io.ReadAll(resp.Body)
	return content, rerr
}

func glob(dir string, validExts []string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if slices.Contains(validExts, filepath.Ext(path)) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	if *help {
		usage("")
		return
	}

	if *ver {
		fmt.Printf("md2pdf version: %s, commit: %s, built on: %s\n", version, commit, date)
		return
	}

	if *output == "" {
		usage("Output PDF filename is required")
	}

	if *hrAsNewPage {
		opts = append(opts, renderer.IsHorizontalRuleNewPage(true))
	}

	// get text for PDF
	var content []byte
	var err error
	var inputBaseURL string
	if *input == "" {
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		httpRegex := regexp.MustCompile("^http(s)?://")
		if httpRegex.Match([]byte(*input)) {
			content, err = processRemoteInputFile(*input)
			if err != nil {
				log.Fatal(err)
			}
			// get the base URL so we can adjust relative links and images
			inputBaseURL = strings.Replace(filepath.Dir(*input), ":/", "://", 1)
		} else {
			fileInfo, err := os.Stat(*input)
			if err != nil {
				log.Fatal(err)
			}

			if fileInfo.IsDir() {
				opts = append(opts, renderer.IsHorizontalRuleNewPage(true))
				validExts := []string{".md", ".markdown"}
				files, err := glob(*input, validExts)
				if err != nil {
					log.Fatal(err)
				}
				for i, filePath := range files {
					fileContents, err := os.ReadFile(filePath)
					if err != nil {
						log.Fatal(err)
					}
					content = append(content, fileContents...)
					if i < len(files)-1 {
						content = append(content, []byte("---\n")...)
					}
				}
			} else {
				content, err = os.ReadFile(*input)
				if err != nil {
					log.Fatal(err)
				}
				if absInput, absErr := filepath.Abs(*input); absErr == nil {
					inputBaseURL = filepath.Dir(absInput)
				}
			}
		}
	}

	theme := renderer.LIGHT
	themeFile := ""
	if *themeArg == "dark" {
		theme = renderer.DARK
	} else if _, err := os.Stat(*themeArg); err == nil {
		theme = renderer.CUSTOM
		themeFile = *themeArg
	}

	params := renderer.PdfRendererParams{
		Orientation:     *orientation,
		Papersz:         *pageSize,
		PdfFile:         *output,
		TracerFile:      *logFile,
		Opts:            opts,
		Theme:           theme,
		CustomThemeFile: themeFile,
	}

	pf := renderer.NewPdfRenderer(params)

	if inputBaseURL != "" {
		pf.InputBaseURL = inputBaseURL
	}
	pf.Pdf.SetSubject(*title, true)
	pf.Pdf.SetTitle(*title, true)
	pf.Extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode | parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs | parser.BackslashLineBreak | parser.DefinitionLists

	if *printFooter {
		pf.Pdf.SetFooterFunc(func() {
			pf.Pdf.SetFillColor(pf.Theme.BackgroundColor.Red, pf.Theme.BackgroundColor.Green, pf.Theme.BackgroundColor.Blue)
			// Position at 1.5 cm from bottom
			pf.Pdf.SetY(-15)
			// Arial italic 8
			pf.Pdf.SetFont(pf.Theme.Normal.Font, "I", 8)
			// Text color in gray
			pf.Pdf.SetTextColor(128, 128, 128)
			w, h, _ := pf.Pdf.PageSize(pf.Pdf.PageNo())
			pf.Pdf.SetX(4)
			pf.Pdf.CellFormat(0, 10, *author, "", 0, "", true, 0, "")
			middle := w / 2
			if *orientation == "landscape" {
				middle = h / 2
			}
			pf.Pdf.SetX(middle - float64(len(*title)))
			pf.Pdf.CellFormat(0, 10, *title, "", 0, "", true, 0, "")
			pf.Pdf.SetX(-40)
			pf.Pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pf.Pdf.PageNo()), "", 0, "", true, 0, "")
		})
	}

	var p renderer.Processor = pf
	if *generateTOC {
		p = renderer.NewTOCDecorator(pf)
	}
	err = p.Process(content)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Printf("Usage: %s (%s) [options]\n", filepath.Base(fileName), version)
	flag.PrintDefaults()
	os.Exit(0)
}
