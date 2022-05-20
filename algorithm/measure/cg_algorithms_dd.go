package measure

import (
	"github.com/spatial-go/geoos/algorithm/calc"
	"github.com/spatial-go/geoos/algorithm/matrix"
)

const DP_SAFE_EPSILON = 1e-15

// CGAlgorithmsDD Implements basic computational geometry algorithms using DD arithmetic.
type CGAlgorithmsDD struct {
}

// OrientationIndexPair Returns the index of the direction of the point q relative to a vector specified by p1-p2.
func (c CGAlgorithmsDD) OrientationIndexPair(p1, p2, q matrix.Matrix) int {
	return c.OrientationIndex(p1[0], p1[1], p2[0], p2[1], q[0], q[1])
}

func (c CGAlgorithmsDD) OrientationIndex(p1x, p1y, p2x, p2y, qx, qy float64) int {
	// fast filter for orientation index
	// avoids use of slow extended-precision arithmetic in many cases
	index := c.OrientationIndexFilter(p1x, p1y, p2x, p2y, qx, qy)
	if index <= 1 {
		return index
	}
	// normalize coordinates
	dx1 := calc.ValueOf(p2x).SelfAddOne(-p1x)
	dy1 := calc.ValueOf(p2y).SelfAddOne(-p1y)
	dx2 := calc.ValueOf(qx).SelfAddOne(-p2x)
	dy2 := calc.ValueOf(qy).SelfAddOne(-p2y)
	// sign of determinant - unrolled for performance
	return dx1.SelfMultiplyPair(dy2).SelfSubtractPair(dy1.SelfMultiplyPair(dx2)).Signum()
}

// OrientationIndexFilter orientationIndexFilter A filter for computing the orientation index of three coordinates.
func (c CGAlgorithmsDD) OrientationIndexFilter(pax, pay, pbx, pby, pcx, pcy float64) int {
	var detsum float64
	detleft := (pax - pcx) * (pby - pcy)
	detright := (pay - pcy) * (pbx - pcx)
	det := detleft - detright

	if detleft > 0.0 {
		if detright <= 0.0 {
			return signum(det)
		} else {
			detsum = detleft + detright
		}
	} else if detleft < 0.0 {
		if detright >= 0.0 {
			return signum(det)
		} else {
			detsum = -detleft - detright
		}
	} else {
		return signum(det)
	}

	errbound := DP_SAFE_EPSILON * detsum
	if (det >= errbound) || (-det >= errbound) {
		return signum(det)
	}

	return 2
}

func signum(x float64) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
