package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/edgegraph"
)

type OverlayEdge struct {
	*edgegraph.HalfEdge
	origin, dirPt matrix.Matrix
	// true indicates direction is forward along segString false is
	// reverse direction The label must be interpreted accordingly.
	direction bool
	pts       matrix.LineMatrix
	label     OverlayLabel
}

// createEdgePair...
func (o *OverlayEdge) createEdgePair(pts matrix.LineMatrix, lbl OverlayLabel) edgegraph.IHalfEdge {
	e0 := o.createEdge(pts, lbl, true)
	e1 := o.createEdge(pts, lbl, false)
	e0.(*edgegraph.HalfEdge).Link(e1) // todo 转换问题
	return e0
}

// createEdge...
func (o *OverlayEdge) createEdge(pts matrix.LineMatrix, lbl OverlayLabel, direction bool) edgegraph.IHalfEdge {
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
	//return o.HalfEdge.Sym() // todo 结构
	return nil
}

// directionPt...
func (o *OverlayEdge) DirectionPt() matrix.Matrix {
	return o.dirPt
}

// ToString
func (o *OverlayEdge) ToString() string {

	return ""
}