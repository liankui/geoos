package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/overlay/chain"
	"github.com/spatial-go/geoos/index"
	"github.com/spatial-go/geoos/index/strtree"
	"strconv"
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
		index: &strtree.STRtree{
			AbstractSTRtree: &strtree.AbstractSTRtree{
				NodeCapacity: strtree.DEFAULT_NODE_CAPACITY,
			}},
	}
}

// NewMCIndexNoderByTolerance Creates a new noder with a given SegmentIntersector and
// an overlap tolerance distance to expand intersection tests with.
func NewMCIndexNoderByTolerance(si SegmentIntersector, overlapTolerance float64) *MCIndexNoder {
	mc := MCIndexNoder{
		index: &strtree.STRtree{
			AbstractSTRtree: &strtree.AbstractSTRtree{
				NodeCapacity: strtree.DEFAULT_NODE_CAPACITY,
			}},
		overlapTolerance: overlapTolerance,
	}
	mc.SetSinglePassNoder(si)
	return &mc
}

// ComputeNodes ...
func (m *MCIndexNoder) ComputeNodes(inputSegStrings interface{}) {
	fmt.Printf("====computeNodes2\n")
	m.nodedSegStrings = inputSegStrings
	inputSS := inputSegStrings.([]*NodedSegmentString)
	for i, _ := range inputSS {
		m.add(inputSS[i])
	}
	m.intersectChains()
}

// add ...
func (m *MCIndexNoder) add(segStr *NodedSegmentString) {
	tmpPts := make([][]float64, 0)
	for _, pt := range segStr.pts {
		tmpPts = append(tmpPts, pt)
	}
	fmt.Printf("=====segStr.pts:%+v\n", tmpPts)
	segChains := chain.ChainsContext(tmpPts, segStr)
	for _, mc := range segChains {
		fmt.Printf("=====segChain1.env:%#v\n", mc.Env)
		fmt.Printf("=====segChain1.edge:%#v\n", mc.Edge)

		mc.ID = m.idCounter + 1
		expansion := mc.EnvelopeExpansion(m.overlapTolerance)
		fmt.Printf("====expansion=%v\n", expansion)
		_ = m.index.Insert(expansion, mc)
		//fmt.Printf("insert s.itemBoundables=%#v\n", m.index.(*strtree.STRtree).Query(expansion).(*chain.MonotoneChain))
		m.monoChains = append(m.monoChains, mc)
	}
}

// intersectChains...
func (m *MCIndexNoder) intersectChains() {
	overlapAction := NewSegmentOverlapAction(m.segInt)
	for _, queryChain := range m.monoChains {
		queryEnv := queryChain.EnvelopeExpansion(m.overlapTolerance)
		fmt.Printf("intersectChains1,queryEnv=%+v\n", queryEnv)
		overlapChains := m.index.Query(queryEnv) // STRtree
		//fmt.Printf("000 overlapChains:%#v\n", overlapChains)
		for _, testChain := range overlapChains.([]interface{}) {
			//fmt.Println("intersectChains3")
			/**
			 * following test makes sure we only compare each pair of chains once
			 * and that we don't compare a chain to itself
			 */
			if testChain.(*chain.MonotoneChain).ID > queryChain.ID {
				queryChain.ComputeOverlapsTolerance(testChain.(*chain.MonotoneChain), m.overlapTolerance, overlapAction)
				m.nOverlaps++
			}
			// short-circuit if possible
			if m.segInt.isDone() {
				return
			}
			//fmt.Println("intersectChains4")
		}
	}
}

// todo remove
func PrintlnNodedSegStrings(req []*NodedSegmentString) {
	fmt.Print("===getNodedSubstrings2: ")
	for i, _ := range req {
		fmt.Print(strconv.Itoa(i)+":", req[i], "| ")
	}
	fmt.Println()
}

// getNodedSubstrings...
func (m *MCIndexNoder) GetNodedSubstrings() interface{} {
	PrintlnNodedSegStrings(m.nodedSegStrings.([]*NodedSegmentString))
	var nodeSS NodedSegmentString
	substrings := nodeSS.GetNodedSubstrings(m.nodedSegStrings) // []SegmentString
	for i, _ := range substrings {
		fmt.Println("====substrings=", substrings[i])
	}
	return substrings
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
