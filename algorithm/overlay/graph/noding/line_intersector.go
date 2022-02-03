package noding

// A LineIntersector is an algorithm that can both test whether two line segments
// intersect and compute the intersection point(s) if they do.
// There are three possible outcomes when determining whether two line segments intersect:
//		NO_INTERSECTION - the segments do not intersect
//		POINT_INTERSECTION - the segments intersect in a single point
//		COLLINEAR_INTERSECTION - the segments are collinear and they intersect in a line segment
// For segments which intersect in a single point, the point may be either an endpoint or
// in the interior of each segment. If the point lies in the interior of both segments,
// this is termed a proper intersection. The method isProper() test for this situation.
// The intersection point(s) may be computed in a precise or non-precise manner.
// Computing an intersection point precisely involves rounding it via a supplied Pm.
// LineIntersectors do not perform an initial envelope intersection test to determine
// if the segments are disjoint. This is because this class is likely to be used in a context
// where envelope overlap is already known to occur (or be likely).
type LineIntersector struct {

}

func (l *LineIntersector) processIntersections(e0, e1 SegmentString, segIndex0, segIndex1 int) {

}

func (l *LineIntersector) isDone() {

}

