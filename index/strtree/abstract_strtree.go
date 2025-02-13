package strtree

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"sort"
)

const DEFAULT_NODE_CAPACITY = 10

// A query-only R-tree created using the Sort-Tile-Recursive (STR) algorithm.
// For two-dimensional spatial data. The STR packed R-tree is simple to implement
// and maximizes space utilization; that is, as many leaves as possible are filled
// to capacity. Overlap between nodes is far less than in a basic R-tree.
// However, the index is semi-static; once the tree has been built (which happens
// automatically upon the first query), items may not be added.
// Items may be removed from the tree using remove(Envelope, Object).
type AbstractSTRtree struct {
	Root           *AbstractNode
	NodeCapacity   int // default=DEFAULT_NODE_CAPACITY
	itemBoundables []Boundable
	built          bool
}

// getItemBoundables...
func (s *AbstractSTRtree) GetItemBoundables() []Boundable {
	return s.itemBoundables
}

// build Creates parent nodes, grandparent nodes, and so forth up to the root node,
// for the data that has been inserted into the tree. Can only be called once,
// and thus can be called only after all of the data has been inserted into the tree.
func (s *AbstractSTRtree) build() {
	if s.built {
		return
	}
	if len(s.itemBoundables) == 0 {
		s.Root = s.createNode(0)
	} else {
		s.Root = s.createHigherLevels(s.itemBoundables, -1)
	}
	s.itemBoundables = nil
	s.built = true
}

// createNode Create a node.
func (s *AbstractSTRtree) createNode(level int) *AbstractNode {
	return &AbstractNode{Level: level}
}

// getBounds Gets the bounds of this node.
func (s *AbstractSTRtree) getBounds() *envelope.Envelope {
	if s.Root.Bounds.IsNil() {
		return s.Root.computeBounds()
	}
	return s.Root.Bounds
}

// createHigherLevels Creates the levels higher than the given level.
// Params:
// 		boundablesOfALevel – the level to build on.
// 		level – the level of the Boundables, or -1 if the boundables are item boundables (that is, below level 0).
// Returns:
// 		the root, which may be a ParentNode or a LeafNode.
func (s *AbstractSTRtree) createHigherLevels(boundablesOfALevel []Boundable, level int) *AbstractNode {
	fmt.Printf("==createHigherLevels level=%v\n", level)
	if len(boundablesOfALevel) == 0 {
		return nil
	}
	parentBoundables := s.createParentBoundables(boundablesOfALevel, level+1)
	fmt.Printf("==len=%v\n", len(parentBoundables))
	if len(parentBoundables) == 1 {
		return (parentBoundables[0]).(*AbstractNode)
	}
	return s.createHigherLevels(parentBoundables, level+1)
}

// createParentBoundables Sorts the childBoundables then divides them into groups of getSize M, where M is the node capacity.
func (s *AbstractSTRtree) createParentBoundables(childBoundables []Boundable, newLevel int) []Boundable {
	if len(childBoundables) == 0 {
		return nil
	}
	parentBoundables := make([]Boundable, 0)
	parentBoundables = append(parentBoundables, s.createNode(newLevel))

	sortedChildBoundables := make([]Boundable, len(childBoundables))
	copy(sortedChildBoundables, childBoundables)

	// Sort from largest to smallest based on the averages of MaxY and MinY.
	sort.Slice(sortedChildBoundables, func(i, j int) bool {
		return centreY(childBoundables[i].getBounds()) > centreY(childBoundables[j].getBounds())
	})
	fmt.Println("===len(sortedChildBoundables)=", len(sortedChildBoundables))	// todo 数量不对

	for _, childBoundable := range sortedChildBoundables {
		if len(parentBoundables[len(parentBoundables)-1].(*AbstractNode).ChildBoundables) == s.NodeCapacity {
			parentBoundables = append(parentBoundables, s.createNode(newLevel))
		}
		parentBoundables[len(parentBoundables)-1].(*AbstractNode).addChildBoundable(childBoundable)
	}
	return parentBoundables
}

// getRoot Gets the root node of the tree.
func (s *AbstractSTRtree) getRoot() *AbstractNode {
	s.build()
	return s.Root
}

// isEmpty ...
func (s *AbstractSTRtree) isEmpty() bool {
	if !s.built {
		return len(s.itemBoundables) == 0
	}
	return len(s.Root.ChildBoundables) == 0
}

// Insert ...
func (s *AbstractSTRtree) insert(bounds *envelope.Envelope, item interface{}) error {
	if s.built {
		return index.ErrSTRtreeInsert
	}
	s.itemBoundables = append(s.itemBoundables, &ItemBoundable{Bounds: bounds, Item: item})
	return nil
}

// Query Also builds the tree, if necessary.
func (s *AbstractSTRtree) query(searchBounds *envelope.Envelope) interface{} {
	s.build()
	matches := make([]interface{}, 0)
	if s.isEmpty() {
		fmt.Println("s.isEmpty")
		return matches
	}
	if s.Root.getBounds().IsIntersects(searchBounds) {
		matches, _ = s.queryInternal(searchBounds, s.Root, matches)
		//fmt.Printf("----intersects, matches2=%#v\n", matches)
	}
	return matches
}

// QueryVisitor Also builds the tree, if necessary.
func (s *AbstractSTRtree) queryVisitor(searchBounds *envelope.Envelope, visitor index.ItemVisitor) error {
	s.build()
	if s.isEmpty() {
		return index.ErrSTRtreeIsEmpty
	}
	if s.Root.getBounds().IsIntersects(searchBounds) {
		return s.queryVisitorInternal(searchBounds, s.Root, visitor)
	}
	return nil
}

// queryInternal ...
func (s *AbstractSTRtree) queryInternal(searchBounds *envelope.Envelope, node *AbstractNode, matches []interface{}) ([]interface{}, error) {
	fmt.Println("------queryInternal")
	childBoundables := node.ChildBoundables
	for _, childBoundable := range childBoundables {
		if !childBoundable.getBounds().IsIntersects(searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			matches, _ = s.queryInternal(searchBounds, childBoundable.(*AbstractNode), matches)
		case *ItemBoundable:
			matches = append(matches, childBoundable.(*ItemBoundable).getItem())
		default:
			return nil, index.ErrSTRtreeNeverReach
		}
	}
	return matches, nil
}

// queryVisitorInternal ...
func (s *AbstractSTRtree) queryVisitorInternal(searchBounds *envelope.Envelope, node *AbstractNode, visitor index.ItemVisitor) error {
	childBoundables := node.ChildBoundables
	for _, childBoundable := range childBoundables {
		if !childBoundable.getBounds().IsIntersects(searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			return s.queryVisitorInternal(searchBounds, childBoundable.(*AbstractNode), visitor)
		case *ItemBoundable:
			visitor.VisitItem(childBoundable.(*ItemBoundable).getItem())
		default:
			return index.ErrSTRtreeNeverReach
		}
	}
	return nil
}

// remove Removes an item from the tree. (Builds the tree, if necessary.)
func (s *AbstractSTRtree) remove(searchBounds *envelope.Envelope, item interface{}) bool {
	s.build()
	if s.Root.getBounds().IsIntersects(searchBounds) {
		return s.removeNode(searchBounds, s.Root, item)
	}
	return false
}

// removeItem ...
func (s *AbstractSTRtree) removeItem(node *AbstractNode, item interface{}) bool {
	for i := len(node.ChildBoundables) - 1; i >= 0; i-- {
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

// removeNode ...
func (s *AbstractSTRtree) removeNode(searchBounds *envelope.Envelope, node *AbstractNode, item interface{}) bool {
	// first try removing item from this node
	found := s.removeItem(node, item)
	if found {
		return true
	}
	// next try removing item from lower nodes
	var childToPrune *AbstractNode
	for i := len(node.ChildBoundables) - 1; i >= 0; i-- {
		childBoundable := node.ChildBoundables[i]
		if !childBoundable.getBounds().IsIntersects(searchBounds) {
			continue
		}
		switch childBoundable.(type) {
		case *AbstractNode:
			found = s.removeNode(searchBounds, childBoundable.(*AbstractNode), item)
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

// size Returns the number of items in the tree.
func (s *AbstractSTRtree) size() int {
	if s.isEmpty() {
		return 0
	}
	s.build()
	return s.getSize(s.Root)
}

// getSize ...
func (s *AbstractSTRtree) getSize(node *AbstractNode) (size int) {
	for _, childBoundable := range node.ChildBoundables {
		switch childBoundable.(type) {
		case *AbstractNode:
			size += s.getSize(childBoundable.(*AbstractNode))
		case *ItemBoundable:
			size++
		}
	}
	return
}
