package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/overlay/chain"
	"github.com/spatial-go/geoos/index"
	"github.com/spatial-go/geoos/index/strtree"
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
	nodedSegStrings  interface{}
	nOverlaps        int
	overlapTolerance float64
}

// NewMCIndexNoder...
func NewMCIndexNoder() *MCIndexNoder {
	return &MCIndexNoder{
		index: &strtree.STRtree{AbstractSTRtree: new(strtree.AbstractSTRtree)},
	}
}

// ComputeNodes ...
func (m *MCIndexNoder) ComputeNodes(inputSegStrings interface{}) {
	fmt.Printf("====computeNodes2")
	m.nodedSegStrings = inputSegStrings
	inputSS := inputSegStrings.([]*NodedSegmentString)
	for i, _ := range inputSS {
		m.add(inputSS[i])
	}
	fmt.Printf("====computeNodes3 m.index=%#v\n", m.index)
	m.intersectChains()
	fmt.Printf("====computeNodes4")
}

// add ...
func (m *MCIndexNoder) add(segStr *NodedSegmentString) {
	tmpPts := make([][]float64, 0)
	for _, pt := range segStr.pts {
		tmpPts = append(tmpPts, pt)
	}
	segChains := chain.ChainsContext(tmpPts, segStr)
	fmt.Printf("=====segChains:%#v\n", segChains)
	for _, mc := range segChains {
		mc.ID = m.idCounter + 1
		expansion := mc.EnvelopeExpansion(m.overlapTolerance)
		fmt.Println("---expansion", expansion)
		err := m.index.Insert(expansion, mc)
		fmt.Println("m.index.Insert err=", err)
		m.monoChains = append(m.monoChains, mc)
	}
}

// intersectChains...
func (m *MCIndexNoder) intersectChains() {
	overlapAction := NewSegmentOverlapAction(m.segInt)
	for _, queryChain := range m.monoChains {
		fmt.Printf("intersectChains1,queryChain=%+v\n", queryChain)
		queryEnv := queryChain.EnvelopeExpansion(m.overlapTolerance)
		fmt.Printf("intersectChains1,queryEnv=%+v\n", queryEnv)
		fmt.Printf("intersectChains1,m.index=%T\n", m.index)
		overlapChains := m.index.Query(queryEnv) // STRtree
		fmt.Printf("intersectChains2")
		for _, testChain := range overlapChains.([]*chain.MonotoneChain) {
			fmt.Printf("intersectChains3")
			/**
			 * following test makes sure we only compare each pair of chains once
			 * and that we don't compare a chain to itself
			 */
			if testChain.ID > queryChain.ID {
				queryChain.ComputeOverlapsTolerance(testChain, m.overlapTolerance, overlapAction)
				m.nOverlaps++
			}
			// short-circuit if possible
			if m.segInt.isDone() {
				return
			}
			fmt.Printf("intersectChains4")
		}
	}
}

// getNodedSubstrings...
func (m *MCIndexNoder) GetNodedSubstrings() interface{} {
	return nil
}

type SegmentOverlapAction struct {
	*chain.MonotoneChainOverlapAction
	si SegmentIntersector
}

// NewSegmentOverlapAction...
func NewSegmentOverlapAction(si SegmentIntersector) *SegmentOverlapAction {
	return &SegmentOverlapAction{
		si: si,
	}
}

