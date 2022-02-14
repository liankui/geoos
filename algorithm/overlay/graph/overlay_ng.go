package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/calc"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/noding"
	"github.com/spatial-go/geoos/space"
)

type OverlayNG struct {
	Pm        *noding.PrecisionModel
	OpCode    int
	Noder     noding.Noder
	InputGeom *InputGeometry

	STRICT_MODE_DEFAULT bool // default=false
	isStrictMode        bool // =STRICT_MODE_DEFAULT
	isOptimized         bool // default=true
	isAreaResultOnly    bool
	isOutputEdges       bool
	isOutputResultEdges bool
	isOutputNodedEdges  bool
}

// NewOverlayNG...
func NewOverlayNG(geom0, geom1 space.Geometry, opCode int) *OverlayNG {
	return &OverlayNG{
		Pm:          noding.NewPrecisionModel(),
		OpCode:      opCode,
		InputGeom:   NewInputGeometry(geom0, geom1),
		isOptimized: true,
	}
}

// overlay Computes an overlay operation on the given geometry operands, using the
// precision model of the geometry. and an appropriate noder.
// The noder is chosen according to the precision model specified.
//		For PrecisionModel.FIXED a snap-rounding noder is used, and the computation is robust.
//		For PrecisionModel.FLOATING a non-snapping noder is used, and this computation
//	  		may not be robust. If errors occur a TopologyException is thrown.
// Params:
//		geom0 – the first argument geometry
//		geom1 – the second argument geometry
// 		opCode – the code for the desired overlay operation
// Returns:
//		the result of the overlay operation
func (o *OverlayNG) overlay(g0, g1 space.Geometry, opCode int) space.Geometry {
	ov := NewOverlayNG(g0, g1, opCode)
	return ov.getResult()
}

// getResult Gets the result of the overlay operation.
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
	if o.isOutputNodedEdges {
		return overlayUtil.toLines(graph, o.isOutputEdges)
	}

	//o.labelGraph(graph) // todo

	if o.isOutputEdges || o.isOutputResultEdges {
		return overlayUtil.toLines(graph, o.isOutputEdges)
	}

	return o.extractResult(o.OpCode, graph)
}

// nodeEdges...
func (o *OverlayNG) nodeEdges() []*Edge {
	// Node the edges, using whatever noder is being used
	nodingBuilder := NewEdgeNodingBuilder(o.Pm, o.Noder)

	// Optimize Intersection and Difference by clipping to the
	// result extent, if enabled.
	if o.isOptimized {
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

// todo
//func (o *OverlayNG) labelGraph(graph *OverlayGraph) {
//	labeller := NewOverlayLabeller(graph, o.InputGeom)
//	labeller.computeLabelling()
//	labeller.markResultAreaEdges(opCode)
//	labeller.unmarkDuplicateEdgesFromResultArea()
//}

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
	isAllowMixedIntResult := !o.isStrictMode

	//--- Build polygons
	resultAreaEdges := graph.getResultAreaEdges()
	polyBuilder := NewPolygonBuilder(resultAreaEdges)
	resultPolyList := polyBuilder.getPolygons()
	hasResultAreaComponents := len(resultPolyList) > 0

	resultLineList := make([]space.LineString, 0)
	resultPointList := make([]space.Point, 0)

	if !isAllowMixedIntResult {
		allowResultLines := !hasResultAreaComponents ||
			isAllowMixedIntResult || opCode == SYMDIFFERENCE || opCode == UNION
		if allowResultLines {
			lineBuilder := NewLineBuilder(o.InputGeom, graph, hasResultAreaComponents, opCode)
			lineBuilder.setStrictMode(o.isStrictMode)
			resultLineList = lineBuilder.getLines()
		}
		/**
		 * Operations with point inputs are handled elsewhere.
		 * Only an Intersection op can produce point results
		 * from non-point inputs.
		 */
		hasResultComponents := hasResultAreaComponents || len(resultLineList) > 0
		allowResultPoints := !hasResultComponents || isAllowMixedIntResult
		if o.OpCode == INTERSECTION && allowResultPoints {
			pointBuilder := NewIntersectionPointBuilder(graph, o.isStrictMode)
			pointBuilder.setStrictMode(o.isStrictMode)
			resultPointList = pointBuilder.getPoints()
		}
	}

	if resultPolyList == nil && resultLineList == nil && resultPointList == nil {
		return o.createEmptyResult()
	}

	var overlayUtil OverlayUtil
	resultGeom := overlayUtil.createResultGeometry(resultPolyList, resultLineList, resultPointList)
	return resultGeom
}

// isResultOfOp Tests whether a point with given Locations relative to two geometries
// would be contained in the result of overlaying the geometries using a given overlay operation.
// This is used to determine whether components computed during the overlay process
// should be included in the result geometry.
// The method handles arguments of Location.NONE correctly.
// Params:
//		overlayOpCode – the code for the overlay operation to test
//		loc0 – the code for the location in the first geometry
//		loc1 – the code for the location in the second geometry
func (o *OverlayNG) isResultOfOp(overlayOpCode, loc0, loc1 int) bool {
	if loc0 == calc.ImBoundary {
		loc0 = calc.ImInterior
	}
	if loc1 == calc.ImBoundary {
		loc1 = calc.ImInterior
	}
	switch overlayOpCode {
	case INTERSECTION:
		return loc0 == calc.ImInterior && loc1 == calc.ImInterior
	case UNION:
		return loc0 == calc.ImInterior || loc1 == calc.ImInterior
	case DIFFERENCE:
		return loc0 == calc.ImInterior && loc1 != calc.ImInterior
	case SYMDIFFERENCE:
		return (loc0 == calc.ImInterior && loc1 != calc.ImInterior) ||
			(loc0 != calc.ImInterior && loc1 == calc.ImInterior)
	}
	return false
}

func (o *OverlayNG) createEmptyResult() space.Geometry {
	var overlayUtil OverlayUtil
	return overlayUtil.createEmptyResult(
		overlayUtil.resultDimension(o.OpCode,
			o.InputGeom.getDimension(0),
			o.InputGeom.getDimension(1)))
}
