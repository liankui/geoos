package noding

// Nodes a set of SegmentStrings using a index based on
// MonotoneChains and a SpatialIndex. The SpatialIndex
// used should be something that supports envelope (range)
// queries efficiently (such as a Quadtree} or STRtree
// (which is the default index provided).
// The noder supports using an overlap tolerance distance.
// This allows determining segment intersection using a buffer
// for uses involving snapping with a distance tolerance.
type MCIndexNoder struct {
	overlapTolerance float64
	nodedSegStrings  interface{}
}

func (m *MCIndexNoder) computeNodes(inputSegStrings interface{}) {
	m.nodedSegStrings = inputSegStrings

}

func (m *MCIndexNoder) getNodedSubstrings() interface{} {

	//NodedSegmentString.getNodedSubstrings(nodedSegStrings)
	return nil
}
