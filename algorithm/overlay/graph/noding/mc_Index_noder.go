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
	index            index.SpatialIndex // STRtree()
	idCounter        int
	nodedSegStrings  []matrix.LineMatrix
	nOverlaps        int
	overlapTolerance float64
}

// computeNodes ...
func (m *MCIndexNoder) ComputeNodes(inputSegStrings interface{}) {
	fmt.Printf("====computeNodes2")
	inputSS := inputSegStrings.([]matrix.LineMatrix)	// todo []*noding.NodedSegmentString
	m.nodedSegStrings = inputSS
	for i, _ := range inputSS {
		m.add(inputSS[i])
	}
	m.intersectChains()
}

// add ...
func (m *MCIndexNoder) add(segStr matrix.LineMatrix) {
	segChains := chain.Chains(segStr)
	fmt.Printf("=====segChains:%v\n", segChains)
	for _, mc := range segChains {
		mc.ID = m.idCounter + 1
		_ = m.index.Insert(mc.EnvelopeExpansion(m.overlapTolerance), mc)
		m.monoChains = append(m.monoChains, mc)
	}
}

// intersectChains...
func (m *MCIndexNoder) intersectChains() {
	overlapAction := chain.NewSegmentOverlapAction(m.segInt)
	for _, queryChain := range m.monoChains {
		queryEnv := queryChain.EnvelopeExpansion(m.overlapTolerance)
		overlapChains := m.index.Query(queryEnv) // STRtree
		for _, testChain := range overlapChains.([]*chain.MonotoneChain) {
			/**
			 * following test makes sure we only compare each pair of chains once
			 * and that we don't compare a chain to itself
			 */
			if testChain.ID > queryChain.ID {
				queryChain.ComputeOverlapsTolerance(testChain, m.overlapTolerance, overlapAction)
				m.nOverlaps++
			}
			// short-circuit if possible
			if m.segInt.IsDone() {
				return
			}
		}
	}
}

// getNodedSubstrings...
func (m *MCIndexNoder) GetNodedSubstrings() interface{} {
	return nil
}

