package noding

import "github.com/spatial-go/geoos/algorithm/matrix"

// IntersectionAdder Computes the possible intersections between two line segments in NodedSegmentStrings
// and adds them to each string using NodedSegmentString.addIntersection(LineIntersector, int, int, int).
type IntersectionAdder struct {
	properIntersectionPoint matrix.Matrix
	li                      *LineIntersector
}

// NewIntersectionAdder ...
func NewIntersectionAdder(li *LineIntersector) *IntersectionAdder {
	return &IntersectionAdder{
		li: li,
	}
}

// processIntersections...
func (i *IntersectionAdder) processIntersections(e0, e1 SegmentString, segIndex0, segIndex1 int) {

}

// isDone...
func (i *IntersectionAdder) isDone() bool {
	return false
}
