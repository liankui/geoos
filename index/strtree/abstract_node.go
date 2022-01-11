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
	ChildBoundables []Boundable        `json:"child_boundables,omitempty"`
	Bounds          *envelope.Envelope `json:"bounds"`
	Level           int                `json:"level"`
}

// getBounds Gets the bounds of this node
// Returns:
//		the object representing bounds in this index
func (a *AbstractNode) getBounds() *envelope.Envelope {
	if a.Bounds.IsNil() {
		return a.computeBounds()
	}
	return a.Bounds
}

// computeBounds Returns a representation of space that encloses this Boundable, preferably not much bigger than
// this Boundable's boundary yet fast to test for intersection with the bounds of other Boundables.
// The class of object returned depends on the subclass of AbstractSTRtree.
// Returns: an Envelope (for STRtrees), an Interval (for SIRtrees), or other object (for other subclasses of AbstractSTRtree)
func (a *AbstractNode) computeBounds() *envelope.Envelope {
	var bounds envelope.Envelope
	for _, childBoundable := range a.ChildBoundables {
		if bounds.IsNil() {
			bounds = *childBoundable.getBounds()
		} else {
			bounds.ExpandToIncludeEnv(childBoundable.getBounds())
		}
	}
	return &bounds
}

// addChildBoundable Adds either an AbstractNode, or if this is a leaf node, a data object (wrapped in an ItemBoundable)
// Params:
//		childBoundable – the child to add
func (a *AbstractNode) addChildBoundable(childBoundable Boundable) {
	if childBoundable.getBounds().IsNil() {
		return
	}
	a.ChildBoundables = append(a.ChildBoundables, childBoundable)
}