package graph

import "strings"

type Color string

const (
	Blue   Color = "/blues8/7"
	Green  Color = "/bugn8/7"
	Teal   Color = "/brbg8/7"
	Red    Color = "/orrd8/7"
	Grey   Color = "/greys8/7"
	Purple Color = "/bupu8/7"

	DarkBlue   Color = "/blues8/8"
	DarkGreen  Color = "/bugn8/8"
	DarkTeal   Color = "/brbg8/8"
	DarkRed    Color = "/orrd8/8"
	DarkGrey   Color = "/greys8/8"
	DarkPurple Color = "/bupu8/8"

	White Color = "white"
)

type Shape string

const (
	Box      Shape = "box"
	Diamond  Shape = "diamond"
	Cylinder Shape = "cylinder"
	Hexagon  Shape = "hexagon"
	Record   Shape = "record"
)

type Style string

const (
	Rounded Style = "rounded"
	Filled  Style = "filled"
	Bold    Style = "bold"
	Dotted  Style = "dotted"
	Dashed  Style = "dashed"
)

type Palette struct {
	ColorFill   Color
	ColorBorder Color
	ColorFont   Color
}

type Decorations struct {
	styles  []Style
	shape   Shape
	palette Palette
}

// strings.Join replacement for []Style
func JoinStyles(styles []Style, sep string) string {
	var b strings.Builder

	for i, s := range styles {
		if i != 0 {
			b.WriteString(sep)
		}
		b.WriteString(string(s))
	}

	return b.String()
}
