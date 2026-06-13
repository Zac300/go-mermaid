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
	// ShapeSubroutine is a framed rectangle: A[[label]].
	ShapeSubroutine Shape = "subroutine"
	// ShapeCylinder is a database cylinder: A[(label)].
	ShapeCylinder Shape = "cylinder"
	// ShapeHexagon is a hexagon: A{{label}}.
	ShapeHexagon Shape = "hexagon"
	// ShapeParallelogram slants right: A[/label/].
	ShapeParallelogram Shape = "parallelogram"
	// ShapeParallelogramAlt slants left: A[\label\].
	ShapeParallelogramAlt Shape = "parallelogram_alt"
	// ShapeTrapezoid is wider at the bottom: A[/label\].
	ShapeTrapezoid Shape = "trapezoid"
	// ShapeTrapezoidAlt is wider at the top: A[\label/].
	ShapeTrapezoidAlt Shape = "trapezoid_alt"
)
