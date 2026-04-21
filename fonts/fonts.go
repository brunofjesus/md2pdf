package fonts

import _ "embed"

// Liberation Sans - SIL Open Font License
// https://github.com/liberationfonts/liberation-fonts

//go:embed LiberationSans-Regular.ttf
var LiberationSansRegular []byte

//go:embed LiberationSans-Bold.ttf
var LiberationSansBold []byte

//go:embed LiberationSans-Italic.ttf
var LiberationSansItalic []byte

//go:embed LiberationSans-BoldItalic.ttf
var LiberationSansBoldItalic []byte
