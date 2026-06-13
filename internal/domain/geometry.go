package domain

// Point is a 2D coordinate in SVG user units (pixels).
type Point struct {
	X, Y float64
}

// Size is a width/height pair in SVG user units (pixels).
type Size struct {
	W, H float64
}

// Rect is an axis-aligned box defined by its top-left corner and size.
type Rect struct {
	Min  Point
	Size Size
}
