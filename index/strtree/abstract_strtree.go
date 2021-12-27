package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"sort"
)

const DEFAULT_NODE_CAPACITY = 10

// A query-only R-tree created using the Sort-Tile-Recursive (STR) algorithm. For two-dimensional spatial data.
// The STR packed R-tree is simple to implement and maximizes space utilization;
// that is, as many leaves as possible are filled to capacity.
// Overlap between nodes is far less than in a basic R-tree.
// However, the index is semi-static; once the tree has been built (which happens automatically upon the first query),
// items may not be added.
// Items may be removed from the tree using remove(Envelope, Object).
type AbstractSTRtree struct {
	root           *AbstractNode
	built          bool
	itemBoundables []*AbstractNode
	nodeCapacity   int
}

// build Creates parent nodes, grandparent nodes, and so forth up to the root node,
// for the data that has been inserted into the tree.
// Can only be called once, and thus can be called only after all of the data has been inserted into the tree.
func (s *AbstractSTRtree) build() {
	if s.built {
		return
	}
	if len(s.itemBoundables) == 0 {
		s.root = s.createNode(0)
	} else {
		s.createHigherLevels(s.itemBoundables, -1)
	}
	s.itemBoundables = []*AbstractNode{}
	s.built = true
}

// createNode Create a node.
func (s *AbstractSTRtree) createNode(level int) *AbstractNode {
	return &AbstractNode{
		level: level,
	}
}

// createHigherLevels Creates the levels higher than the given level.
func (s *AbstractSTRtree) createHigherLevels(boundablesOfALevel []*AbstractNode, level int) *AbstractNode {
	parentBoundables := s.createParentBoundables(boundablesOfALevel, level+1)
	if len(parentBoundables) == 1 {
		return parentBoundables[0]
	}
	return s.createHigherLevels(parentBoundables, level+1)
}

// createParentBoundables Sorts the childBoundables then divides them into groups of size M, where M is the node capacity.
func (s *AbstractSTRtree) createParentBoundables(childBoundables []*AbstractNode, newLevel int) []*AbstractNode {
	var parentBoundables []*AbstractNode
	parentBoundables = append(parentBoundables, s.createNode(newLevel))
	sortedChildBoundables := childBoundables
	sort.Slice(sortedChildBoundables, func(i, j int) bool {
		return centreY(*sortedChildBoundables[i].bounds) > centreY(*sortedChildBoundables[j].bounds)
	})
	for _, childBoundable := range sortedChildBoundables {
		if len(parentBoundables[len(parentBoundables)-1].childBoundables) == s.nodeCapacity {
			parentBoundables = append(parentBoundables, s.createNode(newLevel))
		}
		parentBoundables[len(parentBoundables)-1].childBoundables =
			append(parentBoundables[len(parentBoundables)-1].childBoundables, childBoundable)
	}
	return parentBoundables
}


// todo 节点如何用Envelope
func (s *AbstractSTRtree) Insert(itemEnv *envelope.Envelope, item interface{}) error {
	s.itemBoundables = append(s.itemBoundables, item.(*AbstractNode))
	return nil
}

// Query Also builds the tree, if necessary.
func (s *AbstractSTRtree) Query(searchBounds *envelope.Envelope) interface{} {
	s.build()
	matches := make([]*AbstractNode, 0)
	if len(matches) == 0 {
		return matches
	}
	if intersects(s.root.bounds, searchBounds) {

	}

	return nil
}

func (s *AbstractSTRtree) queryInternal(searchBounds *envelope.Envelope, node AbstractNode, matches []*AbstractNode) {
	childBoundables := node.childBoundables
	for i := 0; i < len(childBoundables); i++ {
		childBoundable := childBoundables[i]
		if !intersects(childBoundable.bounds, searchBounds) {
			continue
		}
		// todo
	}
}

// QueryVisitor Queries the index for all items whose extents intersect the given search Envelope,
// and applies an  ItemVisitor to them.
// Note that some kinds of indexes may also return objects which do not in fact
// intersect the query envelope.
func (s *AbstractSTRtree) QueryVisitor(queryEnv *envelope.Envelope, visitor index.ItemVisitor) error {
	return nil
}

// Remove Removes a single item from the tree.
func (s *AbstractSTRtree) Remove(itemEnv *envelope.Envelope, item interface{}) bool {
	return false
}
