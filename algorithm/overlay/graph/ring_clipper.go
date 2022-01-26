package graph

import "github.com/spatial-go/geoos/algorithm/matrix"

// Clips rings of points to a rectangle. Uses a variant of Cohen-Sutherland clipping.
// In general the output is not topologically valid. In particular, the output may
// contain coincident non-noded line segments along the clip rectangle sides. However,
// the output is sufficiently well-structured that it can be used as input to the
// OverlayNG algorithm (which is able to process coincident linework due to the need to
// handle topology collapse under precision reduction).
// Because of the likelihood of creating extraneous line segments along the clipping
// rectangle sides, this class is not suitable for clipping linestrings.
// The clipping envelope should be generated using RobustClipEnvelopeComputer,
// to ensure that intersecting line segments are not perturbed by clipping.
// This is required to ensure that the overlay of the clipped geometry is robust and correct
// (i.e. the same as if clipping was not used).
type RingClipper struct {
}

// clip Clips a list of points to the clipping rectangle box.
func (r *RingClipper) clip(pts []matrix.Matrix) []matrix.Matrix {
	for edgeIndex := 0; edgeIndex < 4; edgeIndex++ {
		closeRing := edgeIndex == 3
		pts = r.clipToBoxEdge(pts, edgeIndex, closeRing)
		if len(pts) == 0 {
			return pts
		}
	}
	return pts
}

// clipToBoxEdge Clips line to the axis-parallel line defined by a single box edge.
func (r *RingClipper) clipToBoxEdge(pts []matrix.Matrix, edgeIndex int, closeRing bool) []matrix.Matrix {
	ptsClip := matrix.CoordinateList{}
	p0 := pts[len(pts)-1]
	for i := 0; i < len(pts); i++ {
		p1 := pts[i]
		if r.isInsideEdge(p1, edgeIndex) {
			if r.isInsideEdge(p0, edgeIndex) {
				intPt := r.intersection(p0, p1, edgeIndex)
				ptsClip = ptsClip.AddToEndList(intPt, false)
			}
			ptsClip = ptsClip.AddToEndList(p1, false)	// p1.copy
		} else if r.isInsideEdge(p0, edgeIndex) {
			intPt := r.intersection(p0, p1, edgeIndex)
			ptsClip = ptsClip.AddToEndList(intPt, false)
		}
		// else p0-p1 is outside box, so it is dropped
		p0 = p1
	}
	// add closing point if required
	if closeRing && len(ptsClip) > 0 {
		start := ptsClip[0]
		end := ptsClip[len(ptsClip)-1]
		if !start.Equals(end) {
			ptsClip = ptsClip.AddToEndList(start, false)
		}
	}
	return ptsClip
}
