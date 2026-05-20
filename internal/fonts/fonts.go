// Package fonts contains the default fonts used by the application.
// These are embedded in the binary using Go's embed package.
// The fonts are licensed under the SIL Open Font License and are sourced from
// the Liberation Fonts project on GitHub.
package fonts

import _ "embed"

// Liberation Sans - SIL Open Font License
// https://github.com/liberationfonts/liberation-fonts

// LiberationSansRegular font
//
//go:embed LiberationSans-Regular.ttf
var LiberationSansRegular []byte

// LiberationSansBold font
//
//go:embed LiberationSans-Bold.ttf
var LiberationSansBold []byte

// LiberationSansItalic font
//
//go:embed LiberationSans-Italic.ttf
var LiberationSansItalic []byte

// LiberationSansBoldItalic font
//
//go:embed LiberationSans-BoldItalic.ttf
var LiberationSansBoldItalic []byte

// Liberation Mono - SIL Open Font License
// https://github.com/liberationfonts/liberation-fonts

// LiberationMonoRegular font
//
//go:embed LiberationMono-Regular.ttf
var LiberationMonoRegular []byte

// LiberationMonoBold font
//
//go:embed LiberationMono-Bold.ttf
var LiberationMonoBold []byte

// LiberationMonoItalic font
//
//go:embed LiberationMono-Italic.ttf
var LiberationMonoItalic []byte

// LiberationMonoBoldItalic font
//
//go:embed LiberationMono-BoldItalic.ttf
var LiberationMonoBoldItalic []byte
