package colors

// Color is a RGB set of ints; for a nice picker
// see https://www.w3schools.com/colors/colors_picker.asp
type Color struct {
	Red, Green, Blue int
}

func New(red, green, blue int) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}
