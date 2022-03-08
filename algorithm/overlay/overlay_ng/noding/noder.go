package noding

// Computes all intersections between segments in a set of SegmentStrings.
// Intersections found are represented as SegmentNodes and added to the
// SegmentStrings in which they occur. As a final step in the noding
// a new set of segment strings split at the nodes may be returned.
type Noder interface {
	// Computes the noding for a collection of SegmentStrings.
	// Some Noders may add all these nodes to the input SegmentStrings;
	// others may only add some or none at all.
	// Params:
	//		segStrings – a collection of SegmentStrings to node
	ComputeNodes(segStrings interface{})

	// Returns a Collection of fully noded SegmentStrings.
	// The SegmentStrings have the same context as their parent.
	// Returns:
	//		a Collection of SegmentStrings
	GetNodedSubstrings() interface{}
}
