package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// Finds intersections between line segments which will be snap-rounded, and adds them as nodes to the segments.
// Intersections are detected and computed using full precision. Snapping takes place in a subsequent phase.
// The intersection points are recorded, so that HotPixels can be created for them.
// To avoid robustness issues with vertices which lie very close to line segments a heuristic is used:
// nodes are created if a vertex lies within a tolerance distance of the interior of a segment. The tolerance
// distance is chosen to be significantly below the snap-rounding grid size. This has empirically proven to
// eliminate noding failures.
type SnapRoundingIntersectionAdder struct {
	li            *LineIntersector
	intersections []matrix.Matrix
	nearnessTol   float64
}

// NewSnapRoundingIntersectionAdder...
func NewSnapRoundingIntersectionAdder(nearnessTol float64) *SnapRoundingIntersectionAdder {
	/**
	 * Intersections are detected and computed using full precision.
	 * They are snapped in a subsequent phase.
	 */
	return &SnapRoundingIntersectionAdder{
		li:            new(LineIntersector), // todo
		intersections: make([]matrix.Matrix, 0),
		nearnessTol:   nearnessTol,
	}
}

// processIntersections...
func (s *SnapRoundingIntersectionAdder) processIntersections(
	e0, e1 SegmentString, segIndex0, segIndex1 int) {

}

// isDone...
func (s *SnapRoundingIntersectionAdder) isDone() bool {
	return false
}
