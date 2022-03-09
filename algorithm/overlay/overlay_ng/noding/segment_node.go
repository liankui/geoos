package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
)

// Represents an intersection point between two SegmentStrings.
type SegmentNode struct {
	segString     *NodedSegmentString
	coord         matrix.Matrix
	segmentIndex  int
	segmentOctant int
	isInterior    bool
}

func NewSegmentNode(segString *NodedSegmentString, coord matrix.Matrix, segmentIndex, segmentOctant int) *SegmentNode {
	return &SegmentNode{
		segString:     segString,
		coord:         coord,
		segmentIndex:  segmentIndex,
		segmentOctant: segmentOctant,
		isInterior:    !coord.Equals(segString.GetCoordinate(segmentIndex)),
	}
}
