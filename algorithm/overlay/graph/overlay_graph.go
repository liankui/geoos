package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// A planar graph of edges, representing the topology resulting from
// an overlay operation. Each source edge is represented by a pair of
// OverlayEdges, with opposite (symmetric) orientation. The pair of
// OverlayEdges share the edge coordinates and a single OverlayLabel.
type OverlayGraph struct {
	edges   []*OverlayEdge
	nodeMap map[interface{}]matrix.LineMatrix
}

// addEdge Adds a new edge to this graph, for the given linework and
// topology information. A pair of OverlayEdges with opposite (symmetric)
// orientation is added, sharing the same OverlayLabel.
// Params:
//		pts – the edge vertices
//		label – the edge topology information
// Returns:
//		the created graph edge with same orientation as the linework
func (o *OverlayGraph) addEdge(pts matrix.LineMatrix, label OverlayLabel) *OverlayEdge {
	var overlayEdge OverlayEdge
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
		nodeEdge.
	}
}