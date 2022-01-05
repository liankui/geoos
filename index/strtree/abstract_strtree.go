package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"sort"
)

const DEFAULT_NODE_CAPACITY = 10

var (
	itemBoundables []Boundable
	built          bool
)

// A query-only R-tree created using the Sort-Tile-Recursive (STR) algorithm. For two-dimensional spatial data.
// The STR packed R-tree is simple to implement and maximizes space utilization;
// that is, as many leaves as possible are filled to capacity.
// Overlap between nodes is far less than in a basic R-tree.
// However, the index is semi-static; once the tree has been built (which happens automatically upon the first query),
// items may not be added.
// Items may be removed from the tree using remove(Envelope, Object).
type AbstractSTRtree struct {
	Root         *AbstractNode `json:"root"`
	NodeCapacity int           `json:"node_capacity"`
}

// build Creates parent nodes, grandparent nodes, and so forth up to the root node,
// for the data that has been inserted into the tree.
// Can only be called once, and thus can be called only after all of the data has been inserted into the tree.
func (s *AbstractSTRtree) build() {
	if built {
		return
	}
	if len(itemBoundables) == 0 {
		s.Root = s.createNode(0)
	} else {
		s.Root = s.createHigherLevels(itemBoundables, -1)
	}
	itemBoundables = nil
	built = true
}

// createNode Create a node.
func (s *AbstractSTRtree) createNode(level int) *AbstractNode {
	abstractNode := &AbstractNode{Level: level}
	abstractNode.Bounds = s.getBounds()
	return abstractNode
}

// getBounds Gets the bounds of this node.
func (s *AbstractSTRtree) getBounds() *envelope.Envelope {
	if s.Root.Bounds.IsNil() {
		return s.computeBounds()
	}
	return s.Root.Bounds
}

// computeBounds Returns a representation of space that encloses this Boundable, preferably not much bigger than
// this Boundable's boundary yet fast to test for intersection with the bounds of other Boundables.
// The class of object returned depends on the subclass of AbstractSTRtree.
// Returns: an Envelope (for STRtrees), an Interval (for SIRtrees), or other object (for other subclasses of AbstractSTRtree)
func (s *AbstractSTRtree) computeBounds() *envelope.Envelope {
	var bounds envelope.Envelope
	for _, childBoundable := range s.Root.ChildBoundables {
		if bounds.IsNil() {
			bounds = *childBoundable.getBounds()
		} else {
			bounds.ExpandToIncludeEnv(childBoundable.getBounds())
		}
	}
	return &bounds
}

// createHigherLevels Creates the levels higher than the given level.
// Params:
// 	boundablesOfALevel – the level to build on
// 	level – the level of the Boundables, or -1 if the boundables are item boundables (that is, below level 0)
// Returns:
// 	the root, which may be a ParentNode or a LeafNode
func (s *AbstractSTRtree) createHigherLevels(boundablesOfALevel []Boundable, level int) *AbstractNode {
	parentBoundables := s.createParentBoundables(boundablesOfALevel, level+1)
	if len(parentBoundables) == 1 {
		return (parentBoundables[0]).(*AbstractNode)
	}

	return s.createHigherLevels(parentBoundables, level+1)
}

// createParentBoundables Sorts the childBoundables then divides them into groups of size M, where M is the node capacity.
func (s *AbstractSTRtree) createParentBoundables(childBoundables []Boundable, newLevel int) []Boundable {
	var parentBoundablesNode []Boundable
	parentBoundablesNode = append(parentBoundablesNode, s.createNode(newLevel))

	sortedChildBoundables := childBoundables
	// Sort from largest to smallest based on the averages of MaxY and MinY.
	sort.Slice(sortedChildBoundables, func(i, j int) bool {
		return centreY(*sortedChildBoundables[i].getBounds()) > centreY(*sortedChildBoundables[j].getBounds())
	})

	lastNode := parentBoundablesNode[len(parentBoundablesNode)-1].(*AbstractNode)
	for _, childBoundable := range sortedChildBoundables {
		if len(lastNode.ChildBoundables) == s.NodeCapacity {
			parentBoundablesNode = append(parentBoundablesNode, s.createNode(newLevel))
		}
		lastNode.addChildBoundable(childBoundable)
	}
	return parentBoundablesNode
}

// getRoot Gets the root node of the tree.
func (s *AbstractSTRtree) getRoot() *AbstractNode {
	s.build()
	return s.Root
}

func (s *AbstractSTRtree) isEmpty() bool {
	if !built {
		return len(itemBoundables) == 0
	}
	return len(s.getRoot().ChildBoundables) == 0
}

// Insert ...
func (s *AbstractSTRtree) Insert(bounds *envelope.Envelope, item interface{}) error {
	if !built {
		return index.ErrSTRtreeInsert
	}
	itemBoundables = append(itemBoundables, &ItemBoundable{Bounds: bounds, Item: item})
	return nil
}

// Query Also builds the tree, if necessary.
func (s *AbstractSTRtree) Query(searchBounds *envelope.Envelope) interface{} {
	s.build()
	matches := make([]interface{}, 0) // todo 结构未知
	if s.isEmpty() {
		return matches
	}
	if intersects(s.Root.getBounds(), searchBounds) {
		s.queryInternal(searchBounds, s.Root, matches)
	}
	return matches
}

// QueryVisitor Also builds the tree, if necessary.
func (s *AbstractSTRtree) QueryVisitor(searchBounds *envelope.Envelope, visitor index.ItemVisitor) error {
	s.build()
	if s.isEmpty() {
		return nil // todo 没有返回值
	}
	if intersects(s.Root.getBounds(), searchBounds) {
		s.queryVisitorInternal(searchBounds, s.Root, visitor)
	}
	return nil
}

// queryInternal ...
func (s *AbstractSTRtree) queryInternal(searchBounds *envelope.Envelope, node *AbstractNode, matches []interface{}) {
	childBoundables := node.ChildBoundables
	for i := 0; i < len(childBoundables); i++ {
		childBoundable := childBoundables[i]
		if !intersects(childBoundable.getBounds(), searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			s.queryInternal(searchBounds, childBoundable.(*AbstractNode), matches)
		case *ItemBoundable:
			matches = append(matches, childBoundable.(*ItemBoundable).getItem())
		}
	}
}

// queryVisitorInternal ...
func (s *AbstractSTRtree) queryVisitorInternal(searchBounds *envelope.Envelope, node *AbstractNode, visitor index.ItemVisitor) {
	childBoundables := node.ChildBoundables
	for i := 0; i < len(childBoundables); i++ {
		childBoundable := childBoundables[i]
		if !intersects(childBoundable.getBounds(), searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			s.queryVisitorInternal(searchBounds, childBoundable.(*AbstractNode), visitor)
		case *ItemBoundable:
			visitor.VisitItem(childBoundable.(*ItemBoundable).getItem())
		}
	}
}

// Remove Removes an item from the tree. (Builds the tree, if necessary.)
func (s *AbstractSTRtree) Remove(searchBounds *envelope.Envelope, item interface{}) bool {
	s.build()
	if intersects(s.Root.getBounds(), searchBounds) {
		return s.remove(searchBounds, s.Root, item)
	}
	return false
}

// removeItem ...
func (s *AbstractSTRtree) removeItem(node *AbstractNode, item interface{}) bool {
	for i := len(node.ChildBoundables); i >= 0; i++ {
		childBoundable := node.ChildBoundables[i]
		switch childBoundable.(type) {
		case *ItemBoundable:
			if childBoundable.(*ItemBoundable).getItem() == item {
				node.ChildBoundables = append(node.ChildBoundables[:i], node.ChildBoundables[i+1:]...)
				return true
			}
		}
	}
	return false
}

// remove ...
func (s *AbstractSTRtree) remove(searchBounds *envelope.Envelope, node *AbstractNode, item interface{}) bool {
	found := s.removeItem(node, item)
	if found {
		return true
	}

	var childToPrune *AbstractNode
	for i := len(node.ChildBoundables); i >= 0; i++ {
		childBoundable := node.ChildBoundables[i]
		if !intersects(childBoundable.getBounds(), searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			found = s.remove(searchBounds, childBoundable.(*AbstractNode), item)
			if found {
				childToPrune = childBoundable.(*AbstractNode)
				if childToPrune != nil {
					if len(childToPrune.ChildBoundables) == 0 {
						node.ChildBoundables = append(node.ChildBoundables[:i], node.ChildBoundables[i+1:]...)
					}
				}
				break
			}
		}
	}
	return found
}
