package noding

import "github.com/spatial-go/geoos/space"

// Represents a list of contiguous line segments, and supports noding
// the segments. The line segments are represented by an array of Coordinates.
// Intended to optimize the noding of contiguous segments by reducing the
// number of allocated objects. SegmentStrings can carry a context object,
// which is useful for preserving topological or parentage information.
// All noded substrings are initialized with the same context object.
type NodedSegmentString struct {
	pts      []space.Coordinate
	data     interface{}
	nodeList SegmentNodeList
}

// getNodedSubstrings Gets the SegmentStrings which result from
// splitting this string at node points.
// Params:
//		segStrings – a Collection of NodedSegmentStrings
// Returns:
//		a Collection of NodedSegmentStrings representing the substrings
func (n *NodedSegmentString) getNodedSubstrings(segStrings interface{}) {

}

func (n *NodedSegmentString) getNodedSubstrings2(segStrings, resultEdgelist interface{}) {

}

// getData Gets the user-defined data for this segment string.
// Returns:
//		the user-defined data
func (n *NodedSegmentString) getData() interface{} { return n.data }

// setData Sets the user-defined data for this segment string.
// Params:
//		data – an Object containing user-defined data
func (n *NodedSegmentString) setData(data interface{})             { n.data = data }
func (n *NodedSegmentString) size() int                            { return len(n.pts) }
func (n *NodedSegmentString) getCoordinate(i int) space.Coordinate { return n.pts[i] }
func (n *NodedSegmentString) getCoordinates() []space.Coordinate   { return n.pts }
func (n *NodedSegmentString) isClosed() bool {
	return n.pts[0] == n.pts[len(n.pts)-1]
}
