package noding

import "github.com/spatial-go/geoos/algorithm/matrix"

// Computes the possible intersections between two line segments in NodedSegmentStrings
// and adds them to each string using NodedSegmentString.addIntersection(LineIntersector, int, int, int).
type IntersectionAdder struct {
	properIntersectionPoint matrix.Matrix
	li                      *LineIntersector
}

// NewIntersectionAdder...
func NewIntersectionAdder(li *LineIntersector) *IntersectionAdder {
	return &IntersectionAdder{
		li: li,
	}
}

// processIntersections...
func (i *IntersectionAdder) ProcessIntersections(
	e0 matrix.LineMatrix, segIndex0 int,
	e1 matrix.LineMatrix, segIndex1 int) {
}

// isDone...
func (i *IntersectionAdder) IsDone() bool { return false }

func (i *IntersectionAdder) Result() interface{} { return nil }
