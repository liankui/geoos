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

func NewSTRtree(nodeCapacity int) STRtree {
	return STRtree{&AbstractSTRtree{
		Root:         new(AbstractNode),
		NodeCapacity: nodeCapacity,
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
	doTestCreateParentsFromVerticalSlice(t, 3, 2, 2, 1)
	doTestCreateParentsFromVerticalSlice(t, 4, 2, 2, 2)
	doTestCreateParentsFromVerticalSlice(t, 5, 2, 2, 1)
}

func TestDisallowedInserts(t *testing.T) {
	tree := NewSTRtree(5)
	_ = tree.Insert(&envelope.Envelope{MaxX: 0, MinX: 0, MaxY: 0, MinY: 0}, struct{}{})
	_ = tree.Insert(&envelope.Envelope{MaxX: 0, MinX: 0, MaxY: 0, MinY: 0}, struct{}{})
	_ = tree.Query(&envelope.Envelope{})
	fmt.Printf("----tree:%#v\n", tree.AbstractSTRtree.Root)
	err := tree.Insert(&envelope.Envelope{MaxX: 0, MinX: 0, MaxY: 0, MinY: 0}, struct{}{})
	if err != nil {
		fmt.Println("err:", err)
		assert.True(t, true)
	} else {
		assert.True(t, false)
	}
}

func TestQuery(t *testing.T) {
	// todo
}

func TestVerticalSlices(t *testing.T) {
	doTestVerticalSlices(t, 3, 2, 2, 1)
	doTestVerticalSlices(t, 4, 2, 2, 2)
	doTestVerticalSlices(t, 5, 3, 2, 1)
}

func TestRemove(t *testing.T) {
	tree := NewDefaultSTRtree()
	tree.Insert(&envelope.Envelope{MaxX: 10, MinX: 0, MaxY: 10, MinY: 0}, "1")
	tree.Insert(&envelope.Envelope{MaxX: 15, MinX: 5, MaxY: 15, MinY: 5}, "2")
	tree.Insert(&envelope.Envelope{MaxX: 20, MinX: 10, MaxY: 20, MinY: 10}, "3")
	tree.Insert(&envelope.Envelope{MaxX: 25, MinX: 15, MaxY: 25, MinY: 15}, "4")
	tree.remove(&envelope.Envelope{MaxX: 20, MinX: 10, MaxY: 20, MinY: 10}, "4")
	fmt.Printf("----tree:%+v\n", tree.AbstractSTRtree.Root)
	assert.Equal(t, 3, tree.Size())
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

func doTestVerticalSlices(t *testing.T, itemCount, sliceCount,
	expectedBoundablesPerSlice, expectedBoundablesOnLastSlice int) {
	tree := NewSTRtree(2)
	slices := tree.verticalSlices(itemWrappers(itemCount), sliceCount)
	assert.Equal(t, sliceCount, len(slices))
	for i := 0; i < sliceCount-1; i++ {
		assert.Equal(t, expectedBoundablesPerSlice, len(slices[i]))
	}
	assert.Equal(t, expectedBoundablesOnLastSlice, len(slices[sliceCount-1]))
}

func itemWrappers(size int) (itemWrappers []Boundable) {
	for i := 0; i < size; i++ {
		itemWrappers = append(itemWrappers, &ItemBoundable{&envelope.Envelope{0, 0, 0, 0}, nil})
	}
	return
}
