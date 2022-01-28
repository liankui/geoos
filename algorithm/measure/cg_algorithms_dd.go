package measure

import "github.com/spatial-go/geoos/algorithm/matrix"

// Implements basic computational geometry algorithms using DD arithmetic.
type CGAlgorithmsDD struct {
}

// Returns the index of the direction of the point q relative to a vector specified by p1-p2.
// Params:
//		p1 – the origin point of the vector
//		p2 – the final point of the vector
//		q – the point to compute the direction to
// Returns:
//		1 if q is counter-clockwise (left) from p1-p2 -1 if q is clockwise (right) from p1-p2 0 if q is collinear with p1-p2
func (c CGAlgorithmsDD) orientationIndex(p1, p2, q matrix.Matrix) int {
	return 0	// todo 未实现
}
