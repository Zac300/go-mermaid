package domain

// Node is a vertex in a flowchart.
type Node struct {
	// ID is the unique identifier from the source (e.g. "A").
	ID string
	// Label is the display text; defaults to ID when not given.
	Label string
	// Shape is the node outline style.
	Shape Shape

	// Pos is the laid-out top-left position. Zero until layout runs.
	Pos Point
	// Size is the laid-out box size. Zero until layout runs.
	Size Size
}

// Center returns the geometric center of the node's box.
func (n *Node) Center() Point {
	return Point{
		X: n.Pos.X + n.Size.W/2,
		Y: n.Pos.Y + n.Size.H/2,
	}
}
