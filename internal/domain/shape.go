package domain

// Shape is a node outline style. The zero value is ShapeRect.
type Shape string

const (
	// ShapeRect is a rectangle: A[label].
	ShapeRect Shape = "rect"
	// ShapeRound is a rounded rectangle: A(label).
	ShapeRound Shape = "round"
	// ShapeStadium is a stadium/pill: A([label]).
	ShapeStadium Shape = "stadium"
	// ShapeCircle is a circle: A((label)).
	ShapeCircle Shape = "circle"
	// ShapeDiamond is a rhombus/decision: A{label}.
	ShapeDiamond Shape = "diamond"
)
