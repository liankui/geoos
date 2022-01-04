package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"sort"
)

const DEFAULT_NODE_CAPACITY = 10

var (
	itemBoundables []*AbstractNode
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
		s.Root = s.createNode(0) // 创建根节点
	} else {
		s.Root = s.createHigherLevels(itemBoundables, -1) // 创建父节点、祖节点，直至根节点
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
			bounds = *childBoundable.Bounds
		} else {
			bounds.ExpandToIncludeEnv(childBoundable.Bounds)
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
func (s *AbstractSTRtree) createHigherLevels(boundablesOfALevel []*AbstractNode, level int) *AbstractNode {
	parentBoundables := s.createParentBoundables(boundablesOfALevel, level+1)
	if len(parentBoundables) == 1 {
		return parentBoundables[0]
	}
	return s.createHigherLevels(parentBoundables, level+1)
}

// createParentBoundables Sorts the childBoundables then divides them into groups of size M, where M is the node capacity.
func (s *AbstractSTRtree) createParentBoundables(childBoundables []*AbstractNode, newLevel int) []*AbstractNode {
	var parentBoundablesNode []*AbstractNode
	parentBoundablesNode = append(parentBoundablesNode, s.createNode(newLevel))

	sortedChildBoundables := childBoundables
	sort.Slice(sortedChildBoundables, func(i, j int) bool {	// 根据MaxY和MinY的平均值从大到小排序
		return centreY(*sortedChildBoundables[i].Bounds) > centreY(*sortedChildBoundables[j].Bounds)
	})

	for _, childBoundable := range sortedChildBoundables {
		lastNode := parentBoundablesNode[len(parentBoundablesNode)-1]
		if len(lastNode.ChildBoundables) == s.NodeCapacity {
			parentBoundablesNode = append(parentBoundablesNode, s.createNode(newLevel))
		}
		lastNode.addChildBoundable(childBoundable)
	}
	return parentBoundablesNode
}

// Insert ...
func (s *AbstractSTRtree) Insert(itemEnv *envelope.Envelope, item interface{}) error {
	//s.itemBoundables = append(s.itemBoundables, item.(*AbstractNode))
	return nil
}

// Query Also builds the tree, if necessary.
func (s *AbstractSTRtree) Query(searchBounds *envelope.Envelope) interface{} {
	s.build()
	matches := make([]*AbstractNode, 0)
	if len(matches) == 0 {
		return matches
	}
	if intersects(s.Root.Bounds, searchBounds) {

	}

	return nil
}

// queryInternal ...
func (s *AbstractSTRtree) queryInternal(searchBounds *envelope.Envelope, node AbstractNode, matches []*AbstractNode) {
	childBoundables := node.ChildBoundables
	for i := 0; i < len(childBoundables); i++ {
		childBoundable := childBoundables[i]
		if !intersects(childBoundable.Bounds, searchBounds) {
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
