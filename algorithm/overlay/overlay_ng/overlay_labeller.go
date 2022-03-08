package overlay_ng

// mplements the logic to compute the full labeling for the edges in an OverlayGraph.
type OverlayLabeller struct {
	graph         *OverlayGraph
	inputGeometry *InputGeometry
	edges         []*OverlayEdge
}

// NewOverlayLabeller...
func NewOverlayLabeller(graph *OverlayGraph, inputGeometry *InputGeometry) *OverlayLabeller {
	return &OverlayLabeller{
		graph: graph,
		inputGeometry: inputGeometry,
		edges: graph.edges,
	}
}
