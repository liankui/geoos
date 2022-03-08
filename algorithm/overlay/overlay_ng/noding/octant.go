package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"log"
	"math"
)

// Methods for computing and working with octants of the Cartesian plane Octants are numbered as follows:
//    \2|1/
//   3 \|/ 0
//   ---+--
//   4 /|\ 7
//    /5|6\
//
// If line segments lie along a coordinate axis, the octant is the lower of the two possible values.
type Octant struct {
}

// Octant Returns the octant of a directed line segment from p0 to p1.
func (o *Octant) Octant(p0, p1 matrix.Matrix) int {
	dx := p1[0] - p0[0]
	dy := p1[1] - p0[1]
	if dx == 0.0 && dy == 0.0 {
		log.Printf("Cannot compute the octant for two identical points p0=%v\n", p0)
		return 0
	}
	return o.octant(dx, dy)
}

// octant Returns the octant of a directed line segment (specified as x and y displacements,
// which cannot both be 0).
func (o *Octant) octant(dx, dy float64) int {
	adx := math.Abs(dx)
	ady := math.Abs(dy)

	if dx >= 0 {
		if dy >= 0 {
			if adx >= ady {
				return 0
			} else {
				return 1
			}
		} else { // dy < 0
			if adx >= ady {
				return 7
			} else {
				return 6
			}
		}
	} else { // dx < 0
		if dy >= 0 {
			if adx >= ady {
				return 3
			} else {
				return 2
			}
		} else { // dy < 0
			if adx >= ady {
				return 4
			} else {
				return 5
			}
		}
	}
}
