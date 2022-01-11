package strtree

import (
	"fmt"
	"github.com/spatial-go/geoos/index"
	"testing"

	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/stretchr/testify/assert"
)

func NewDefaultSTRtree() STRtree {
	return STRtree{&AbstractSTRtree{
		Root:         new(AbstractNode),
		NodeCapacity: DEFAULT_NODE_CAPACITY,
	}}
}

func TestEmptyTreeUsingListQuery(t *testing.T) {
	tree := NewDefaultSTRtree()
	list := tree.Query(&envelope.Envelope{MaxX: 0, MinX: 0, MaxY: 1, MinY: 1})
	assert.Equal(t, 0, len(list.([]interface{})), "list's length is equal to 0")
}

func TestEmptyTreeUsingItemVisitorQuery(t *testing.T) {
	tree := NewDefaultSTRtree()
	_ = tree.QueryVisitor(&envelope.Envelope{MaxX: 0, MinX: 0, MaxY: 1, MinY: 1}, *new(index.ItemVisitor))
	assert.True(t, true, "Should never reach here")
}

func TestCreateParentsFromVerticalSlice(t *testing.T) {
	doTestCreateParentsFromVerticalSlice(t, 3,2,2,1)
	doTestCreateParentsFromVerticalSlice(t, 4,2,2,1)
	doTestCreateParentsFromVerticalSlice(t, 5,2,2,1)
}

func doTestCreateParentsFromVerticalSlice(t *testing.T, childCount, nodeCapacity,
	expectedChildrenPerParentBoundable, expectedChildrenOfLastParent int) {
	tree := STRtree{&AbstractSTRtree{
		Root:         new(AbstractNode),
		NodeCapacity: nodeCapacity,
	}}

	parentBoundables := tree.createParentBoundablesFromVerticalSlice(itemWrappers(childCount), 0)
	for i := 0; i < len(parentBoundables)-1; i++ {
		parentBoundable := parentBoundables[i].(*AbstractNode)
		assert.Equal(t, expectedChildrenPerParentBoundable, len(parentBoundable.ChildBoundables))
	}
	lastParent := parentBoundables[len(parentBoundables)-1].(*AbstractNode)
	assert.Equal(t, expectedChildrenOfLastParent, len(lastParent.ChildBoundables))
}

func itemWrappers(size int) (itemWrappers []Boundable) {
	for i := 0; i < size; i++ {
		itemWrappers = append(itemWrappers, &ItemBoundable{&envelope.Envelope{0, 0, 0, 0}, nil})
	}
	return
}

func TestBuildTree(t *testing.T) {
	aNode := &AbstractNode{
		ChildBoundables: nil,
		Bounds: &envelope.Envelope{
			MaxX: 10,
			MinX: 0,
			MaxY: 10,
			MinY: 0,
		},
		Level: 0,
	}
	abSTRtree := &AbstractSTRtree{
		Root:         aNode,
		NodeCapacity: 0,
	}
	abSTRtree.build()

	fmt.Printf("abSTRtree.root:%+v\n", *abSTRtree.Root)
	fmt.Printf("abSTRtree.root.Bounds:%+v\n", abSTRtree.Root.Bounds)
}
