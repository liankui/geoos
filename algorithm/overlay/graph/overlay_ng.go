package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/chain"
)

type OverlayNG struct {
	g0, g1         matrix.Steric
	PrecisionModel string
	Noder          interface{}
}

// overlay 主函数入口，得到计算后的多边形
func (*OverlayNG) overlay(g0, g1 matrix.Steric) {

}

// nodeEdges new edge
func (o *OverlayNG) nodeEdges() (edge matrix.LineMatrix) {
	var enb EdgeNodingBuilder
	mergedEdges := enb.build(o.g0, o.g1)
	fmt.Printf("mergedEdges:%v\n", mergedEdges)

	// todo

	return
}

type EdgeNodingBuilder struct {
	InputEdges []matrix.LineMatrix
}

// build Creates a set of labelled {Edge}s.
// Representing the fully noded edges of the input geometries.
// Coincident edges (from the same or both geometries) are merged along with their labels into a single unique, fully labelled edge.
func (e *EdgeNodingBuilder) build(g0, g1 matrix.Steric) (mergedEdges []matrix.Steric) {
	mergedEdges = append(mergedEdges, g0)
	mergedEdges = append(mergedEdges, g1)
	return
}

// add ...
func (e *EdgeNodingBuilder) add(g matrix.LineMatrix, geomIndex int) (matrix.Steric, error) {
	if g == nil {
		return nil, nil
	}

	// todo
	//var t simplify.Trans
	//t.transformRing(g, nil)

	return nil, nil
}

// Nodes a set of segment strings and creates Edges from the result.
// The input segment strings each carry a EdgeSourceInfo object,
// which is used to provide source topology info to the constructed Edges (and is then discarded).
func (e *EdgeNodingBuilder) node() {
	// todo e.InputEdges == inputSegStrings

}

// MCIndexNoder Nodes a set of SegmentStrings using a index based on MonotoneChains and a SpatialIndex.
// The SpatialIndex used should be something that supports envelope (range) queries efficiently
// (such as a Quadtree} or STRtree (which is the default index provided).
// The noder supports using an overlap tolerance distance.
// This allows determining segment intersection using a buffer for uses involving snapping with a distance tolerance.
type MCIndexNoder struct {
	monoChains       []*chain.MonotoneChain
	index            int // todo SpatialIndex
	idCounter        int
	nodedSegStrings  []matrix.LineMatrix
	nOverlaps        int
	overlapTolerance float64
}

// computeNodes ...
func (m *MCIndexNoder) computeNodes(inputSegStrings []matrix.LineMatrix) {
	m.nodedSegStrings = inputSegStrings
	for i, _ := range inputSegStrings {
		m.add(inputSegStrings[i])
	}
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
