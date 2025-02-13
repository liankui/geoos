package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"strconv"
)

// NodedSegmentString Represents a list of contiguous line segments, and supports noding
// the segments. The line segments are represented by an array of Coordinates.
// Intended to optimize the noding of contiguous segments by reducing the
// number of allocated objects. SegmentStrings can carry a context object,
// which is useful for preserving topological or parentage information.
// All noded substrings are initialized with the same context object.
type NodedSegmentString struct {
	pts      []matrix.Matrix
	data     interface{}
	nodeList *SegmentNodeList
}

// NewNodedSegmentString Creates a instance from a list of vertices and optional data object.
func NewNodedSegmentString(pts []matrix.Matrix, data interface{}) *NodedSegmentString {
	nodedSegmentString := NodedSegmentString{
		pts:  pts,
		data: data,
	}
	nodedSegmentString.nodeList = NewSegmentNodeList(&nodedSegmentString) // todo 验证写法
	return &nodedSegmentString
}

// GetNodedSubstrings Gets the SegmentStrings which result from
// splitting this string at node points.
func (n *NodedSegmentString) GetNodedSubstrings(segStrings interface{}) []SegmentString {
	resultEdgeList := make([]SegmentString, 0)
	resultEdgeList = n.getNodedSubstrings(segStrings, resultEdgeList).([]SegmentString)
	return resultEdgeList
}

// getNodedSubstrings Adds the noded SegmentStrings which result from splitting this string at node points.
func (n *NodedSegmentString) getNodedSubstrings(segStrings, resultEdgeList interface{}) interface{} {
	for _, ss := range segStrings.([]*NodedSegmentString) {
		fmt.Println("===ss.pts", ss.nodeList.edge.GetCoordinates())
		resultEdgeList = ss.nodeList.addSplitEdges(resultEdgeList)
		fmt.Print("===ss.resultEdgeList")
		PrintlnEdgeList(resultEdgeList)
	}
	fmt.Print("===getNodedSubstrings:")
	PrintlnEdgeList(resultEdgeList)
	return resultEdgeList
}

// todo remove
func PrintlnEdgeList(resultEdgeList interface{}) {
	for i, _ := range resultEdgeList.([]SegmentString) {
		fmt.Print(strconv.Itoa(i)+":", resultEdgeList.([]SegmentString)[i], " ")
	}
	fmt.Println()
}

// getData Gets the user-defined data for this segment string.
func (n *NodedSegmentString) GetData() interface{} { return n.data }

// setData Sets the user-defined data for this segment string.
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

// addIntersectionNode Adds an intersection node for a given point and segment to this segment string.
// If an intersection already exists for this exact location, the existing node will be returned.
func (n *NodedSegmentString) addIntersectionNode(intPt matrix.Matrix, segmentIndex int) *SegmentNode {
	normalizedSegmentIndex := segmentIndex
	//Debug.println("edge intpt: " + intPt + " dist: " + dist);
	// normalize the intersection point location
	nextSegIndex := normalizedSegmentIndex + 1
	if nextSegIndex < len(n.pts) {
		nextPt := n.pts[nextSegIndex]
		// Normalize segment index if intPt falls on vertex
		// The check for point equality is 2D only - Z values are ignored
		if intPt.Equals(nextPt) {
			//Debug.println("normalized distance");
			normalizedSegmentIndex = nextSegIndex
		}
	}
	/*
	  Add the intersection point to edge intersection list.
	*/
	ei := n.nodeList.add(intPt, normalizedSegmentIndex)
	return ei
}
