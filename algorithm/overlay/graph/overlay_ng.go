package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
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
	PrecisionModel *PrecisionModel
	OpCode         int
	Noder          noding.Noder
}

// overlay 主函数入口，得到计算后的多边形
func (o *OverlayNG) overlay(g0, g1 space.Geometry, opCode int) space.Geometry {
	ov := OverlayNG{
		G0:             g0,
		G1:             g1,
		PrecisionModel: NewPrecisionModel(),
		OpCode:         opCode,
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
	// 1
	edges := o.nodeEdges()
	o.buildGraph(edges)

	fmt.Println(edges)
	return nil
}

// nodeEdges...
func (o *OverlayNG) nodeEdges() (edges []matrix.LineMatrix) {
	// Node the edges, using whatever noder is being used
	// 1。1
	nodingBuilder := NewEdgeNodingBuilder(o.PrecisionModel, o.Noder)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if isOptimized {

	}

	// 1。3 DONE
	mergedEdges := nodingBuilder.build(o.G0, o.G1)
	fmt.Printf("mergedEdges:%v\n", mergedEdges)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if isOptimized {

	}

	return
}

// buildGraph...
func (o *OverlayNG) buildGraph(edges []matrix.LineMatrix) {
	var graph OverlayGraph
	for _, e := range edges {
		graph.addEdge(e.Bound(), e.cre)
	}
}
