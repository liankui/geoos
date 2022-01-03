package strtree

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"testing"
)

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
		Root: aNode,
		NodeCapacity: 0,
	}
	abSTRtree.build()

	fmt.Printf("abSTRtree.root:%+v\n", *abSTRtree.Root)
	fmt.Printf("abSTRtree.root.Bounds:%+v\n", abSTRtree.Root.Bounds)
}
