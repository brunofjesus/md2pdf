// Package highlight provides syntax highlighting for code blocks in markdown
// documents.
package highlight

import "embed"

// Files embeds the syntax definition YAML files used for code block
// syntax highlighting. The files were sourced from:
// https://github.com/jessp01/gohighlight/tree/2b769d0/syntax_files
//
//go:embed *.yaml *.yml
var Files embed.FS
