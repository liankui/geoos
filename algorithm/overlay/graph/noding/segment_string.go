package noding

import "github.com/spatial-go/geoos/space"

// An interface for classes which represent a sequence of
// contiguous line segments. SegmentStrings can carry a
// context object, which is useful for preserving topological
// or parentage information.
type SegmentString interface {
	// Gets the user-defined data for this segment string.
	// Returns:
	//		the user-defined data
	getData() interface{}
	// Sets the user-defined data for this segment string.
	// Params:
	//		data – an Object containing user-defined data
	setData(data interface{})
	size() int
	getCoordinate(i int) space.Coordinate
	getCoordinates() []space.Coordinate
	isClosed() bool
}
