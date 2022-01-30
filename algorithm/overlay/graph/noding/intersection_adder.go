package noding

import "github.com/spatial-go/geoos/algorithm/matrix"

// Computes the possible intersections between two line segments in NodedSegmentStrings
// and adds them to each string using NodedSegmentString.addIntersection(LineIntersector, int, int, int).
type IntersectionAdder struct {
	properIntersectionPoint matrix.Matrix
	li                      *LineIntersector
}

func NewIntersectionAdder(li *LineIntersector) *IntersectionAdder {
	return &IntersectionAdder{
		li: li,
	}
}

func (i *IntersectionAdder) processIntersections(e0, e1 SegmentString, segIndex0, segIndex1 int) {

}

func (i *IntersectionAdder) isDone() {

}
