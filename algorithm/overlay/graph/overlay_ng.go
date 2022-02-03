package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/noding"
	"github.com/spatial-go/geoos/space"
)

const (
	STRICT_MODE_DEFAULT = false
	isStrictMode        = STRICT_MODE_DEFAULT
	isOptimized         = true
	isAreaResultOnly    = false
	isOutputEdges       = false
	isOutputResultEdges = false
	isOutputNodedEdges  = false
)

type OverlayNG struct {
	G0, G1    space.Geometry
	Pm        *PrecisionModel
	OpCode    int
	Noder     noding.Noder
	InputGeom *InputGeometry
	//geomFact GeometryFactory ;
}

// overlay 主函数入口，得到计算后的多边形
func (o *OverlayNG) overlay(g0, g1 space.Geometry, opCode int) space.Geometry {
	ov := OverlayNG{
		G0:     g0,
		G1:     g1,
		Pm:     NewPrecisionModel(),
		OpCode: opCode,
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
	graph := o.buildGraph(edges)

	var overlayUtil OverlayUtil
	if isOutputNodedEdges {
		return overlayUtil.toLines(graph, isOutputEdges)
	}

	o.labelGraph(graph)

	if isOutputEdges || isOutputResultEdges {
		return overlayUtil.toLines(graph, isOutputEdges)
	}

	return o.extractResult(o.OpCode, graph)
}

// nodeEdges...
func (o *OverlayNG) nodeEdges() []*Edge {
	// Node the edges, using whatever noder is being used
	nodingBuilder := NewEdgeNodingBuilder(o.Pm, o.Noder)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if isOptimized {
		var overlayUtil OverlayUtil
		clipEnv := overlayUtil.clippingEnvelope(o.OpCode, o.InputGeom, o.Pm)
		if clipEnv != nil {
			nodingBuilder.setClipEnvelope(clipEnv)
		}
	}

	mergedEdges := nodingBuilder.build(
		o.InputGeom.getGeometry(0),
		o.InputGeom.getGeometry(1))
	fmt.Printf("mergedEdges:%v\n", mergedEdges)

	/**
	 * Record if an input geometry has collapsed.
	 * This is used to avoid trying to locate disconnected edges
	 * against a geometry which has collapsed completely.
	 */
	o.InputGeom.setCollapsed(0, !nodingBuilder.hasEdgesFor(0))
	o.InputGeom.setCollapsed(1, !nodingBuilder.hasEdgesFor(1))

	return mergedEdges
}

// buildGraph...
func (o *OverlayNG) buildGraph(edges []*Edge) *OverlayGraph {
	graph := new(OverlayGraph)
	for _, e := range edges {
		graph.addEdge(e.pts, e.createLabel())
	}
	return graph
}

func (o *OverlayNG) labelGraph(graph *OverlayGraph) {
	labeller := NewOverlayLabeller(graph, o.InputGeom)
	labeller.computeLabelling()
	labeller.markResultAreaEdges(opCode)
	labeller.unmarkDuplicateEdgesFromResultArea()
}

// extractResult Extracts the result geometry components from the fully labelled topology graph.
// This method implements the semantic that the result of an intersection operation
// is homogeneous with highest dimension. In other words, if an intersection has
// components of a given dimension no lower-dimension components are output.
// For example, if two polygons intersect in an area, no linestrings or points are
// included in the result, even if portions of the input do meet in lines or points.
// This semantic choice makes more sense for typical usage, in which only the highest
// dimension components are of interest.
// Params:
//		opCode – the overlay operation
//		graph – the topology graph
// Returns:
//		the result geometry
func (o *OverlayNG) extractResult(opCode int, graph *OverlayGraph) space.Geometry {

}
