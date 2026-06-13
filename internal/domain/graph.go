// Package domain holds the pure diagram model: the types every other
// stage produces or consumes. It has no I/O and no third-party
// dependencies, so it can be reasoned about and tested in isolation.
package domain

// Direction is the flow direction of a flowchart.
type Direction string

const (
	// TopBottom lays ranks out top to bottom (graph TD / graph TB).
	TopBottom Direction = "TB"
	// BottomTop lays ranks out bottom to top (graph BT).
	BottomTop Direction = "BT"
	// LeftRight lays ranks out left to right (graph LR).
	LeftRight Direction = "LR"
	// RightLeft lays ranks out right to left (graph RL).
	RightLeft Direction = "RL"
)

// Graph is a parsed flowchart, independent of layout or rendering.
// Coordinates are not set until the layout stage populates them.
type Graph struct {
	Direction Direction
	Nodes     []*Node
	Edges     []*Edge
}

// NodeByID returns the node with the given id, or nil if absent.
func (g *Graph) NodeByID(id string) *Node {
	for _, n := range g.Nodes {
		if n.ID == id {
			return n
		}
	}
	return nil
}
