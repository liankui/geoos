package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
)

// A node of an AbstractSTRtree. A node is one of:
// 	1.empty
// 	2.an interior node containing child AbstractNodes
// 	3.a leaf node containing data items (ItemBoundables).
// A node stores the bounds of its children, and its level within the index tree.
type AbstractNode struct {
	ChildBoundables []Boundable       `json:"child_boundables,omitempty"`
	Bounds          *envelope.Envelope `json:"bounds"`
	Level           int                `json:"level"`
}

func (a *AbstractNode) addChildBoundable(childBoundable Boundable) {
	if abstractNode, ok := childBoundable.(*AbstractNode); ok {
		if abstractNode.Bounds.IsNil() {
			return
		}
		a.ChildBoundables = append(a.ChildBoundables, childBoundable)
	}
}

func (a *AbstractNode) getBounds() *envelope.Envelope {
	if a != nil {
		return a.Bounds
	}
	return nil
}