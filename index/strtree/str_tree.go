package strtree

import (
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"math"
	"sort"
)

type STRtree struct {
	*AbstractSTRtree
}

func centreX(e envelope.Envelope) float64 {
	return avg(e.MinX, e.MaxX)
}

func centreY(e envelope.Envelope) float64 {
	return avg(e.MinY, e.MaxY)
}

func avg(a, b float64) float64 {
	return (a + b) / 2.0
}

func intersects(a, b *envelope.Envelope) bool {
	return a.IsIntersects(b)
}

// CreateParentBoundables Creates the parent level for the given child level.
// First, orders the items by the x-values of the midpoints, and groups them into vertical slices.
// For each slice, orders the items by the y-values of the midpoints,
// and group them into runs of Size M (the node capacity).
// For each run, creates a new (parent) node.
func (s *STRtree) CreateParentBoundables(childBoundables []Boundable, newLevel int) []Boundable {
	if len(childBoundables) == 0 {
		return nil
	}
	minLeafCount := int(math.Ceil(float64(len(childBoundables))) / float64(s.NodeCapacity))
	sortedChildBoundables := childBoundables
	// Sort from largest to smallest based on the averages of MaxX and MinX.
	sort.Slice(sortedChildBoundables, func(i, j int) bool {
		return centreX(*sortedChildBoundables[i].getBounds()) > centreX(*sortedChildBoundables[j].getBounds())
	})
	verticalSlices := s.verticalSlices(sortedChildBoundables, int(math.Ceil(math.Sqrt(float64(minLeafCount)))))
	return s.createParentBoundablesFromVerticalSlices(verticalSlices, newLevel)
}

// verticalSlices...
func (s *STRtree) verticalSlices(childBoundables []Boundable, sliceCount int) [][]Boundable {
	sliceCapacity := int(math.Ceil(float64(len(childBoundables)) / float64(sliceCount)))
	slices := make([][]Boundable, sliceCount)
	for i, j := 0, 0; j < sliceCount; j++ {
		slices[j] = []Boundable{}
		boundablesAddedToSlice := 0
		for i < len(childBoundables) && boundablesAddedToSlice < sliceCapacity {
			slices[j] = append(slices[j], childBoundables[i])
			boundablesAddedToSlice++
			i++
		}
	}
	return slices
}

// createParentBoundablesFromVerticalSlices...
func (s *STRtree) createParentBoundablesFromVerticalSlices(verticalSlices [][]Boundable, newLevel int) []Boundable {
	if len(verticalSlices) == 0 {
		return nil
	}
	var parentBoundables []Boundable
	for i := 0; i < len(verticalSlices); i++ {
		parentBoundables = append(parentBoundables,
			s.createParentBoundablesFromVerticalSlice(verticalSlices[i], newLevel)...)
	}
	return parentBoundables
}

// createParentBoundablesFromVerticalSlice...
func (s *STRtree) createParentBoundablesFromVerticalSlice(childBoundables []Boundable, newLevel int) []Boundable {
	return s.createParentBoundables(childBoundables, newLevel)
}

func (s *STRtree) CreateNode(level int) *AbstractNode {
	return s.createNode(level)
}

// Insert Inserts an item having the given bounds into the tree.
func (s *STRtree) Insert(bounds *envelope.Envelope, item interface{}) error {
	if bounds.IsNil() {
		return index.ErrSTRtreeBoundsIsNil
	}
	return s.insert(bounds, item)
}

// Query Returns items whose bounds intersect the given envelope.
func (s *STRtree) Query(searchBounds *envelope.Envelope) interface{} {
	return s.query(searchBounds)
}

// QueryVisitor Returns items whose bounds intersect the given envelope.
func (s *STRtree) QueryVisitor(searchBounds *envelope.Envelope, visitor index.ItemVisitor) error {
	return s.queryVisitor(searchBounds, visitor)
}

// Remove Removes a single item from the tree.
// Params:
//		itemEnv – the Envelope of the item to remove
//		item – the item to remove
// Returns:
//		true if the item was found
func (s *STRtree) Remove(searchBounds *envelope.Envelope, item interface{}) bool {
	return s.remove(searchBounds, item)
}

// Size Returns the number of items in the tree.
func (s *STRtree) Size() int {
	return s.size()
}
