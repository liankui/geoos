package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// An interface for classes which represent a sequence of
// contiguous line segments. SegmentStrings can carry a
// context object, which is useful for preserving topological
// or parentage information.
type SegmentString interface {
	// Gets the user-defined data for this segment string.
	// Returns:
	//		the user-defined data
	GetData() interface{}
	// Sets the user-defined data for this segment string.
	// Params:
	//		data – an Object containing user-defined data
	SetData(data interface{})
	Size() int
	GetCoordinate(i int) matrix.Matrix
	GetCoordinates() []matrix.Matrix
	IsClosed() bool
}
