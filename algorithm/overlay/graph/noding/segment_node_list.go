package noding

type SegmentNodeList struct {
	nodeMap map[string]interface{}
	edge    *NodedSegmentString // the parent edge
}

func NewSegmentNodeList(edge *NodedSegmentString) *SegmentNodeList {
	return &SegmentNodeList{
		edge: edge,
	}
}
