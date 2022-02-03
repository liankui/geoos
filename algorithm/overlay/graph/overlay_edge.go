package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/edgegraph"
)

const (
	isInResultArea = false
	isInResultLine = false
	isVisited      = false
)

type OverlayEdge struct {
	*edgegraph.HalfEdge
	origin, dirPt matrix.Matrix
	// true indicates direction is forward along segString false is
	// reverse direction The label must be interpreted accordingly.
	direction bool
	pts       []matrix.Matrix
	label     *OverlayLabel
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
func (o *OverlayEdge) createEdgePair(pts []matrix.Matrix, lbl *OverlayLabel) edgegraph.IHalfEdge {
	e0 := o.createEdge(pts, lbl, true)
	e1 := o.createEdge(pts, lbl, false)
	e0.(*edgegraph.HalfEdge).Link(e1) // todo 转换问题
	return e0
}

// createEdge...
func (o *OverlayEdge) createEdge(pts []matrix.Matrix, lbl *OverlayLabel, direction bool) edgegraph.IHalfEdge {
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
