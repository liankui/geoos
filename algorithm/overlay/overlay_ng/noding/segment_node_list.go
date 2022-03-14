package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"strconv"
)

// A list of the SegmentNodes present along a noded SegmentString.
type SegmentNodeList struct {
	nodeSlice []*SegmentNode
	//nodeMap   map[*SegmentNode]*SegmentNode // todo 验证 treeMap
	edge      *NodedSegmentString           // the parent edge
}

// NewSegmentNodeList...
func NewSegmentNodeList(edge *NodedSegmentString) *SegmentNodeList {
	return &SegmentNodeList{
		edge:    edge,
		//nodeMap: make(map[*SegmentNode]*SegmentNode),
		nodeSlice: make([]*SegmentNode, 0),	// todo 如何初始化
	}
}

// addSplitEdges Creates new edges for all the edges that the intersections
// in this list split the parent edge into. Adds the edges to the provided
// argument list (this is so a single list can be used to accumulate all split
// edges for a set of SegmentStrings).
func (s *SegmentNodeList) addSplitEdges(edgeList interface{}) interface{} {
	// ensure that the list has entries for the first and last point of the edge
	s.printSnodeMap()
	fmt.Println("====edge=", s.edge.pts)
	s.addEndpoints()
	fmt.Printf("====s.nodeSlice=%v\n", s.nodeSlice[0])
	s.addCollapsedNodes()

	// there should always be at least two entries in the list, since the endpoints are nodes
	eiPrev := s.nodeSlice[0]	// todo array out of index
	for i := 1; i < len(s.nodeSlice); i++ {
		ei := s.nodeSlice[i]
		fmt.Printf("createSplitEdge,eiPrev=%v,ei=%v\n", eiPrev, ei)
		newEdge := s.createSplitEdge(eiPrev, ei)
		/*
		   if (newEdge.size() < 2)
		     throw new RuntimeException("created single point edge: " + newEdge.toString());
		*/
		edgeList = append(edgeList.([]SegmentString), newEdge)
		eiPrev = ei
	}
	return edgeList
}

// addEndpoints Adds nodes for the first and last points of the edge
func (s *SegmentNodeList) addEndpoints() {
	maxSegIndex := len(s.edge.pts) - 1
	fmt.Println("maxSegIndex=", maxSegIndex)
	fmt.Println("s.edge.GetCoordinate(0)=", s.edge.GetCoordinate(0))
	add := s.add(s.edge.GetCoordinate(0), 0)
	fmt.Printf("--add SegmentNode=%#v\n", add)
	fmt.Printf("s.edge.GetCoordinate(maxSegIndex)%#v\n", s.edge.GetCoordinate(maxSegIndex))
	node := s.add(s.edge.GetCoordinate(maxSegIndex), maxSegIndex)
	fmt.Printf("--node SegmentNode=%#v\n", node)
}

// add Adds an intersection into the list, if it isn't already there.
// The input segmentIndex and dist are expected to be normalized.
func (s *SegmentNodeList) add(intPt matrix.Matrix, segmentIndex int) *SegmentNode {
	fmt.Printf("add intPt=%v, segmentIndex=%v\n", intPt, segmentIndex)
	eiNew := NewSegmentNode(s.edge, intPt, segmentIndex, s.edge.getSegmentOctant(segmentIndex))
	for _, ei := range s.nodeSlice {
		if ei.coord.Equals(eiNew.coord) {
			fmt.Println("ei != null")
			return ei
		}
	}
	//s.nodeMap[eiNew] = eiNew
	fmt.Println("nodeMap.put=", eiNew)
	s.nodeSlice = append(s.nodeSlice, eiNew)
	return eiNew
}

// addCollapsedNodes Adds nodes for any collapsed edge pairs. Collapsed edge pairs can be
// caused by inserted nodes, or they can be pre-existing in the edge vertex list.
// In order to provide the correct fully noded semantics, the vertex at the base of a
// collapsed pair must also be added as a node.
func (s *SegmentNodeList) addCollapsedNodes() {
	collapsedVertexIndexes := make([]int, 0)

	collapsedVertexIndexes = s.findCollapsesFromInsertedNodes(collapsedVertexIndexes)
	collapsedVertexIndexes = s.findCollapsesFromExistingVertices(collapsedVertexIndexes)

	// node the collapses
	for _, vertexIndex := range collapsedVertexIndexes {
		s.add(s.edge.GetCoordinate(vertexIndex), vertexIndex)
	}
	s.printSnodeMap()
}

func (s *SegmentNodeList) printSnodeMap() {
	fmt.Print("===printSnodeMap:")
	for i, _ := range s.nodeSlice {
		fmt.Print(strconv.Itoa(i)+":", s.nodeSlice[i].segString.pts, " ")
	}
	fmt.Println()
}

// findCollapsesFromInsertedNodes Adds nodes for any collapsed edge pairs caused by inserted nodes
// Collapsed edge pairs occur when the same coordinate is inserted as a node both before and
// after an existing edge vertex. To provide the correct fully noded semantics, the vertex must be
// added as a node as well.
func (s *SegmentNodeList) findCollapsesFromInsertedNodes(collapsedVertexIndexes []int) []int {
	collapsedVertexIndex := make([]int, 0)
	// there should always be at least two entries in the list, since the endpoints are nodes
	eiPrev := s.nodeSlice[0]
	var isCollapsed bool
	for i := 1; i < len(s.nodeSlice); i++ {
		ei := s.nodeSlice[i]
		collapsedVertexIndex, isCollapsed = s.findCollapseIndex(eiPrev, ei, collapsedVertexIndex)
		if isCollapsed {
			collapsedVertexIndexes = append(collapsedVertexIndexes, collapsedVertexIndex[0])
		}
		eiPrev = ei
	}
	fmt.Println("collapsedVertexIndexes1=", collapsedVertexIndexes)
	return collapsedVertexIndexes
}

// findCollapseIndex...
func (s *SegmentNodeList) findCollapseIndex(ei0, ei1 *SegmentNode, collapsedVertexIndex []int) ([]int, bool) {
	// only looking for equal nodes
	if !ei0.coord.Equals(ei1.coord) {
		return collapsedVertexIndex, false
	}
	numVerticesBetween := ei1.segmentIndex - ei0.segmentIndex
	if !ei1.isInterior {
		numVerticesBetween--
	}
	// if there is a single vertex between the two equal nodes, this is a collapse
	if numVerticesBetween == 1 {
		collapsedVertexIndex[0] = ei0.segmentIndex + 1
		return collapsedVertexIndex, true
	}
	return collapsedVertexIndex, false
}

// findCollapsesFromExistingVertices Adds nodes for any collapsed edge pairs which are pre-existing in the vertex list.
func (s *SegmentNodeList) findCollapsesFromExistingVertices(collapsedVertexIndexes []int) []int {
	for i := 0; i < len(s.edge.pts)-2; i++ {
		p0 := s.edge.GetCoordinate(i)
		//p1 := s.edge.GetCoordinate(i+1)
		p2 := s.edge.GetCoordinate(i + 2)
		if p0.Equals(p2) {
			// add base of collapse as node
			collapsedVertexIndexes = append(collapsedVertexIndexes, i+1)
		}
	}
	fmt.Println("collapsedVertexIndexes2=", collapsedVertexIndexes)
	return collapsedVertexIndexes
}

// createSplitEdge Create a new "split edge" with the section of points between (and including)
// the two intersections. The label for the new edge is the same as the label for the parent edge.
func (s *SegmentNodeList) createSplitEdge(ei0, ei1 *SegmentNode) SegmentString {
	pts := s.createSplitEdgePts(ei0, ei1)
	return NewNodedSegmentString(pts, s.edge.GetData())
}

// createSplitEdgePts Extracts the points for a split edge running between two nodes.
// The extracted points should contain no duplicate points. There should always be at least
// two points extracted (which will be the given nodes).
func (s *SegmentNodeList) createSplitEdgePts(ei0, ei1 *SegmentNode) []matrix.Matrix {
	npts := ei1.segmentIndex - ei0.segmentIndex + 2
	fmt.Println("createSplitEdgePts npts=", npts)
	// if only two points in split edge they must be the node points
	if npts == 2 {
		return []matrix.Matrix{ei0.coord, ei1.coord}
	}
	lastSegStartPt := s.edge.GetCoordinate(ei1.segmentIndex)
	/**
	 * If the last intersection point is not equal to the its segment start pt,
	 * add it to the points list as well.
	 * This check is needed because the distance metric is not totally reliable!
	 *
	 * Also ensure that the created edge always has at least 2 points.
	 *
	 * The check for point equality is 2D only - Z values are ignored
	 */
	useIntPt1 := ei1.isInterior || !ei1.coord.Equals(lastSegStartPt)
	if !useIntPt1 {
		fmt.Println("(! useIntPt1)")
		npts--
	}
	pts := make([]matrix.Matrix, npts)
	ipt := 0
	pts[ipt] = ei0.coord
	fmt.Println("createSplitEdgePts npts2=", npts)
	for i := ei0.segmentIndex + 1; i <= ei1.segmentIndex; i++ {
		ipt++
		pts[ipt] = s.edge.GetCoordinate(i)
		fmt.Println("ipt2=", ipt)
	}
	if useIntPt1 {
		pts[ipt] = ei1.coord
	}
	return pts
}
