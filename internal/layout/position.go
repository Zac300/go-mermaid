package layout

import "github.com/Zac300/go-mermaid/internal/domain"

// position assigns pixel coordinates to each node from its rank and order,
// then returns the overall diagram bounds. The cross axis is centered per
// layer. Direction controls which screen axis the ranks grow along.
func position(g *domain.Graph, layers [][]string, opts Options) (width, height float64) {
	horizontal := g.Direction == domain.LeftRight || g.Direction == domain.RightLeft

	// Extent of each layer along the rank axis (max node thickness) and
	// the running offset where each layer starts.
	rankExtent := make([]float64, len(layers))
	for r, ids := range layers {
		for _, id := range ids {
			n := g.NodeByID(id)
			thick := n.Size.H
			if horizontal {
				thick = n.Size.W
			}
			if thick > rankExtent[r] {
				rankExtent[r] = thick
			}
		}
	}

	// Cross-axis total per layer, to center layers against the widest one.
	layerCross := make([]float64, len(layers))
	var maxCross float64
	for r, ids := range layers {
		var sum float64
		for i, id := range ids {
			n := g.NodeByID(id)
			cross := n.Size.W
			if horizontal {
				cross = n.Size.H
			}
			sum += cross
			if i > 0 {
				sum += opts.NodeSep
			}
		}
		layerCross[r] = sum
		if sum > maxCross {
			maxCross = sum
		}
	}

	var rankOffset float64
	for r, ids := range layers {
		cross := (maxCross - layerCross[r]) / 2 // center this layer
		for _, id := range ids {
			n := g.NodeByID(id)
			cw, ch := n.Size.W, n.Size.H
			if horizontal {
				n.Pos = domain.Point{X: rankOffset, Y: cross}
				cross += ch + opts.NodeSep
			} else {
				n.Pos = domain.Point{X: cross, Y: rankOffset}
				cross += cw + opts.NodeSep
			}
		}
		rankOffset += rankExtent[r] + opts.RankSep
	}

	// Bounds.
	for _, n := range g.Nodes {
		if r := n.Pos.X + n.Size.W; r > width {
			width = r
		}
		if b := n.Pos.Y + n.Size.H; b > height {
			height = b
		}
	}
	return width, height
}
