package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/chain"
	"github.com/spatial-go/geoos/index"
)

// MCIndexNoder Nodes a set of SegmentStrings using a index based on MonotoneChains and a SpatialIndex.
// The SpatialIndex used should be something that supports envelope (range) queries efficiently
// (such as a Quadtree} or STRtree (which is the default index provided).
// The noder supports using an overlap tolerance distance.
// This allows determining segment intersection using a buffer for uses involving snapping with a distance tolerance.
type MCIndexNoder struct {
	SinglePassNoder
	monoChains       []*chain.MonotoneChain
	index            index.SpatialIndex	// STRtree()
	idCounter        int
	nodedSegStrings  []matrix.LineMatrix
	nOverlaps        int
	overlapTolerance float64
}

// computeNodes ...
func (m *MCIndexNoder) ComputeNodes(inputSegStrings interface{}) {
	inputSS := inputSegStrings.([]matrix.LineMatrix)
	m.nodedSegStrings = inputSS
	for i, _ := range inputSS {
		m.add(inputSS[i])
	}
}

// getNodedSubstrings...
func (m *MCIndexNoder) GetNodedSubstrings() interface{} {
	return nil
}

// add ...
func (m *MCIndexNoder) add(InputEdge matrix.LineMatrix) {
	segChains := chain.Chains(InputEdge)
	fmt.Printf("=====segChains:%v\n", segChains)
	for _, mc := range segChains {
		mc.ID = m.idCounter + 1
		//m.index
		m.monoChains = append(m.monoChains, mc)
	}
}
