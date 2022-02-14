package graph

import (
	"github.com/spatial-go/geoos/algorithm/calc"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/space"
)

// Finds and builds overlay result lines from the overlay graph. Output linework has the following semantics:
//		Linework is fully noded
//		Nodes in the input are preserved in the output
//		Output may contain more nodes than in the input (in particular, sequences of coincident line
//			segments are noded at each vertex
// Various strategies are possible for how to merge graph edges into lines.
// 		This implementation uses the simplest approach of maintaining all nodes arising from noding
// 			(which includes all nodes in the input, and possibly others). This matches the current
//			JTS overlay output semantics.
// 		Another option is to fully merge output lines from node to node. For rings a node point is
// 			chosen arbitrarily. It would also be possible to output LinearRings, if the input is a LinearRing
// 			and is unchanged. This will require additional info from the input linework.
type LineBuilder struct {
	graph          *OverlayGraph
	opCode         int
	inputAreaIndex int
	hasResultArea  bool
	lines          []space.LineString

	// Indicates whether intersections are allowed to produce heterogeneous results including
	// proper boundary touches. This does not control inclusion of touches along collapses.
	// True provides the original JTS semantics.
	isAllowMixedResult bool // todo =! OverlayNG.STRICT_MODE_DEFAULT
	// Allow lines created by area topology collapses to appear in the result.
	// True provides the original JTS semantics.
	isAllowCollapseLines bool // todo =! OverlayNG.STRICT_MODE_DEFAULT
}

// Creates a builder for linear elements which may be present in the overlay result.
//Params:
//inputGeom – the input geometries
//graph – the topology graph
//hasResultArea – true if an area has been generated for the result
//opCode – the overlay operation code
//geomFact – the output geometry factory
func NewLineBuilder(inputGeom *InputGeometry, graph *OverlayGraph, hasResultArea bool, opCode int) *LineBuilder {
	return &LineBuilder{
		graph:          graph,
		opCode:         opCode,
		hasResultArea:  hasResultArea,
		inputAreaIndex: inputGeom.getAreaIndex(),
	}
}

// setStrictMode...
func (l *LineBuilder) setStrictMode(isStrictResultMode bool) {
	l.isAllowCollapseLines = !isStrictResultMode
	l.isAllowMixedResult = !isStrictResultMode
}

// getLines...
func (l *LineBuilder) getLines() []space.LineString {
	l.markResultLines()
	l.addResultLines()
	return l.lines
}

// markResultLines...
func (l *LineBuilder) markResultLines() {
	edges := l.graph.edges
	for _, edge := range edges {
		/**
		 * If the edge linework is already marked as in the result,
		 * it is not included as a line.
		 * This occurs when an edge either is in a result area
		 * or has already been included as a line.
		 */
		if edge.isInResultEither() {
			continue
		}
		if l.isResultLine(edge.label) {
			edge.markInResultLine()
		}
	}
}

// isResultLine Checks if the topology indicated by an edge label determines that this edge
// should be part of a result line.
// Note that the logic here relies on the semantic that for intersection lines are only returned
// if there is no result area components.
// Params:
//		lbl – the label for an edge
// Returns:
//		true if the edge should be included in the result
func (l *LineBuilder) isResultLine(lbl *OverlayLabel) bool {
	/**
	 * Omit edge which is a boundary of a single geometry
	 * (i.e. not a collapse or line edge as well).
	 * These are only included if part of a result area.
	 * This is a short-circuit for the most common area edge case
	 */
	if lbl.isBoundarySingleton() {
		return false
	}
	/**
	 * Omit edge which is a collapse along a boundary.
	 * I.e a result line edge must be from a input line
	 * OR two coincident area boundaries.
	 *
	 * This logic is only used if not including collapse lines in result.
	 */
	if !l.isAllowCollapseLines && lbl.isBoundaryCollapse() {
		return false
	}
	/**
	 * Omit edge which is a collapse interior to its parent area.
	 * (E.g. a narrow gore, or spike off a hole)
	 */
	if lbl.isInteriorCollapse() {
		return false
	}
	/**
	 * For ops other than Intersection, omit a line edge
	 * if it is interior to the other area.
	 *
	 * For Intersection, a line edge interior to an area is included.
	 */
	if l.opCode != INTERSECTION {
		/**
		 * Omit collapsed edge in other area interior.
		 */
		if lbl.isCollapseAndNotPartInterior() {
			return false
		}
		/**
		 * If there is a result area, omit line edge inside it.
		 * It is sufficient to check against the input area rather
		 * than the result area,
		 * because if line edges are present then there is only one input area,
		 * and the result area must be the same as the input area.
		 */
		if l.hasResultArea && lbl.isLineInArea(l.inputAreaIndex) {
			return false
		}
	}
	/**
	 * Include line edge formed by touching area boundaries,
	 * if enabled.
	 */
	if l.isAllowMixedResult && l.opCode == INTERSECTION && lbl.isBoundaryTouch() {
		return true
	}
	/**
	 * Finally, determine included line edge
	 * according to overlay op boolean logic.
	 */
	aLoc := l.effectiveLocation(lbl, 0)
	bLoc := l.effectiveLocation(lbl, 1)
	var overlayNG OverlayNG
	isInResult := overlayNG.isResultOfOp(l.opCode, aLoc, bLoc)
	return isInResult
}

// effectiveLocation Determines the effective location for a line, for the purpose of overlay operation evaluation. Line edges and Collapses are reported as INTERIOR so they may be included in the result if warranted by the effect of the operation on the two edges. (For instance, the intersection of a line edge and a collapsed boundary is included in the result).
//Params:
//lbl – label of line
//geomIndex – index of input geometry
//Returns:
//the effective location of the line
func (l *LineBuilder) effectiveLocation(lbl *OverlayLabel, geomIndex int) int {
	if lbl.isCollapse(geomIndex) {
		return calc.ImInterior
	}
	if lbl.isLineByIndex(geomIndex) {
		return calc.ImInterior
	}
	return lbl.getLineLocation(geomIndex)
}

// addResultLines...
func (l *LineBuilder) addResultLines() {
	edges := l.graph.edges
	for _, edge := range edges {
		if !edge.isInResultLine {
			continue
		}
		if edge.isVisited {
			continue
		}
		l.lines = append(l.lines, l.toLine(edge))
		edge.markVisitedBoth()
	}
}

// toLine...
func (l *LineBuilder) toLine(edge *OverlayEdge) space.LineString {
	isForward := edge.direction
	pts := new(matrix.CoordinateList)
	pts.AddToEndList(edge.origin, false)
	edge.addCoordinates(pts.ToPts())

	ptsOut := pts.ToCoordinateArray(isForward)
	return space.LineString(ptsOut.ToLineString())
}