/*
 * Markdown to PDF Converter
 * Available at http://github.com/solworktech/md2pdf
 *
 * Copyright © Cecil New <cecil.new@gmail.com>, Jesse Portnoy <jesse@packman.io>.
 * Distributed under the MIT License.
 * See README.md for details.
 *
 * Dependencies
 * This package depends on two other packages:
 *
 * Go markdown processor
 *   Available at https://github.com/gomarkdown/markdown
 *
 * fpdf - a PDF document generator with high level support for
 *   text, drawing and images.
 *   Available at https://codeberg.org/go-pdf/fpdf
 */

package renderer

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
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
	highlight "github.com/jessp01/gohighlight"
	"github.com/mitchellh/go-wordwrap"
	syntaxhighlight "github.com/solworktech/md2pdf/v2/internal/highlight"
)

func (r *PdfRenderer) processText(node *ast.Text) {
	currentStyle := r.cs.peek().textStyle
	r.setStyler(currentStyle)
	s := string(node.Literal)
	s = strings.ReplaceAll(s, "\n", " ")
	r.tracer("Text", s)

	if incell {
		r.cs.peek().cellInnerString += s
		r.cs.peek().cellInnerStringStyle = &currentStyle
		return
	}
	switch node.Parent.(type) {

	case *ast.Link:
		r.writeLink(currentStyle, s, r.cs.peek().destination)
	case *ast.Heading:
		if len(r.tocLinks) > 0 {
			if linkPtr, exists := r.tocLinks[s]; exists {
				// Dereference the pointer to get the actual link ID
				link := *linkPtr
				r.Pdf.SetLink(link, -1, -1)
				r.tracer("Text Heading", fmt.Sprintf("Set link for header '%s' with link ID: %d\n", s, link))
			} else {
				r.tracer("Text Heading", fmt.Sprintf("Header '%s' not found in links map\n", s))
			}
		}
		r.write(currentStyle, s)
	case *ast.BlockQuote:
		r.tracer("Text BlockQuote", s)
		r.multiCell(currentStyle, s)
	default:
		r.write(currentStyle, s)
	}
}

func (r *PdfRenderer) outputUnhighlightedCodeBlock(codeBlock string) {
	r.cr() // start on next line!
	r.setStyler(r.Theme.Backtick)
	if r.Theme.CodeTabWidth > 0 {
		codeBlock = strings.ReplaceAll(codeBlock, "\t", strings.Repeat(" ", r.Theme.CodeTabWidth))
	}
	r.multiCell(r.Theme.Code, codeBlock)
}

// drawCodeFill manages the light-gray background rectangles and page breaks
// for a highlighted code block.  It disables fpdf's automatic page breaking
// for the duration of the block so it can draw a correctly-sized rectangle on
// every page the block spans.  renderLine is called for each logical line and
// is responsible for writing the actual text characters.
func (r *PdfRenderer) drawCodeFill(lines []string, lineHeights []float64, renderLine func(lineN int, l string)) {
	lm, _, rm, bm := r.Pdf.GetMargins()
	pw, ph := r.Pdf.GetPageSize()
	availW := pw - lm - rm
	usableH := ph - bm

	drawBg := func(y, height float64) {
		r.setStyler(r.Theme.Code)
		r.Pdf.Rect(lm, y, availW, height, "F")
	}

	// rectHeightFrom returns the height needed to cover lines[from..] that
	// still fit on a page whose printable area starts at pageTopY.
	rectHeightFrom := func(from int, pageTopY float64) float64 {
		h := 0.0
		for i := from; i < len(lines); i++ {
			if pageTopY+h+lineHeights[i] > usableH {
				break
			}
			h += lineHeights[i]
		}
		return h
	}

	autoBreak, pbMargin := r.Pdf.GetAutoPageBreak()
	r.Pdf.SetAutoPageBreak(false, pbMargin)
	defer r.Pdf.SetAutoPageBreak(autoBreak, pbMargin)

	startX, startY := r.Pdf.GetXY()
	if r.Theme.Code.FillColor != r.Theme.BackgroundColor {
		if h := rectHeightFrom(0, startY); h > 0 {
			drawBg(startY, h)
			r.Pdf.SetXY(startX, startY)
		}
	}

	for lineN, l := range lines {
		// If this line would exceed the printable area, break to a new page
		// and draw a fresh background rectangle for the remaining lines.
		if r.Pdf.GetY()+lineHeights[lineN] > usableH {
			r.Pdf.AddPage()
			newY := r.Pdf.GetY()
			if r.Theme.Code.FillColor != r.Theme.BackgroundColor {
				if h := rectHeightFrom(lineN, newY); h > 0 {
					drawBg(newY, h)
				}
			}
			r.Pdf.SetX(lm)
		}

		renderLine(lineN, l)
		r.cr()
	}
}

func (r *PdfRenderer) processCodeblock(node ast.CodeBlock) {
	r.tracer("Codeblock", fmt.Sprintf("%v", ast.ToString(node.AsLeaf())))

	currentStyle := r.cs.peek().textStyle
	r.setStyler(currentStyle)

	if len(node.Info) < 1 {
		r.outputUnhighlightedCodeBlock(string(node.Literal))
		return
	}

	if strings.HasPrefix(string(node.Literal), "<script") && string(node.Info) == "html" {
		node.Info = []byte("javascript")
	}
	syntaxFile, lerr := syntaxhighlight.Files.ReadFile(string(node.Info) + ".yaml")
	if lerr != nil {
		r.outputUnhighlightedCodeBlock(string(node.Literal))
		return
	}
	syntaxDef, _ := highlight.ParseDef(syntaxFile)
	h := highlight.NewHighlighter(syntaxDef)
	linesWrapped := wordwrap.WrapString(string(node.Literal), 90)
	if r.Theme.CodeTabWidth > 0 {
		linesWrapped = strings.ReplaceAll(linesWrapped, "\t", strings.Repeat(" ", r.Theme.CodeTabWidth))
	}
	matches := h.HighlightString(linesWrapped)

	r.setStyler(r.Theme.Code)
	r.cr()

	lines := strings.Split(linesWrapped, "\n")
	// Trim the trailing empty element produced by a newline-terminated string.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	lineH := r.Theme.Code.Size + r.Theme.Code.Spacing
	lm, _, rm, _ := r.Pdf.GetMargins()
	pw, _ := r.Pdf.GetPageSize()
	availW := pw - lm - rm

	// Pre-compute per-line rendered heights (accounting for soft-wrapping).
	lineHeights := make([]float64, len(lines))
	for i, l := range lines {
		w := r.Pdf.GetStringWidth(l)
		if w <= 0 {
			lineHeights[i] = lineH
		} else {
			lineHeights[i] = math.Ceil(w/availW) * lineH
		}
	}

	r.drawCodeFill(lines, lineHeights, func(lineN int, l string) {
		colN := 0
		for _, c := range l {
			if group, ok := matches[lineN][colN]; ok {
				switch group {
				case highlight.Groups["default"]:
					fallthrough
				case highlight.Groups[""]:
					r.setStyler(r.Theme.Code)
				case highlight.Groups["statement"]:
					fallthrough
				case highlight.Groups["green"]:
					r.Pdf.SetTextColor(42, 170, 138)
				case highlight.Groups["identifier"]:
					fallthrough
				case highlight.Groups["blue"]:
					r.Pdf.SetTextColor(137, 207, 240)
				case highlight.Groups["preproc"]:
					r.Pdf.SetTextColor(255, 80, 80)
				case highlight.Groups["special"]:
					fallthrough
				case highlight.Groups["type.keyword"]:
					fallthrough
				case highlight.Groups["red"]:
					r.Pdf.SetTextColor(255, 80, 80)
				case highlight.Groups["constant"]:
					fallthrough
				case highlight.Groups["constant.number"]:
					fallthrough
				case highlight.Groups["constant.bool"]:
					fallthrough
				case highlight.Groups["symbol.brackets"]:
					fallthrough
				case highlight.Groups["identifier.var"]:
					fallthrough
				case highlight.Groups["cyan"]:
					r.Pdf.SetTextColor(0, 136, 163)
				case highlight.Groups["constant.specialChar"]:
					fallthrough
				case highlight.Groups["constant.string.url"]:
					fallthrough
				case highlight.Groups["constant.string"]:
					fallthrough
				case highlight.Groups["magenta"]:
					r.Pdf.SetTextColor(255, 0, 255)
				case highlight.Groups["type"]:
					fallthrough
				case highlight.Groups["symbol"]:
					fallthrough
				case highlight.Groups["symbol.operator"]:
					fallthrough
				case highlight.Groups["symbol.tag.extended"]:
					fallthrough
				case highlight.Groups["yellow"]:
					r.Pdf.SetTextColor(255, 165, 0)
				case highlight.Groups["comment"]:
					fallthrough
				case highlight.Groups["high.green"]:
					r.Pdf.SetTextColor(82, 204, 0)
				default:
					fmt.Printf("Unknown group: %s\n", group)
					r.setStyler(r.Theme.Code)
				}
			}
			r.Pdf.Write(lineH, string(c))
			colN++
		}
	})

	// Restore fill color to what the theme expects.
	r.setStyler(r.Theme.Code)
}

func (r *PdfRenderer) processList(node ast.List, entering bool) {
	kind := unordered
	if node.ListFlags&ast.ListTypeOrdered != 0 {
		kind = ordered
	}
	if node.ListFlags&ast.ListTypeDefinition != 0 {
		kind = definition
	}
	r.setStyler(r.Theme.Normal)
	if entering {
		r.tracer(fmt.Sprintf("%v List (entering)", kind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		r.Pdf.SetLeftMargin(r.cs.peek().leftMargin + r.IndentValue)
		r.tracer("... List Left Margin",
			fmt.Sprintf("set to %v", r.cs.peek().leftMargin+r.IndentValue))
		x := &containerState{
			textStyle: r.Theme.Normal, itemNumber: 0,
			listkind:   kind,
			leftMargin: r.cs.peek().leftMargin + r.IndentValue,
		}
		r.cs.push(x)
	} else {
		r.tracer(fmt.Sprintf("%v List (leaving)", kind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		r.Pdf.SetLeftMargin(r.cs.peek().leftMargin - r.IndentValue)
		r.tracer("... Reset List Left Margin",
			fmt.Sprintf("re-set to %v", r.cs.peek().leftMargin-r.IndentValue))
		r.cs.pop()
		if len(r.cs.stack) < 2 {
			r.cr()
		}
	}
}

func isListItem(node ast.Node) bool {
	_, ok := node.(*ast.ListItem)
	return ok
}

func (r *PdfRenderer) processItem(node ast.ListItem, entering bool) {
	if entering {
		r.tracer(fmt.Sprintf("%v Item (entering) #%v",
			r.cs.peek().listkind, r.cs.peek().itemNumber+1),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		r.cr() // newline before getting started
		x := &containerState{
			textStyle: r.Theme.Normal, itemNumber: r.cs.peek().itemNumber + 1,
			listkind:       r.cs.peek().listkind,
			firstParagraph: true,
			leftMargin:     r.cs.peek().leftMargin,
		}
		// add bullet or itemnumber; then set left margin for the
		// text/paragraphs in the item
		r.cs.push(x)
		if r.cs.peek().listkind == unordered {
			bulletChar := "•"
			currFontSize, _ := r.Pdf.GetFontSize()
			if node.BulletChar != 45 { // if the bullet char is not '-'
				bulletChar = "▪"
				r.Pdf.SetFont("", "", 25)
			}
			r.Pdf.CellFormat(4*r.NormalEm, r.Theme.Normal.Size+r.Theme.Normal.Spacing,
				bulletChar,
				"", 0, "RB", false, 0, "")
			r.Pdf.SetFont("", "", currFontSize)
		} else if r.cs.peek().listkind == ordered {
			r.Pdf.CellFormat(4*r.NormalEm, r.Theme.Normal.Size+r.Theme.Normal.Spacing,
				fmt.Sprintf("%v.", r.cs.peek().itemNumber),
				"", 0, "RB", false, 0, "")
		}
		// with the bullet done, now set the left margin for the text
		r.Pdf.SetLeftMargin(r.cs.peek().leftMargin + (4 * r.NormalEm))
		// set the cursor to this point
		r.Pdf.SetX(r.cs.peek().leftMargin + (4 * r.NormalEm))
	} else {
		r.tracer(fmt.Sprintf("%v Item (leaving)",
			r.cs.peek().listkind),
			fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
		// before we output the new line, reset left margin
		r.Pdf.SetLeftMargin(r.cs.peek().leftMargin)
		r.cs.parent().itemNumber++
		r.cs.pop()
	}
}

func (r *PdfRenderer) processEmph(node ast.Node, entering bool) {
	if entering {
		r.tracer("Emph (entering)", "")
		r.cs.peek().textStyle.Style += "i"
	} else {
		r.tracer("Emph (leaving)", "")
		r.cs.peek().textStyle.Style = strings.ReplaceAll(
			r.cs.peek().textStyle.Style, "i", "")
	}
}

func (r *PdfRenderer) processStrong(node ast.Node, entering bool) {
	if entering {
		r.cs.peek().textStyle.Style += "b"
		r.tracer("Strong (entering)", "")
	} else {
		r.tracer("Strong (leaving)", "")
		r.cs.peek().textStyle.Style = strings.ReplaceAll(
			r.cs.peek().textStyle.Style, "b", "")
	}
}

func (r *PdfRenderer) processLink(node ast.Link, entering bool) {
	destination := string(node.Destination)
	if entering {
		if r.InputBaseURL != "" && !strings.HasPrefix(destination, "http") {
			destination = r.InputBaseURL + "/" + strings.Replace(destination, "./", "", 1)
		}
		x := &containerState{
			textStyle: r.Theme.Link, listkind: notlist,
			leftMargin:  r.cs.peek().leftMargin,
			destination: destination,
		}
		r.cs.push(x)
		r.tracer("Link (entering)",
			fmt.Sprintf("Destination[%v] Title[%v]",
				string(node.Destination),
				string(node.Title)))
	} else {
		r.tracer("Link (leaving)", "")
		r.cs.pop()
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
	// Get the response bytes from the url
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code: " + fmt.Sprintf("HTTP %d", response.StatusCode))
	}
	// Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *PdfRenderer) processImage(node ast.Image, entering bool) {
	// while this has entering and leaving states, it doesn't appear
	// to be useful except for other markup languages to close the tag
	if entering {
		r.cr() // newline before getting started
		destination := string(node.Destination)
		tempDir := os.TempDir() + "/" + filepath.Base(os.Args[0])
		_, err := os.Stat(destination)
		if errors.Is(err, os.ErrNotExist) &&
			!strings.HasPrefix(destination, "http") &&
			r.InputBaseURL != "" &&
			!strings.HasPrefix(r.InputBaseURL, "http") {
			localPath := filepath.Join(r.InputBaseURL, destination)
			if _, lerr := os.Stat(localPath); lerr == nil {
				destination = localPath
				err = nil
			}
		}
		if errors.Is(err, os.ErrNotExist) {
			// download the image so we can use it
			var source string = destination
			if !strings.HasPrefix(destination, "http") {
				if r.InputBaseURL != "" {
					source = r.InputBaseURL + "/" + destination
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
		r.tracer("Image (entering)",
			fmt.Sprintf("Destination[%v] Title[%v]",
				destination,
				string(node.Title)))
		// following changes suggested by @sirnewton01, issue #6
		// does file exist?
		imgPath := destination
		_, err = os.Stat(imgPath)
		if err == nil {
			r.Pdf.ImageOptions(destination,
				-1, 0, 0, 0, true,
				fpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
		} else {
			r.tracer("Image (file error)", err.Error())
		}
	} else {
		r.tracer("Image (leaving)", "")
	}
}

func (r *PdfRenderer) processCode(node ast.Node) {
	r.tracer("processCode", fmt.Sprintf("%s", string(node.AsLeaf().Literal)))
	r.write(r.Theme.Normal, " ") // fix: no margin
	r.tracer("Code (entering)", "")
	r.setStyler(r.Theme.Code)
	s := string(node.AsLeaf().Literal)
	hw := r.Pdf.GetStringWidth(s)
	h := r.Theme.Code.Size
	r.Pdf.CellFormat(hw, h, s, "", 0, "C", true, 0, "")
}

func (r *PdfRenderer) processParagraph(node *ast.Paragraph, entering bool) {
	r.setStyler(r.Theme.Normal)
	if entering {
		r.tracer("Paragraph (entering)", "")
		lm, tm, rm, bm := r.Pdf.GetMargins()
		r.tracer("... Margins (left, top, right, bottom:",
			fmt.Sprintf("%v %v %v %v", lm, tm, rm, bm))
		if isListItem(node.Parent) {
			t := r.cs.peek().listkind
			if t == unordered || t == ordered || t == definition {
				if r.cs.peek().firstParagraph {
					r.tracer("First Para within a list", "breaking")
				} else {
					r.tracer("Not First Para within a list", "indent etc.")
					r.cr()
				}
			}
			return
		}
		r.cr()
	} else {
		r.tracer("Paragraph (leaving)", "")
		lm, tm, rm, bm := r.Pdf.GetMargins()
		r.tracer("... Margins (left, top, right, bottom:",
			fmt.Sprintf("%v %v %v %v", lm, tm, rm, bm))
		if isListItem(node.Parent) {
			t := r.cs.peek().listkind
			if t == unordered || t == ordered || t == definition {
				if r.cs.peek().firstParagraph {
					r.cs.peek().firstParagraph = false
				} else {
					r.tracer("Not First Para within a list", "")
					r.cr()
				}
			}
			return
		}
		r.cr()
	}
}

func (r *PdfRenderer) processBlockQuote(node ast.Node, entering bool) {
	if entering {
		r.tracer("BlockQuote (entering)", "")
		curleftmargin, _, _, _ := r.Pdf.GetMargins()
		x := &containerState{
			textStyle: r.Theme.Blockquote, listkind: notlist,
			leftMargin: curleftmargin + r.IndentValue,
		}
		r.cs.push(x)
		r.Pdf.SetLeftMargin(curleftmargin + r.IndentValue)
	} else {
		r.tracer("BlockQuote (leaving)", "")
		curleftmargin, _, _, _ := r.Pdf.GetMargins()
		r.Pdf.SetLeftMargin(curleftmargin - r.IndentValue)
		r.cs.pop()
		r.cr()
	}
}

func (r *PdfRenderer) processHeading(node ast.Heading, entering bool) {
	if entering {
		r.cr()
		switch node.Level {
		case 1:
			r.tracer("Heading (1, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H1, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		case 2:
			r.tracer("Heading (2, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H2, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		case 3:
			r.tracer("Heading (3, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H3, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		case 4:
			r.tracer("Heading (4, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H4, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		case 5:
			r.tracer("Heading (5, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H5, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		case 6:
			r.tracer("Heading (6, entering)", fmt.Sprintf("%v", ast.ToString(node.AsContainer())))
			x := &containerState{
				textStyle: r.Theme.H6, listkind: notlist,
				leftMargin: r.cs.peek().leftMargin,
			}
			r.cs.push(x)
		}
	} else {
		r.tracer("Heading (leaving)", "")
		r.cr()
		r.cs.pop()
	}
}

func (r *PdfRenderer) processHorizontalRule(node ast.Node) {
	r.tracer("HorizontalRule", "")
	if r.HorizontalRuleNewPage {
		r.Pdf.AddPage()
	} else {
		// do a newline
		r.cr()
		// get the current x and y (assume left margin in ok)
		x, y := r.Pdf.GetXY()
		// get the page margins
		lm, _, _, _ := r.Pdf.GetMargins()
		// get the page size
		w, _ := r.Pdf.GetPageSize()
		// now compute the x value of the right side of page
		newx := w - lm
		r.tracer("... From X,Y", fmt.Sprintf("%v,%v", x, y))
		r.Pdf.MoveTo(x, y)
		r.tracer("...   To X,Y", fmt.Sprintf("%v,%v", newx, y))
		r.Pdf.LineTo(newx, y)
		r.Pdf.SetLineWidth(3)
		r.Pdf.SetFillColor(200, 200, 200)
		r.Pdf.DrawPath("F")
		// another newline
		r.cr()
	}
}

func (r *PdfRenderer) processHTMLBlock(node ast.Node) {
	r.tracer("HTMLBlock", string(node.AsLeaf().Literal))
	r.cr()
	r.setStyler(r.Theme.Backtick)
	r.Pdf.CellFormat(0, r.Theme.Backtick.Size,
		string(node.AsLeaf().Literal), "", 1, "LT", true, 0, "")
	r.cr()
}

func (r *PdfRenderer) processTable(node ast.Node, entering bool) {
	if entering {
		r.tracer("Table (entering)", "")
		x := &containerState{
			textStyle: r.Theme.THeader, listkind: notlist,
			leftMargin: r.cs.peek().leftMargin,
		}
		r.cr()
		r.cs.push(x)
		fill = false
		cellwidths = r.ColumnWidths[node]
	} else {
		wSum := 0.0
		for _, w := range cellwidths {
			wSum += w
		}
		r.Pdf.CellFormat(wSum, 0, "", "T", 0, "", false, 0, "")

		r.cs.pop()
		r.tracer("Table (leaving)", "")
		r.cr()
	}
}

func (r *PdfRenderer) processTableHead(node ast.Node, entering bool) {
	if entering {
		r.tracer("TableHead (entering)", "")
		x := &containerState{
			textStyle: r.Theme.THeader, listkind: notlist,
			leftMargin: r.cs.peek().leftMargin,
		}
		r.cs.push(x)
	} else {
		r.cs.pop()
		r.tracer("TableHead (leaving)", "")
	}
}

func (r *PdfRenderer) processTableBody(node ast.Node, entering bool) {
	if entering {
		r.tracer("TableBody (entering)", "")
		x := &containerState{
			textStyle: r.Theme.TBody, listkind: notlist,
			leftMargin: r.cs.peek().leftMargin,
		}
		r.cs.push(x)
	} else {
		r.cs.pop()
		r.tracer("TableBody (leaving)", "")
		r.Pdf.Ln(-1)
	}
}

func (r *PdfRenderer) processTableRow(node ast.Node, entering bool) {
	if entering {
		r.tracer("TableRow (entering)", "")
		x := &containerState{
			textStyle: r.Theme.TBody, listkind: notlist,
			leftMargin: r.cs.peek().leftMargin,
		}
		if r.cs.peek().isHeader {
			x.textStyle = r.Theme.THeader
		}
		r.Pdf.Ln(-1)

		// initialize cell widths slice; only one table at a time!
		curdatacell = 0
		r.cs.push(x)
	} else {
		r.cs.pop()
		r.tracer("TableRow (leaving)", "")
		fill = !fill
	}
}

func (r *PdfRenderer) processTableCell(node ast.TableCell, entering bool) {
	if entering {

		r.tracer("TableCell (entering)", "")
		x := &containerState{
			textStyle: r.Theme.Normal, listkind: notlist,
			leftMargin: r.cs.peek().leftMargin,
		}
		if node.IsHeader {
			x.isHeader = true
			x.textStyle = r.Theme.THeader
			r.setStyler(r.Theme.THeader)
		} else {
			x.textStyle = r.Theme.TBody
			r.setStyler(r.Theme.TBody)
			x.isHeader = false
		}
		r.cs.push(x)
		incell = true
	} else {
		incell = false
		cs := r.cs.pop()
		currentStyle := cs.textStyle
		if cs.cellInnerStringStyle != nil {
			currentStyle = *cs.cellInnerStringStyle
		}
		s := cs.cellInnerString
		w := cellwidths[curdatacell]
		if cs.isHeader {
			h, _ := r.Pdf.GetFontSize()
			h += currentStyle.Spacing
			r.tracer("... table header cell",
				fmt.Sprintf("Width=%v, height=%v", w, h))

			r.Pdf.CellFormat(w, h, s, "1", 0, "C", true, 0, "")
		} else {
			h := currentStyle.Size + currentStyle.Spacing
			r.Pdf.CellFormat(w, h, s, "LR", 0, "", fill, 0, "")
		}
		r.tracer("TableCell (leaving)", "")
		curdatacell++
	}
}
