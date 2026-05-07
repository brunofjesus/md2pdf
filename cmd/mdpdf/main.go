package main

import (
	"context"
	"fmt"
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

			params := renderer.PdfRendererParams{
				Title:           flagTitle,
				Orientation:     "",
				Papersz:         "",
				PdfFile:         flatOutput,
				TracerFile:      "",
				Opts:            opts,
				Theme:           renderer.LIGHT,
				CustomThemeFile: "",
			}

			pf := renderer.NewPdfRenderer(params)

			var p renderer.Processor = pf
			if flagTOC {
				p = renderer.NewTOCDecorator(pf)
			}

			err = p.Process(content)
			if err != nil {
				fmt.Printf("error: %v\n", err)
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
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output PDF filename; required",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "PDF title; default is empty string",
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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
