package overlay_ng

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/overlay_ng/edgegraph"
)

// OverlayGraph A planar graph of edges, representing the topology resulting from
// an overlay operation. Each source edge is represented by a pair of
// OverlayEdges, with opposite (symmetric) orientation. The pair of
// OverlayEdges share the edge coordinates and a single OverlayLabel.
type OverlayGraph struct {
	edges   []*OverlayEdge
	nodeMap map[interface{}]*OverlayEdge
}

// addEdge Adds a new edge to this graph, for the given linework and
// topology information. A pair of OverlayEdges with opposite (symmetric)
// orientation is added, sharing the same OverlayLabel.
// Params:
//		pts – the edge vertices
//		label – the edge topology information
// Returns:
//		the created graph edge with same orientation as the linework
func (o *OverlayGraph) addEdge(pts []matrix.Matrix, label *OverlayLabel) *OverlayEdge {
	overlayEdge := new(OverlayEdge)
	e := overlayEdge.createEdgePair(pts, label)
	o.insert(e)
	o.insert(e.symOE())
	return e
}

// insert Inserts a single half-edge into the graph. The sym edge must also be inserted.
// Params:
//		e – the half-edge to insert
func (o *OverlayGraph) insert(e *OverlayEdge) {
	o.edges = append(o.edges, e)
	/**
	 * If the edge origin node is already in the graph,
	 * insert the edge into the star of edges around the node.
	 * Otherwise, add a new node for the origin.
	 */
	nodeEdge, ok := o.nodeMap[e.origin]
	if ok {
		edgegraph.HalfEdgerInsert(nodeEdge, e)
	} else {
		o.nodeMap[e.origin] = e
	}
}

// getResultAreaEdges Gets the representative edges marked as being in the result area.
func (o *OverlayGraph) getResultAreaEdges() []*OverlayEdge {
	resultEdges := make([]*OverlayEdge, 0)
	for _, edge := range o.edges {
		if edge.isInResultArea {
			resultEdges = append(resultEdges, edge)
		}
	}
	return resultEdges
}
