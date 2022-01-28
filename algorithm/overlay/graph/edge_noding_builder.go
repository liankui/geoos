package graph

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/algorithm/measure"
	"github.com/spatial-go/geoos/algorithm/overlay/graph/noding"
	"github.com/spatial-go/geoos/space"
)

const (
	MIN_LIMIT_PTS       = 20
	IS_NODING_VALIDATED = true
)

type EdgeNodingBuilder struct {
	precisionModel string
	inputEdges     []*noding.NodedSegmentString
	customNoder    noding.Noder
	clipEnv        *envelope.Envelope
	clipper        *RingClipper
	limiter        *LineLimiter
}

// NewEdgeNodingBuilder...
func NewEdgeNodingBuilder(pm string, noder noding.Noder) *EdgeNodingBuilder {
	return &EdgeNodingBuilder{
		precisionModel: pm,
		customNoder:    noder,
	}
}

// createFloatingPrecisionNoder...
func (e *EdgeNodingBuilder) createFloatingPrecisionNoder(doValidation bool) {
	mcNoder := MCIndexNoder{}
	fmt.Println(mcNoder)
}

func (e *EdgeNodingBuilder) getNoder() noding.Noder {
	if e.customNoder != nil {
		return e.customNoder
	}
	if e.precisionModel == "FLOATING" { // 简化写法
		return nil
	}

	return nil
}

// build Creates a set of labelled {Edge}s.
// Representing the fully noded edges of the input geometries.
// Coincident edges (from the same or both geometries) are merged
// along with their labels into a single unique, fully labelled edge.
// Params:
//		geom0 – the first geometry
//		geom1 – the second geometry
// Returns:
//		the noded, merged, labelled edges
func (e *EdgeNodingBuilder) build(g0, g1 space.Geometry) (mergedEdges []space.Geometry) {
	e.add(g0, 0)
	e.add(g1, 1)
	e.node(e.inputEdges)

	mergedEdges = append(mergedEdges, g0)
	mergedEdges = append(mergedEdges, g1)
	return
}

// node Nodes a set of segment strings and creates Edges from the result.
// The input segment strings each carry a EdgeSourceInfo object, which is
// used to provide source topology info to the constructed Edges (and is then discarded).
func (e *EdgeNodingBuilder) node(segStrings []*noding.NodedSegmentString) {
	noder := e.getNoder()

	fmt.Println(noder)
}

// add...
func (e *EdgeNodingBuilder) add(g space.Geometry, geomIndex int) {
	if g == nil || g.IsEmpty() {
		return
	}
	switch g.(type) {
	case space.LineString:
	case space.MultiLineString:
	case space.Polygon:
		e.addPolygon(g.(space.Polygon), geomIndex)
	case space.MultiPolygon:
	case space.Collection:
	}
}

// addPolygon...
func (e *EdgeNodingBuilder) addPolygon(poly space.Polygon, geomIndex int) {
	shell := poly.Shell()
	e.addPolygonRing(shell, false, geomIndex)
	for _, hole := range poly.Holes() {
		e.addPolygonRing(hole, true, geomIndex)
	}
}

// addPolygonRing Adds a polygon ring to the graph. Empty rings are ignored.
func (e *EdgeNodingBuilder) addPolygonRing(ring space.Ring, isHole bool, index int) {
	if ring.IsEmpty() {
		return
	}
	if e.isClippedCompletely(ring.GetEnvelopeInternal()) {
		return
	}
	pts := e.clip(ring)
	// Don't add edges that collapse to a point
	if len(pts) < 2 {
		return
	}
	depthDelta := e.computeDepthDelta(ring, isHole)
	info := NewEdgeSourceInfo(index, depthDelta, isHole)
	e.addEdge(pts, info)
}

// clip If a clipper is present, clip the line to the clip extent.
// Otherwise, remove duplicate points from the ring.
// If clipping is enabled, then every ring MUST be clipped,
// to ensure that holes are clipped to be inside the shell.
// This means it is not possible to skip clipping for rings with few vertices.
// Params:
//		ring – the line to clip
// Returns:
//		the points in the clipped line
func (e *EdgeNodingBuilder) clip(ring space.Ring) []matrix.Matrix {
	pts := ring.ToMatrix().Bound() // todo Coordinate,xyz坐标系
	coordList := matrix.CoordinateList{}
	env := ring.GetEnvelopeInternal()
	/**
	 * If no clipper or ring is completely contained then no need to clip.
	 * But repeated points must be removed to ensure correct noding.
	 */
	if e.clipper == nil || e.clipEnv.Covers(env) {
		return coordList.RemoveRepeatedPoints(pts)
	}
	return e.clipper.clip(pts)
}

// isClippedCompletely Tests whether a geometry (represented by its envelope)
// lies completely outside the clip extent(if any).
// Params:
//		env – the geometry envelope
// Returns:
//		true if the geometry envelope is outside the clip extent.
func (e *EdgeNodingBuilder) isClippedCompletely(env *envelope.Envelope) bool {
	if e.clipEnv.IsNil() {
		return false
	}
	return e.clipEnv.Disjoint(env)
}

// computeDepthDelta...
func (e *EdgeNodingBuilder) computeDepthDelta(ring space.Ring, isHole bool) int {
	/**
	 * Compute the orientation of the ring, to
	 * allow assigning side interior/exterior labels correctly.
	 * JTS canonical orientation is that shells are CW, holes are CCW.
	 *
	 * It is important to compute orientation on the original ring,
	 * since topology collapse can make the orientation computation give the wrong answer.
	 */
	var o measure.Orientation
	isCCW := o.IsCCW(ring)
	/**
	 * Compute whether ring is in canonical orientation or not.
	 * Canonical orientation for the overlay process is
	 * Shells : CW, Holes: CCW
	 */
	isOriented := true
	if isHole {
		isOriented = !isCCW
	} else {
		isOriented = isCCW
	}
	/**
	 * Depth delta can now be computed.
	 * Canonical depth delta is 1 (Exterior on L, Interior on R).
	 * It is flipped to -1 if the ring is oppositely oriented.
	 */
	depthDelta := 0
	if isOriented {
		depthDelta = 1
	} else {
		depthDelta = -1
	}
	return depthDelta
}

func (e *EdgeNodingBuilder) addEdge(pts []matrix.Matrix, info *EdgeSourceInfo) {
	ss := noding.NewNodedSegmentString(pts, info)
	e.inputEdges = append(e.inputEdges, ss)
}
