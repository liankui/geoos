package strtree

import "github.com/spatial-go/geoos/algorithm/matrix/envelope"

// Boundable A spatial object in an AbstractSTRtree.
type Boundable interface {
	// Returns a representation of space that encloses this Boundable,
	// preferably not much bigger than this Boundable's boundary yet
	// fast to test for intersection with the bounds of other Boundables.
	// The class of object returned depends on the subclass of AbstractSTRtree.
	// Returns:
	// 	an Envelope (for STRtrees), an Interval (for SIRtrees),
	//	or other object (for other subclasses of AbstractSTRtree)
	getBounds() *envelope.Envelope
}
