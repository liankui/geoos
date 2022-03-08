package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// Represents a list of contiguous line segments, and supports noding
// the segments. The line segments are represented by an array of Coordinates.
// Intended to optimize the noding of contiguous segments by reducing the
// number of allocated objects. SegmentStrings can carry a context object,
// which is useful for preserving topological or parentage information.
// All noded substrings are initialized with the same context object.
type NodedSegmentString struct {
	pts      []matrix.Matrix
	data     interface{}
	nodeList SegmentNodeList
}

// NewNodedSegmentString Creates a instance from a list of vertices and optional data object.
func NewNodedSegmentString(pts []matrix.Matrix, data interface{}) *NodedSegmentString {
	return &NodedSegmentString{
		pts:  pts,
		data: data,
	}
}

// GetNodedSubstrings Gets the SegmentStrings which result from
// splitting this string at node points.
func (n *NodedSegmentString) GetNodedSubstrings(segStrings interface{}) []SegmentString {
	resultEdgelist := make([]SegmentString, 0)
	n.getNodedSubstrings(segStrings, resultEdgelist)
	return resultEdgelist
}

// getNodedSubstrings Adds the noded SegmentStrings which result from splitting this string at node points.
func (n *NodedSegmentString) getNodedSubstrings(segStrings, resultEdgelist interface{}) {
	for _, ss := range segStrings.([]*NodedSegmentString) {
		ss.nodeList.addSplitEdges(resultEdgelist)
	}
}

// getData Gets the user-defined data for this segment string.
// Returns:
//		the user-defined data
func (n *NodedSegmentString) GetData() interface{} { return n.data }

// setData Sets the user-defined data for this segment string.
// Params:
//		data – an Object containing user-defined data
func (n *NodedSegmentString) SetData(data interface{})          { n.data = data }
func (n *NodedSegmentString) Size() int                         { return len(n.pts) }
func (n *NodedSegmentString) GetCoordinate(i int) matrix.Matrix { return n.pts[i] }
func (n *NodedSegmentString) GetCoordinates() []matrix.Matrix   { return n.pts }
func (n *NodedSegmentString) IsClosed() bool {
	return n.pts[0].Equals(n.pts[len(n.pts)-1])
}

// getSegmentOctant Gets the octant of the segment starting at vertex index.
func (n *NodedSegmentString) getSegmentOctant(index int) int {
	if index == len(n.pts)-1 {
		return -1
	}
	return n.safeOctant(n.GetCoordinate(index), n.GetCoordinate(index+1))
}

// safeOctant...
func (n *NodedSegmentString) safeOctant(p0, p1 matrix.Matrix) int {
	if p0.Equals(p1) {
		return 0
	}
	var oct Octant
	return oct.Octant(p0, p1)
}
