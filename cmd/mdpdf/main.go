package main

import (
	"context"
	"errors"
	"log"
	"os"

	cliutil "github.com/brunofjesus/md2pdf/internal/cli"
	"github.com/brunofjesus/md2pdf/internal/renderer"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "md2pdf",
		Usage: "md2pdf", // TODO: improve this
		Action: func(_ context.Context, cmd *cli.Command) error {
			flagInput := cmd.String("input")
			flatOutput := cmd.String("output")
			flagTitle := cmd.String("title")
			flagTOC := cmd.Bool("table-of-contents")
			flagHRNewPage := cmd.Bool("horizontal-rule-new-page")
			flagTheme := cmd.String("theme")
			flagForceOverwrite := cmd.Bool("force-overwrite")
			flagFooter := cmd.Bool("footer")
			flagOrientation := cmd.String("orientation")
			flagAuthor := cmd.String("author")

			if !flagForceOverwrite {
				outFile, err := os.Stat(flatOutput)
				if err != nil && !errors.Is(err, os.ErrNotExist) {
					log.Fatalf("error: failed to check output file: %v\n", err)
				}
				if outFile != nil {
					log.Fatalf("error: output file already exists: %s; use -f to overwrite.\n", flatOutput)
				}
			}

			inputProcessor, err := cliutil.GetInputProcessor(flagInput)
			if err != nil {
				return err
			}

			opts, content, err := inputProcessor(flagInput)
			if err != nil {
				return err
			}

			if flagHRNewPage {
				opts = append(opts, renderer.WithHorizontalRuleAsNewPage())
			}

			if flagFooter {
				opts = append(opts, renderer.WithDefaultFooter(flagOrientation, flagAuthor, flagTitle))
			}

			if flagTOC {
				opts = append(opts, renderer.WithTableOfContents())
			}

			params := renderer.PdfRendererParams{
				Title:           flagTitle,
				Orientation:     flagOrientation,
				Papersz:         "",
				PdfFile:         flatOutput,
				TracerFile:      "",
				Opts:            opts,
				Theme:           renderer.LIGHT,
				CustomThemeFile: "",
			}

			switch flagTheme {
			case "light":
				params.Theme = renderer.LIGHT
			case "dark":
				params.Theme = renderer.DARK
			case "custom":
				params.Theme = renderer.CUSTOM
				params.CustomThemeFile = flagTheme
			}

			pf := renderer.NewPdfRenderer(params)

			err = pf.Process(content)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"i"},
				Usage:   "Input filename, dir consisting of .md|.markdown files or HTTP(s) URL; default is os.Stdin",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output PDF filename; required",
				Value:   "out.pdf",
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "PDF title; default is empty string",
			},
			&cli.StringFlag{
				Name:  "theme",
				Usage: "Theme to use for the PDF; Can be 'light', 'dark' or the path for a custom theme file; default is 'light'",
				Value: "light",
			},
			&cli.BoolFlag{
				Name:    "table-of-contents",
				Aliases: []string{"toc"},
				Usage:   "Generate a table of contents page based on the headings in the input markdown",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "horizontal-rule-new-page",
				Aliases: []string{"hr-new-page"},
				Usage:   "Start a new page on horizontal rules (---); useful for presentations",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "force-overwrite",
				Aliases: []string{"f"},
				Usage:   "Force overwrite of output file if it already exists",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "footer",
				Usage: "Print doc footer (<author>  <title>  <page number>)",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "orientation",
				Usage: "Page orientation for the PDF; can be 'portrait' or 'landscape'; default is 'portrait'",
				Value: "portrait",
				Validator: func(value string) error {
					if value != "portrait" && value != "landscape" {
						return errors.New("orientation must be either 'portrait' or 'landscape'")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "author",
				Usage: "Author's name",
				Value: "",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
