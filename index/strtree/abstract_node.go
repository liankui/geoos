package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

// A node of an AbstractSTRtree. A node is one of:
// 	- empty
// 	- an interior node containing child AbstractNodes
// 	- a leaf node containing data items (ItemBoundables).
// A node stores the bounds of its children, and its level within the index tree.
type AbstractNode struct {
	ChildBoundables []*AbstractNode    `json:"child_boundables,omitempty"` // 如果是叶子节点，结构为ItemBoundables
	Bounds          *envelope.Envelope `json:"bounds"`
	Level           int                `json:"level"`
}

func (a *AbstractNode) addChildBoundable(childBoundable *AbstractNode) {
	if childBoundable.Bounds.IsNil() {
		return
	}
	a.ChildBoundables = append(a.ChildBoundables, childBoundable)
}
