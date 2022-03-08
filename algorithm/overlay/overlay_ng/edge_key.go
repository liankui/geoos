package overlay_ng

import "github.com/spatial-go/geoos/algorithm/matrix"

// A key for sorting and comparing edges in a noded arrangement. Relies on the
// fact that in a correctly noded arrangement edges are identical (up to direction)
// if they have their first segment in common.
type EdgeKey struct {
	p0x float64
	p0y float64
	p1x float64
	p1y float64
}

// create...
func (e *EdgeKey) create(edge *Edge) {
	e.initPoints(edge)
}

// initPoints...
func (e *EdgeKey) initPoints(edge *Edge) {
	direction, _ := edge.direction()
	if direction {
		e.init(edge.pts[0], edge.pts[1])
	} else {
		l := len(edge.pts)
		e.init(edge.pts[l-1], edge.pts[l-2])
	}
}

// init...
func (e *EdgeKey) init(p0, p1 matrix.Matrix) {
	e.p0x = p0[0]
	e.p0y = p0[1]
	e.p1x = p1[0]
	e.p1y = p1[1]
}
