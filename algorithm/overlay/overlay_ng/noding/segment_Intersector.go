package noding

// Processes possible intersections detected by a Noder.
// The SegmentIntersector is passed to a Noder. The
// processIntersections(SegmentString, int, SegmentString, int)
// method is called whenever the Noder detects that two
// SegmentStrings might intersect. This class may be used
// either to find all intersections, or to detect the presence
// of an intersection. In the latter case, Noders may choose to
// short-circuit their computation by calling the isDone() method.
// This class is an example of the Strategy pattern.
type SegmentIntersector interface {
	// This method is called by clients of the SegmentIntersector
	// interface to process intersections for two segments of the
	// SegmentStrings being intersected.
	processIntersections(e0, e1 SegmentString, segIndex0, segIndex1 int)

	// Reports whether the client of this class needs to continue
	// testing all intersections in an arrangement.
	// Returns: true if there is no need to continue testing segments
	isDone() bool
}
