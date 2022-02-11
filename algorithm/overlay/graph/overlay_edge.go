package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/edgegraph"
)

type OverlayEdge struct {
	*edgegraph.HalfEdge
	origin, dirPt matrix.Matrix
	/*
	 * true indicates direction is forward along segString false is reverse direction
	 * The label must be interpreted accordingly.
	*/
	direction      bool
	pts            []matrix.Matrix
	label          *OverlayLabel
	isInResultArea bool
	isInResultLine bool
	isVisited      bool

	/**
	 * Link to next edge in the result ring.
	 * The origin of the edge is the dest of this edge.
	 */
	nextResultEdge    *OverlayEdge
	edgeRing          *OverlayEdgeRing
	maxEdgeRing       *MaximalEdgeRing
	nextResultMaxEdge *OverlayEdge
}

// getCoordinatesOriented...
func (o *OverlayEdge) getCoordinatesOriented() []matrix.Matrix {
	if o.direction {
		return o.pts
	}
	var co []matrix.Matrix
	copy(co, o.pts)
	co = o.reverse(co)
	return co
}

// reverse...
func (o *OverlayEdge) reverse(coord []matrix.Matrix) []matrix.Matrix {
	if len(coord) <= 1 {
		return coord
	}
	last := len(coord) - 1
	mid := last / 2
	for i := 0; i <= mid; i++ {
		tmp := coord[i]
		coord[i] = coord[last-i]
		coord[last-i] = tmp
	}
	return coord
}

// createEdgePair...
func (o *OverlayEdge) createEdgePair(pts []matrix.Matrix, lbl *OverlayLabel) *OverlayEdge {
	e0 := o.createEdge(pts, lbl, true)
	e1 := o.createEdge(pts, lbl, false)
	edgegraph.HalfEdgerLink(e0, e1) // todo 待验证
	return e0
}

// createEdge...
func (o *OverlayEdge) createEdge(pts []matrix.Matrix, lbl *OverlayLabel, direction bool) *OverlayEdge {
	var origin, dirPt matrix.Matrix
	if direction {
		origin = pts[0]
		dirPt = pts[1]
	} else {
		ilast := len(pts) - 1
		origin = pts[ilast]
		dirPt = pts[ilast-1]
	}
	return &OverlayEdge{
		origin:    origin,
		dirPt:     dirPt,
		direction: direction,
		pts:       pts,
		label:     lbl,
	}
}

// symOE Gets the symmetric pair edge of this edge.
// Returns:
//		the symmetric pair edge
func (o *OverlayEdge) symOE() *OverlayEdge {
	sym := edgegraph.HalfEdgerSym(o)
	return sym.(*OverlayEdge)
}

// oNextOE Gets the next edge CCW around the origin of this edge, with the same origin.
// If the origin vertex has degree 1 then this is the edge itself.
// Returns:
//		the next edge around the origin
func (o *OverlayEdge) oNextOE() *OverlayEdge {
	oNext := edgegraph.HalfEdgerONext(o)
	return oNext.(*OverlayEdge)
}

// addCoordinates Adds the coordinates of this edge to the given list, in the direction
// of the edge. Duplicate coordinates are removed (which means that this is safe to use
// for a path of connected edges in the topology graph).
// Params:
//		coords – the coordinate list to add to
func (o *OverlayEdge) addCoordinates(coords []matrix.Matrix) {
	var coordList matrix.CoordinateList = coords
	isFirstEdge := len(coordList) > 0
	if o.direction {
		startIndex := 1
		if isFirstEdge {
			startIndex = 0
		}
		for i := startIndex; i < len(o.pts); i++ {
			coordList.AddToEndList(o.pts[i], false)
		}
	} else {
		startIndex := len(o.pts) - 2
		if isFirstEdge {
			startIndex = len(o.pts) - 1
		}
		for i := startIndex; i >= 0; i-- {
			coordList.AddToEndList(o.pts[i], false)
		}
	}
}

// isInResult...
func (o *OverlayEdge) isInResult() bool {
	return o.isInResultArea || o.isInResultLine
}

// isInResultEither...
func (o *OverlayEdge) isInResultEither() bool {
	return o.isInResult() || o.symOE().isInResult()
}

// markInResultLine...
func (o *OverlayEdge) markInResultLine() {
	o.isInResultLine = true
	o.symOE().isInResultLine = true
}

// markVisited...
func (o *OverlayEdge) markVisited() {
	o.isVisited = true
}

// markVisitedBoth...
func (o *OverlayEdge) markVisitedBoth() {
	o.markVisited()
	o.symOE().markVisited()
}
