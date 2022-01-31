package graph

import (
	"github.com/spatial-go/geoos/algorithm"
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// Represents the linework for edges in the topology derived from (up to) two
// parent geometries. An edge may be the result of the merging of two or more
// edges which have the same linework (although possibly different orientations).
// In this case the topology information is derived from the merging of the
// information in the source edges. Merged edges can occur in the following situations
// Due to coincident edges of polygonal or linear geometries.
// Due to topology collapse caused by snapping or rounding of polygonal geometries.
// The source edges may have the same parent geometry, or different ones, or a mix of the two.
type Edge struct {
	pts         []matrix.Matrix
	aDim        int // default=-1
	aDepthDelta int
	aIsHole     bool
	bDim        int // default=-1
	bDepthDelta int
	bIsHole     bool
}

// IsCollapsed Tests if the given point sequence is a collapsed line. A collapsed edge has fewer than two distinct points.
// Params:
//		pts – the point sequence to check
// Returns:
//		true if the points form a collapsed line
func (e *Edge) IsCollapsed(pts []matrix.Matrix) bool {
	if len(pts) < 2 {
		return true
	}
	if pts[0].Equals(pts[1]) {
		return true
	}
	if len(pts) > 2 {
		if pts[len(pts)-1].Equals(pts[len(pts)-2]) {
			return true
		}
	}
	return false
}

// NewEdge...
func NewEdge(pts []matrix.Matrix, info *EdgeSourceInfo) *Edge {
	var edge = Edge{pts: pts}
	edge.copyInfo(info)
	return &edge
}

// copyInfo...
func (e *Edge) copyInfo(info *EdgeSourceInfo) {
	if info.index == 0 {
		e.aDim = info.dim
		e.aIsHole = info.isHole
		e.aDepthDelta = info.depthDelta
	} else {
		e.bDim = info.dim
		e.bIsHole = info.isHole
		e.bDepthDelta = info.depthDelta
	}
}

// direction...
func (e *Edge) direction() (bool, error) {
	pts := e.pts
	if len(pts) < 2 {
		// Edge must have >= 2 points
		return false, algorithm.ErrEdgeTooFewPoint
	}
	p0 := pts[0]
	p1 := pts[1]

	pn0 := pts[len(pts)-1]
	pn1 := pts[len(pts)-2]

	cmp := 0
	cmp0, err := p0.Compare(pn0)
	if err != nil {
		return false, err
	}
	if cmp0 != 0 {
		cmp = cmp0
	}

	if cmp == 0 {
		cmp1, err := p1.Compare(pn1)
		if err != nil {
			return false, err
		}
		if cmp1 != 0 {
			cmp = cmp1
		}
	}

	if cmp == 0 {
		// Edge direction cannot be determined because endpoints are equal
		return false, algorithm.ErrEdgeEndPointEqual
	}

	return cmp == -1, nil
}

// merge Merges an edge into this edge, updating the topology info accordingly.
func (e *Edge) merge(edge *Edge) {
	/**
	 * Marks this as a shell edge if any contributing edge is a shell.
	 * Update hole status first, since it depends on edge dim
	 */
	e.aIsHole = e.isHoleMerged(0, e, edge)
	e.bIsHole = e.isHoleMerged(1, e, edge)

	if edge.aDim > e.aDim {
		e.aDim = edge.aDim
	}
	if edge.bDim > e.bDim {
		e.bDim = edge.bDim
	}

	relDir := e.relativeDirection(edge)
	var flipFactor int
	if relDir {
		flipFactor = 1
	} else {
		flipFactor = -1
	}

	e.aDepthDelta += flipFactor * edge.aDepthDelta
	e.bDepthDelta += flipFactor * edge.bDepthDelta
}

// isHoleMerged...
func (e *Edge) isHoleMerged(geomIndex int, edge1, edge2 *Edge) bool {
	// TODO: this might be clearer with tri-state logic for isHole?
	isShell1 := edge1.isShell(geomIndex)
	isShell2 := edge2.isShell(geomIndex)
	isShellMerged := isShell1 || isShell2
	// flip since isHole is stored
	return !isShellMerged
}

// isShell Tests whether the edge is part of a shell in the given geometry.
// This is only the case if the edge is a boundary.
// Params:
//		geomIndex – the index of the geometry
// Returns:
//		true if this edge is a boundary and part of a shell
func (e *Edge) isShell(geomIndex int) bool {
	if geomIndex == 0 {
		return e.aDim == DIM_BOUNDARY && !e.aIsHole
	}
	return e.bDim == DIM_BOUNDARY && !e.bIsHole
}

// relativeDirection Compares two coincident edges to determine whether they have
// the same or opposite direction.
// Params:
//		edge2 – an edge
// Returns:
//		true if the edges have the same direction, false if not
func (e *Edge) relativeDirection(edge2 *Edge) bool {
	// assert: the edges match (have the same coordinates up to direction)
	if !e.pts[0].Equals(edge2.pts[0]) {
		return false
	}
	if !e.pts[1].Equals(edge2.pts[1]) {

		return false
	}
	return true
}
