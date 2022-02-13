package graph

import "github.com/spatial-go/geoos/space"

// Extracts Point resultants from an overlay graph created by an Intersection operation
// between non-Point inputs. Points may be created during intersection if lines or areas
// touch one another at single points. Intersection is the only overlay operation which
// can result in Points from non-Point inputs.
// Overlay operations where one or more inputs are Points are handled via a different code path.
type IntersectionPointBuilder struct {
	graph  *OverlayGraph
	points []space.Point

	// Controls whether lines created by area topology collapses to participate in the result
	// computation. True provides the original JTS semantics.
	isAllowCollapseLines bool // = ! OverlayNG.STRICT_MODE_DEFAULT
}

// NewIntersectionPointBuilder...
func NewIntersectionPointBuilder(graph *OverlayGraph, strictMode bool) *IntersectionPointBuilder {
	return &IntersectionPointBuilder{
		graph:                graph,
		isAllowCollapseLines: !strictMode,
	}
}

// setStrictMode...
func (i *IntersectionPointBuilder) setStrictMode(isStrictMode bool) {
	i.isAllowCollapseLines = !isStrictMode
}

// getPoints...
func (i *IntersectionPointBuilder) getPoints() []space.Point {
	i.addResultPoints()
	return i.points
}

// addResultPoints...
func (i *IntersectionPointBuilder) addResultPoints() {
	for _, nodeEdge := range i.graph.nodeMap {
		if i.isResultPoint(nodeEdge) {
			pt := nodeEdge.origin
			i.points = append(i.points, space.Point(pt))
		}
	}
}

// isResultPoint Tests if a node is a result point. This is the case if the node
// is incident on edges from both inputs, and none of the edges are themselves in the result.
// Params:
// 		nodeEdge – an edge originating at the node
// Returns:
//		true if this node is a result point
func (i *IntersectionPointBuilder) isResultPoint(nodeEdge *OverlayEdge) bool {
	isEdgeOfA := false
	isEdgeOfB := false

	edge := nodeEdge
	for edge != nodeEdge {
		if edge.isInResult() {
			return false
		}
		label := edge.label
		isEdgeOfA = isEdgeOfA || i.isEdgeOf(label, 0)
		isEdgeOfB = isEdgeOfB || i.isEdgeOf(label, 1)
		edge = edge.oNextOE()
	}
	isNodeInBoth := isEdgeOfA && isEdgeOfB
	return isNodeInBoth
}

// isEdgeOf...
func (i *IntersectionPointBuilder) isEdgeOf(label *OverlayLabel, index int) bool {
	if !i.isAllowCollapseLines && label.isBoundaryCollapse() {
		return false
	}
	return label.isBoundary(index) || label.isLineByIndex(index)
}
