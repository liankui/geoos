package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/chain"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/noding"
	"github.com/spatial-go/geoos/space"
)

const (
	STRICT_MODE_DEFAULT = false // todo
	isStrictMode        = STRICT_MODE_DEFAULT
	isOptimized         = true
	isAreaResultOnly    = false
	isOutputEdges       = false
	isOutputResultEdges = false
	isOutputNodedEdges  = false
)

type OverlayNG struct {
	G0, G1         space.Geometry
	PrecisionModel string
	OpCode         int
	Noder          noding.Noder
}

// overlay 主函数入口，得到计算后的多边形
func (o *OverlayNG) overlay(g0, g1 space.Geometry, opCode int) space.Geometry {
	ov := OverlayNG{ // todo 类型字段的确定
		G0:             g0,
		G1:             g1,
		PrecisionModel: "Floating",
		OpCode:         opCode,
		Noder:          nil,
	}
	return ov.getResult()
}

// getResult...
func (o *OverlayNG) getResult() space.Geometry {
	// 步骤1： handle empty inputs which determine result

	// 步骤2： The elevation model is only computed if the input geometries have Z values.

	// handle case where both inputs are formed of edges (Lines and Polygons)
	return o.computeEdgeOverlay()
}

// computeEdgeOverlay...
func (o *OverlayNG) computeEdgeOverlay() space.Geometry {
	edges := o.nodeEdges()

	fmt.Println(edges)
	return nil
}

// nodeEdges new edges
func (o *OverlayNG) nodeEdges() (edge matrix.LineMatrix) {
	// Node the edges, using whatever noder is being used
	nodingBuilder := NewEdgeNodingBuilder(o.PrecisionModel, o.Noder)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if isOptimized {

	}

	//
	mergedEdges := nodingBuilder.build(o.G0, o.G1)
	fmt.Printf("mergedEdges:%v\n", mergedEdges)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if isOptimized {

	}

	return
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
