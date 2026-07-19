## Markdown to PDF

A CLI utility which, as the name implies, generates a PDF from Markdown.

This is a fork of [solworktech/md2pdf](https://github.com/solworktech/md2pdf). I started doing some little tweaks for my own usage, it ended up deviating a lot from the original project, so I decided to start refactoring some more stuff, share it and probably mantain my own forked version.

### Key differences from the original

- **Built-in UTF-8 support** — Liberation Sans and Liberation Mono fonts are embedded; no need for external font files or unicode encoding flags
- **Embedded syntax highlighting** — syntax definition files are bundled into the binary; no git submodule or `-s` flag needed
- **Code block improvements** — code blocks use Liberation Mono with background highlighting
- **Restructured internals** — renderer, theme, colors, fonts and highlight logic moved into `internal/` packages
- **Customisable headers and footers** — configurable page headers and footers (marginals) with left/center/right sections, text placeholders (`%TITLE%`, `%AUTHOR%`, `%PAGE_NUMBER%`, `%PAGE_TOTAL%`) and background images, defined via JSON; includes a built-in default header and footer.

This package depends on two other packages:
- [gomarkdown](https://github.com/gomarkdown/markdown) parser to read the markdown source
- [fpdf](https://codeberg.org/go-pdf/fpdf) to generate the PDF

---

## Features

- Syntax highlighting (for code blocks) — built-in, no configuration needed
- [Dark and light themes](#custom-themes)
- [Customised themes (by passing a JSON file to `md2pdf`)](#custom-themes)
- [Auto Generation of Table of Contents](#auto-generation-of-table-of-contents)
- Built-in UTF-8 support (Liberation Sans / Liberation Mono)
- [Pagination control (using horizontal lines - especially useful for presentations)](#additional-options)
- [Page headers and footers with sensible defaults](#headers-and-footers)
- [Fully customisable headers and footers via JSON (marginals)](#custom-headers-and-footers)

## Supported Markdown elements

- Emphasised and strong text 
- Headings 1-6
- Ordered and unordered lists
- Nested lists
- Images
- Tables
- Links
- Code blocks and backticked text

## Installation 

Build from source:

```sh
$ go install github.com/brunofjesus/md2pdf/cmd/md2pdf@latest
```

## Syntax highlighting

`md2pdf` includes built-in syntax highlighting for code blocks via embedded [gohighlight](https://github.com/jessp01/gohighlight) syntax files. No external files or configuration needed — just annotate your code blocks with the language name.

*Note: when annotating the code block to specify the language, the
annotation name must match the syntax base filename.*

## Custom themes

`md2pdf` supports both light and dark themes out of the box (use `--theme light` or `--theme dark` - no config required). 

However, if you wish to customise the font faces, sizes and colours, you can use the JSONs in
[custom_themes](./custom_themes) as a starting point. Edit to your liking and pass `--theme /path/to/json` to `md2pdf`

## Auto Generation of Table of Contents

`md2pdf` can automatically generate a TOC where each item corresponds to a header in the doc and include it in the first page.
TOC items can then be clicked to navigate to the relevant section (similar to HTML `<a>` anchors).

To make use of this feature, simply pass `--generate-toc` as an argument.

## Quick start

```sh
$ go run ./cmd/md2pdf -i input.md -o output.pdf
```

To convert multiple MD files into a single PDF, use:
```sh
$ go run ./cmd/md2pdf -i /path/to/md/directory -o output.pdf
```

### Additional options

```sh
--input string, -i string                  Input filename, dir consisting of .md|.markdown files or HTTP(s) URL; default is os.Stdin
--output string, -o string                 Output PDF filename; required (default: "out.pdf")
--title string, -t string                  PDF title
--theme string                             Theme to use for the PDF; Can be 'light', 'dark' or the path for a custom theme file (default: "light")
--table-of-contents, --toc                 Generate a table of contents page based on the headings in the input markdown
--horizontal-rule-new-page, --hr-new-page  Start a new page on horizontal rules (---); useful for presentations
--force-overwrite, -f                      Force overwrite of output file if it already exists
--header                                    Print doc header with the title on the left and the author on the right
--footer                                    Print doc footer with the current page number and the total page count
--header-file string                       Load a custom header configuration from a JSON file (see Custom headers and footers)
--footer-file string                       Load a custom footer configuration from a JSON file (see Custom headers and footers)
--page-size string                         Page size for the PDF; can be 'A1', 'A2', 'A3', 'A4', 'A5', 'A6', 'A7', 'Letter', 'Legal' or 'Tabloid' (default: A4)
--orientation string                       Page orientation for the PDF; can be 'portrait' or 'landscape'; default is 'portrait' (default: "portrait")
--author string                            Author name
--log-file string                          Path to log file
--help, -h                                 show help
--version, -v                              print the version
```

For example, the below will:

- Set the title to `My Grand Title`
- Set `Random Bloke` as the author (used in the header/footer)
- Set the dark theme
- Start a new page when encountering an HR (`---`); useful for creating presentations
- Print the default header (title + author) and the default footer (page number / total pages)

```sh
$ go run ./cmd/md2pdf -i /path/to/md \
    -o /path/to/pdf --title "My Grand Title" --author "Random Bloke" \
    --theme dark --hr-new-page --header --footer
```

## Headers and Footers

`md2pdf` can render a header and/or a footer on every page.

For the common case, just pass `--header` and/or `--footer`:

- `--header` prints the document **title** on the left and the **author** on the right.
- `--footer` prints `page number / total pages`, centered.

```sh
$ go run ./cmd/md2pdf -i input.md -o out.pdf --header --footer \
    --title "My Grand Title" --author "Random Bloke"
```

### Custom headers and footers

For full control over content, positioning and styling, pass a JSON file via
`--header-file` and/or `--footer-file`. Headers and footers share the same
schema (a "marginal").

```sh
$ go run ./cmd/md2pdf -i input.md -o out.pdf \
    --header-file header.json --footer-file footer.json
```

A marginal is split into three columns — `left`, `center` and `right` — each
holding a list of sections. A section renders either a text snippet or a
background image.

**Marginal**

| Field             | Type     | Description                                              |
| ----------------- | -------- | -------------------------------------------------------- |
| `backgroundColor` | hex str  | Optional band background colour, e.g. `"#eba0ac"`        |
| `height`          | number   | Height of the band in mm                                 |
| `left`            | array    | Sections aligned to the left edge                        |
| `center`          | array    | Sections centered horizontally                           |
| `right`           | array    | Sections aligned to the right edge                       |

**Section**

| Field             | Type     | Description                                              |
| ----------------- | -------- | -------------------------------------------------------- |
| `width`           | number   | Optional fixed width in mm (auto-sized from text if 0)   |
| `height`          | number   | Section height in mm                                     |
| `relativeX`       | number   | Horizontal offset from the column anchor in mm           |
| `relativeY`       | number   | Vertical offset from the top of the band in mm           |
| `backgroundImage` | path     | Optional image drawn in the section (needs `width`/`height`) |
| `text`            | object   | Optional text content (see below)                        |

**Text**

| Field                 | Type       | Description                                                       |
| --------------------- | ---------- | ---------------------------------------------------------------- |
| `text`                | string     | The text to render; supports placeholders (see below)            |
| `fontSize`            | number     | Font size in points                                              |
| `fontStyle`           | str array  | Any of `"B"` (bold), `"I"` (italic), `"U"` (underline), `"S"` (strikethrough) |
| `color`               | hex str    | Text colour, e.g. `"#1a1a1a"`                                    |
| `horizontalAlignment` | str array  | Any of `"L"` (left), `"C"` (center), `"R"` (right)              |
| `verticalAlignment`   | str array  | Any of `"T"` (top), `"M"` (middle), `"B"` (bottom), `"A"` (baseline) |

**Placeholders** (substituted at render time):

| Placeholder     | Replaced with            |
| --------------- | ------------------------ |
| `%TITLE%`       | The document title       |
| `%AUTHOR%`      | The author name          |
| `%PAGE_NUMBER%` | The current page number  |
| `%PAGE_TOTAL%`  | The total number of pages |

#### Example header

A header with a logo and title on the left, a centered "CONFIDENTIAL" note, and
the author on the right:

```json
{
  "backgroundColor": "#eba0ac",
  "height": 40,
  "left": [
    {
      "width": 71.3,
      "height": 30,
      "relativeX": 10,
      "relativeY": 5,
      "backgroundImage": "./assets/logo.png"
    },
    {
      "height": 10,
      "relativeX": 84,
      "relativeY": 5,
      "text": {
        "text": "%TITLE%",
        "fontSize": 11,
        "fontStyle": ["B"],
        "color": "#1a1a1a",
        "horizontalAlignment": ["L"],
        "verticalAlignment": ["M"]
      }
    }
  ],
  "center": [
    {
      "height": 10,
      "relativeY": 5,
      "text": {
        "text": "CONFIDENTIAL",
        "fontSize": 9,
        "fontStyle": ["I"],
        "color": "#d20f39",
        "horizontalAlignment": ["C"],
        "verticalAlignment": ["M"]
      }
    }
  ],
  "right": [
    {
      "height": 10,
      "relativeX": 10,
      "relativeY": 5,
      "text": {
        "text": "%AUTHOR%",
        "fontSize": 9,
        "fontStyle": ["I"],
        "color": "#666666",
        "horizontalAlignment": ["R"],
        "verticalAlignment": ["M"]
      }
    }
  ]
}
```

#### Example footer

A footer with the title on the left, `Page X of Y` in the center, and the author
on the right:

```json
{
  "backgroundColor": "#f2cdcd",
  "height": 15,
  "left": [
    {
      "height": 10,
      "relativeX": 10,
      "relativeY": 2,
      "text": {
        "text": "%TITLE%",
        "fontSize": 8,
        "fontStyle": ["I"],
        "color": "#888888",
        "horizontalAlignment": ["L"],
        "verticalAlignment": ["M"]
      }
    }
  ],
  "center": [
    {
      "height": 10,
      "relativeY": 2,
      "text": {
        "text": "Page %PAGE_NUMBER% of %PAGE_TOTAL%",
        "fontSize": 9,
        "color": "#333333",
        "horizontalAlignment": ["C"],
        "verticalAlignment": ["M"]
      }
    }
  ],
  "right": [
    {
      "height": 10,
      "relativeX": 10,
      "relativeY": 2,
      "text": {
        "text": "%AUTHOR%",
        "fontSize": 8,
        "fontStyle": ["I"],
        "color": "#888888",
        "horizontalAlignment": ["R"],
        "verticalAlignment": ["M"]
      }
    }
  ]
}
```

## Tests

The tests included in this repo (see the `testdata` folder) were taken from the BlackFriday package.
While the tests may complete without errors, visual inspection of the created PDF is the
only way to determine if the tests *really* pass!

The tests create log files that trace the [gomarkdown](https://github.com/gomarkdown/markdown) parser
callbacks. This is a valuable debugging tool, showing each callback 
and the data provided while the AST is presented.

## Limitations and Known Issues

- It is common for Markdown to include HTML. HTML is treated as a "code block". *There is no attempt to convert raw HTML to PDF.*
- Github-flavoured Markdown permits strikethrough using tildes. This is not supported by `fpdf` as a font style at present.
- The markdown link title (which would show when converted to HTML as hover-over text) is not supported. The generated PDF will show the URL, but this is a function of the PDF viewer.
- Definition lists are not supported
- Text styling (font, size, spacing, style, fill colour, text colour) can be customised via the theme system. Note that fill colour only works when using `CellFormat()`. This is the case for tables, code blocks, and backticked text.

