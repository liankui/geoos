package strtree

import "github.com/spatial-go/geoos/algorithm/matrix/envelope"

// A node of an AbstractSTRtree. A node is one of:
// 	- empty
// 	- an interior node containing child AbstractNodes
// 	- a leaf node containing data items (ItemBoundables).
// A node stores the bounds of its children, and its level within the index tree.
type AbstractNode struct {
	childBoundables []*AbstractNode
	bounds          *envelope.Envelope
	level           int
}
