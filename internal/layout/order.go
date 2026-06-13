package layout

import "github.com/Zac300/go-mermaid/internal/domain"

// orderLayers groups node IDs by rank, preserving the input order within
// each layer. This is the seed order; crossing minimization (median /
// barycenter heuristic) is a planned refinement that will reorder these.
func orderLayers(g *domain.Graph, ranks map[string]int) [][]string {
	maxRank := 0
	for _, r := range ranks {
		if r > maxRank {
			maxRank = r
		}
	}
	layers := make([][]string, maxRank+1)
	for _, n := range g.Nodes {
		r := ranks[n.ID]
		layers[r] = append(layers[r], n.ID)
	}
	return layers
}
