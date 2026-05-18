package renderer

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gomarkdown/markdown/parser"
)

func testit(t *testing.T, inputf string, gohighlight bool) {
	t.Helper()

	inputd := "../../testdata/"
	input := path.Join(inputd, inputf)

	tracerfile := path.Join(inputd, strings.TrimSuffix(path.Base(input), path.Ext(input)))
	tracerfile += ".log"

	file, err := os.Open(input) //nolint:gosec
	if err != nil {
		t.Errorf("%v:%v", input, err)
	}
	defer func() { _ = file.Close() }()

	var opts []RenderOption
	if gohighlight {
		opts = []RenderOption{WithHorizontalRuleAsNewPage()}
	}

	params := PdfRendererParams{
		Title:           "",
		Orientation:     "",
		PageSize:        "",
		TracerFile:      tracerfile,
		Opts:            opts,
		Theme:           LIGHT,
		CustomThemeFile: "",
	}

	r := NewPdfRenderer(params)
	if absInput, err := filepath.Abs(input); err == nil {
		r.InputBaseURL = filepath.Dir(absInput)
	}

	r.Extensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
		parser.Autolink | parser.Strikethrough | parser.SpaceHeadings |
		parser.HeadingIDs | parser.BackslashLineBreak | parser.DefinitionLists

	err = r.Process(file)
	if err != nil {
		t.Error(err)
	}
}

func TestTables(t *testing.T) {
	t.Parallel()
	testit(t, "Tables.text", false)
}

func TestMarkdownDocumenationBasic(t *testing.T) {
	t.Parallel()
	testit(t, "Markdown Documentation - Basics.text", false)
}

func TestMarkdownDocumenationSyntax(t *testing.T) {
	t.Parallel()
	testit(t, "syntax.md", false)
}

func TestMarkdownDocumenationColourSyntax(t *testing.T) {
	t.Parallel()
	testit(t, "syntax_highlighting.md", true)
}

func TestImage(t *testing.T) {
	t.Parallel()
	testit(t, "Image.text", false)
}

func TestAutoLinks(t *testing.T) {
	t.Parallel()
	testit(t, "Auto links.text", false)
}

func TestAmpersandEncoding(t *testing.T) {
	t.Parallel()
	testit(t, "Amps and angle encoding.text", false)
}

func TestInlineLinks(t *testing.T) {
	t.Parallel()
	testit(t, "Links, inline style.text", false)
}

func TestLists(t *testing.T) {
	t.Parallel()
	testit(t, "Ordered and unordered lists.md", false)
}

func TestStringEmph(t *testing.T) {
	t.Parallel()
	testit(t, "Strong and em together.text", false)
}

func TestTabs(t *testing.T) {
	t.Parallel()
	testit(t, "Tabs.text", false)
}

func TestBackslashEscapes(t *testing.T) {
	t.Parallel()
	testit(t, "Backslash escapes.text", false)
}

func TestBackquotes(t *testing.T) {
	t.Parallel()
	testit(t, "Blockquotes with code blocks.text", false)
}

func TestCodeBlocks(t *testing.T) {
	t.Parallel()
	testit(t, "Code Blocks.text", false)
}

func TestCodeSpans(t *testing.T) {
	t.Parallel()
	testit(t, "Code Spans.text", false)
}

func TestHardWrappedPara(t *testing.T) {
	t.Parallel()
	testit(t, "Hard-wrapped paragraphs with list-like lines no empty line before block.text", false)
}

func TestHardWrappedPara2(t *testing.T) {
	t.Parallel()
	testit(t, "Hard-wrapped paragraphs with list-like lines.text", false)
}

func TestHorizontalRules(t *testing.T) {
	t.Parallel()
	testit(t, "Horizontal rules.text", false)
}

func TestInlineHtmlSimple(t *testing.T) {
	t.Parallel()
	testit(t, "Inline HTML (Simple).text", false)
}

func TestInlineHtmlAdvanced(t *testing.T) {
	t.Parallel()
	testit(t, "Inline HTML (Advanced).text", false)
}

func TestInlineHtmlComments(t *testing.T) {
	t.Parallel()
	testit(t, "Inline HTML comments.text", false)
}

func TestTitleWithQuotes(t *testing.T) {
	t.Parallel()
	testit(t, "Literal quotes in titles.text", false)
}

func TestNestedBlockquotes(t *testing.T) {
	t.Parallel()
	testit(t, "Nested blockquotes.text", false)
}

func TestLinksReference(t *testing.T) {
	t.Parallel()
	testit(t, "Links, reference style.text", false)
}

func TestLinksShortcut(t *testing.T) {
	t.Parallel()
	testit(t, "Links, shortcut references.text", false)
}

func TestTidyness(t *testing.T) {
	t.Parallel()
	testit(t, "Tidyness.text", false)
}
