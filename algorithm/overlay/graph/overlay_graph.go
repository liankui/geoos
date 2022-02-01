package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/edgegraph"
)

// A planar graph of edges, representing the topology resulting from
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
	o.insert(e.(*OverlayEdge).symOE())
	return e.(*OverlayEdge)
}

// insert Inserts a single half-edge into the graph. The sym edge must also be inserted.
// Params:
//		e – the half-edge to insert
func (o *OverlayGraph) insert(e edgegraph.IHalfEdge) {
	o.edges = append(o.edges, e.(*OverlayEdge))
	/**
	 * If the edge origin node is already in the graph,
	 * insert the edge into the star of edges around the node.
	 * Otherwise, add a new node for the origin.
	 */
	nodeEdge, ok := o.nodeMap[e.(*OverlayEdge).origin]
	if ok {
		nodeEdge.Insert(e.(*edgegraph.HalfEdge)) // todo 方式不太对
	} else {
		o.nodeMap[e.(*OverlayEdge).origin] = e.(*OverlayEdge)
	}
}
