package measure

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

const (
	// A value that indicates an orientation of clockwise, or a right turn.
	CLOCKWISE = -1
	// A value that indicates an orientation of clockwise, or a right turn.
	RIGHT = CLOCKWISE
	// A value that indicates an orientation of counterclockwise, or a left turn.
	COUNTERCLOCKWISE = 1
	// A value that indicates an orientation of counterclockwise, or a left turn.
	LEFT = COUNTERCLOCKWISE
	// A value that indicates an orientation of collinear, or no turn (straight).
	COLLINEAR = 0
	// A value that indicates an orientation of collinear, or no turn (straight).
	STRAIGHT = COLLINEAR
)

// Functions to compute the orientation of basic geometric structures including
// point triplets (triangles) and rings. Orientation is a fundamental property
// of planar geometries (and more generally geometry on two-dimensional manifolds).
// Determining triangle orientation is notoriously subject to numerical precision
// errors in the case of collinear or nearly collinear points.
type Orientation struct {
}

// Tests if a ring defined by a CoordinateSequence is oriented counter-clockwise.
// The list of points is assumed to have the first and last points equal.
// This handles coordinate lists which contain repeated points.
// This handles rings which contain collapsed segments (in particular, along the top of the ring).
// This algorithm is guaranteed to work with valid rings. It also works with "mildly invalid" rings which contain collapsed (coincident) flat segments along the top of the ring. If the ring is "more" invalid (e.g. self-crosses or touches), the computed result may not be correct.
// Params:
//		ring – a CoordinateSequence forming a ring (with first and last point identical)
// Returns:
//		true if the ring is oriented counter-clockwise.
func (o Orientation) IsCCW(ring matrix.LineMatrix) bool {
	// of points without closing endpoint
	nPts := len(ring) - 1
	// return default value if ring is flat
	if nPts < 3 {
		return false
	}

	/**
	 * Find first highest point after a lower point, if one exists
	 * (e.g. a rising segment)
	 * If one does not exist, hiIndex will remain 0
	 * and the ring must be flat.
	 * Note this relies on the convention that
	 * rings have the same stxart and end point.
	 */
	upHiPt := matrix.Matrix(ring[0])
	prevY := upHiPt[1]
	upLowPt := matrix.Matrix{}
	iUpHi := 0
	for i := 1; i <= nPts; i++ {
		py := ring[i][1]
		/**
		 * If segment is upwards and endpoint is higher, record it
		 */
		if py > prevY && py >= upHiPt[1] {
			upHiPt = ring[i]
			iUpHi = i
			upLowPt = ring[i-1]
		}
		prevY = py
	}
	/**
	 * Check if ring is flat and return default value if so
	 */
	if iUpHi == 0 {
		return false
	}

	/**
	 * Find the next lower point after the high point
	 * (e.g. a falling segment).
	 * This must exist since ring is not flat.
	 */
	iDownLow := iUpHi
	for iDownLow != iUpHi && ring[iDownLow][1] == upHiPt[1] {
		iDownLow = (iDownLow + 1) % nPts
	}

	downLowPt := ring[iDownLow]
	iDownHi := 0
	if iDownLow > 0 {
		iDownHi = iDownLow - 1
	} else {
		iDownHi = nPts - 1
	}
	downHiPt := ring[iDownHi]

	/**
	 * Two cases can occur:
	 * 1) the hiPt and the downPrevPt are the same.
	 *    This is the general position case of a "pointed cap".
	 *    The ring orientation is determined by the orientation of the cap
	 * 2) The hiPt and the downPrevPt are different.
	 *    In this case the top of the cap is flat.
	 *    The ring orientation is given by the direction of the flat segment
	 */
	if upHiPt.Equals(matrix.Matrix(downHiPt)) {
		/**
		 * Check for the case where the cap has configuration A-B-A.
		 * This can happen if the ring does not contain 3 distinct points
		 * (including the case where the input array has fewer than 4 elements), or
		 * it contains coincident line segments.
		 */
		if upLowPt.Equals(upHiPt) || matrix.Matrix(downLowPt).Equals(upHiPt) || upLowPt.Equals(matrix.Matrix(downLowPt)) {
			return false
		}

		/**
		 * It can happen that the top segments are coincident.
		 * This is an invalid ring, which cannot be computed correctly.
		 * In this case the orientation is 0, and the result is false.
		 */
		index := o.index(upLowPt, upHiPt, downLowPt)	// todo 涉及到cg dd算法
		return index == COUNTERCLOCKWISE
	} else {
		/**
		 * Flat cap - direction of flat top determines orientation
		 */
		delX := downHiPt[0] - upHiPt[0]
		return delX < 0
	}
}

// index Returns the orientation index of the direction of the point q relative to
// a directed infinite line specified by p1-p2. The index indicates whether the point
// lies to the LEFT or RIGHT of the line, or lies on it COLLINEAR.
// The index also indicates the orientation of the triangle formed by the three points
// ( COUNTERCLOCKWISE, CLOCKWISE, or STRAIGHT )
// Params:
//		p1 – the origin point of the line vector
//		p2 – the final point of the line vector
//		q – the point to compute the direction to
// Returns:
//		-1 ( CLOCKWISE or RIGHT ) if q is clockwise (right) from p1-p2; 1 ( COUNTERCLOCKWISE or LEFT ) if q is counter-clockwise (left) from p1-p2; 0 ( COLLINEAR or STRAIGHT ) if q is collinear with p1-p2
func (o Orientation) index(p1, p2, q matrix.Matrix) int {
	/*
	 * MD - 9 Aug 2010 It seems that the basic algorithm is slightly orientation
	 * dependent, when computing the orientation of a point very close to a
	 * line. This is possibly due to the arithmetic in the translation to the
	 * origin.
	 *
	 * For instance, the following situation produces identical results in spite
	 * of the inverse orientation of the line segment:
	 *
	 * Coordinate p0 = new Coordinate(219.3649559090992, 140.84159161824724);
	 * Coordinate p1 = new Coordinate(168.9018919682399, -5.713787599646864);
	 *
	 * Coordinate p = new Coordinate(186.80814046338352, 46.28973405831556); int
	 * orient = orientationIndex(p0, p1, p); int orientInv =
	 * orientationIndex(p1, p0, p);
	 *
	 * A way to force consistent results is to normalize the orientation of the
	 * vector using the following code. However, this may make the results of
	 * orientationIndex inconsistent through the triangle of points, so it's not
	 * clear this is an appropriate patch.
	 *
	 */
	var cgDD CGAlgorithmsDD
	return cgDD.orientationIndex(p1, p2, q)

	// testing only
	//return ShewchuksDeterminant.orientationIndex(p1, p2, q);
	// previous implementation - not quite fully robust
	//return RobustDeterminant.orientationIndex(p1, p2, q);
}
