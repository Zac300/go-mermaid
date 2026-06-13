// Package layout assigns coordinates to a domain.Graph using a layered
// (Sugiyama-style) approach: make the graph acyclic, rank nodes into
// layers, order within layers to reduce crossings, then assign pixel
// positions. v0 uses longest-path ranking; network-simplex can replace
// it behind the same interface.
package layout

import (
	"github.com/Zac300/go-mermaid/internal/domain"
)

// Options tunes spacing and text metrics used during layout.
type Options struct {
	NodeSep  float64 // gap between nodes within a layer
	RankSep  float64 // gap between layers
	FontSize float64 // used to estimate node sizes
}

// Result is a laid-out graph plus its overall bounds.
type Result struct {
	Graph  *domain.Graph
	Width  float64
	Height float64
}

// Compute lays out g in place and returns the result. The input graph's
// nodes and edges are mutated with positions and routed points.
func Compute(g *domain.Graph, opts Options) (*Result, error) {
	sizeNodes(g, opts)

	reversed := makeAcyclic(g)
	ranks := assignRanks(g)
	layers := orderLayers(g, ranks)

	w, h := position(g, layers, opts)
	routeEdges(g)
	restoreReversed(g, reversed)

	return &Result{Graph: g, Width: w, Height: h}, nil
}

// sizeNodes estimates a box size for each node from its label and font size.
func sizeNodes(g *domain.Graph, opts Options) {
	const padX, padY = 20.0, 14.0
	charW := opts.FontSize * 0.6
	for _, n := range g.Nodes {
		label := n.Label
		if label == "" {
			label = n.ID
		}
		w := float64(len([]rune(label)))*charW + padX*2
		h := opts.FontSize + padY*2
		if n.Shape == domain.ShapeCircle {
			if w < h {
				w = h
			}
			h = w
		}
		n.Size = domain.Size{W: w, H: h}
	}
}

// routeEdges sets a straight two-point polyline between node centers.
// Proper orthogonal/spline routing is a later refinement.
func routeEdges(g *domain.Graph) {
	for _, e := range g.Edges {
		from := g.NodeByID(e.From)
		to := g.NodeByID(e.To)
		if from == nil || to == nil {
			continue
		}
		e.Points = []domain.Point{from.Center(), to.Center()}
	}
}
