package noding

// Base class for Noders which make a single pass to
// find intersections. This allows using a custom
// SegmentIntersector (which for instance may simply
// identify intersections, rather than insert them).
type SinglePassNoder struct {
	Noder
	segInt SegmentIntersector
}

// NewSinglePassNoder...
func NewSinglePassNoder(segInt SegmentIntersector) *SinglePassNoder {
	return &SinglePassNoder{
		segInt: segInt,
	}
}

// computeNodes Computes the noding for a collection of SegmentStrings.
// Some Noders may add all these nodes to the input SegmentStrings;
// others may only add some or none at all.
// Params:
//		segStrings – a collection of SegmentStrings to node
func (s SinglePassNoder) computeNodes(segStrings interface{}) {

}

// getNodedSubstrings Returns a Collection of fully noded SegmentStrings.
// The SegmentStrings have the same context as their parent.
// Returns:
//		a Collection of SegmentStrings
func (s SinglePassNoder) getNodedSubstrings() interface{} {return nil}
