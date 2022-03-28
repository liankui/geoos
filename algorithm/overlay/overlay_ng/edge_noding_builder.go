package overlay_ng

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/algorithm/measure"
	"github.com/spatial-go/geoos/algorithm/overlay/overlay_ng/noding"
	"github.com/spatial-go/geoos/space"
)

const (
	MIN_LIMIT_PTS       = 20
	IS_NODING_VALIDATED = true
)

// Builds a set of noded, unique, labelled Edges from the edges of the two input geometries.
// It performs the following steps:
//		Extracts input edges, and attaches topological information
//		if clipping is enabled, handles clipping or limiting input geometry
//		chooses a Noder based on provided precision model, unless a custom one is supplied
//		calls the chosen Noder, with precision model
//		removes any fully collapsed noded edges
//		builds Edges and merges them
type EdgeNodingBuilder struct {
	precisionModel *noding.PrecisionModel
	inputEdges     []*noding.NodedSegmentString
	customNoder    noding.Noder
	clipEnv        *envelope.Envelope
	clipper        *RingClipper
	limiter        *LineLimiter
	hasEdges       [2]bool
}

// NewEdgeNodingBuilder...
func NewEdgeNodingBuilder(pm *noding.PrecisionModel, noder noding.Noder) *EdgeNodingBuilder {
	return &EdgeNodingBuilder{
		precisionModel: pm,
		customNoder:    noder,
	}
}

// createFixedPrecisionNoder...
func (e *EdgeNodingBuilder) createFixedPrecisionNoder(precisionModel *noding.PrecisionModel) noding.Noder {
	noder := noding.NewSnapRoundingNoder(precisionModel)
	return noder
}

// createFloatingPrecisionNoder...
func (e *EdgeNodingBuilder) createFloatingPrecisionNoder(doValidation bool) noding.Noder {
	mcNoder := noding.NewMCIndexNoder()
	li := new(noding.LineIntersector) // todo 待优化
	mcNoder.SetSinglePassNoder(noding.NewIntersectionAdder(li))

	var noder noding.Noder = mcNoder
	if doValidation {
		noder = noding.NewValidatingNoder(noder)
	}
	return noder
}

// setClipEnvelope...
func (e *EdgeNodingBuilder) setClipEnvelope(clipEnv *envelope.Envelope) {
	e.clipEnv = clipEnv
	e.clipper = NewRingClipper(clipEnv)
	e.limiter = NewLineLimiter(clipEnv)
}

// getNoder Gets a noder appropriate for the precision model supplied. This is one of:
// 		Fixed precision: a snap-rounding noder (which should be fully robust)
// 		Floating precision: a conventional nodel (which may be non-robust).
// In this case, a validation step is applied to the output from the noder.
func (e *EdgeNodingBuilder) getNoder() noding.Noder {
	if e.customNoder != nil {
		return e.customNoder
	}

	var overlayUtil OverlayUtil
	if overlayUtil.isFloating(e.precisionModel) {
		fmt.Println("----isFloating")
		return e.createFloatingPrecisionNoder(IS_NODING_VALIDATED)
	}
	fmt.Println("----isFloating2")
	return e.createFixedPrecisionNoder(e.precisionModel)
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
func (e *EdgeNodingBuilder) build(g0, g1 space.Geometry) []*Edge {
	e.add(g0, 0)
	e.add(g1, 1)
	//fmt.Printf("---inputEdges=%#v\n", e.inputEdges[0].GetCoordinates())
	//fmt.Printf("---inputEdges=%#v\n", e.inputEdges[1].GetCoordinates())
	nodedEdges := e.node(e.inputEdges)

	/**
	 * Merge the noded edges to eliminate duplicates.
	 * Labels are combined.
	 */
	edgeMerger := new(EdgeMerger)
	mergedEdges := edgeMerger.merge(nodedEdges)
	fmt.Printf("---mergedEdges=%#v\n", mergedEdges)
	return mergedEdges
}

// node Nodes a set of segment strings and creates Edges from the result.
// The input segment strings each carry a EdgeSourceInfo object, which is
// used to provide source topology info to the constructed Edges (and is then discarded).
func (e *EdgeNodingBuilder) node(segStrings []*noding.NodedSegmentString) []*Edge {
	noder := e.getNoder()
	//fmt.Println("----noder=", noder)
	fmt.Println("----segStrings=", segStrings[0], segStrings[1])
	noder.ComputeNodes(segStrings)
	//fmt.Println("-------nodedSS:pre")
	nodedSS := noder.GetNodedSubstrings()
	fmt.Println("-------nodedSS:", nodedSS)
	nodedEdges := e.createEdges(nodedSS.([]noding.SegmentString))
	return nodedEdges
}

// createEdges...
func (e *EdgeNodingBuilder) createEdges(segStrings []noding.SegmentString) []*Edge {
	edges := make([]*Edge, 0)
	for _, ss := range segStrings {
		pts := ss.GetCoordinates()
		// don't create edges from collapsed lines
		edge := new(Edge)
		if edge.IsCollapsed(pts) {
			continue
		}
		info := ss.GetData().(*EdgeSourceInfo) // 待验证
		// Record that a non-collapsed edge exists for the parent geometry
		e.hasEdges[info.index] = true
		edges = append(edges, NewEdge(pts, info))
	}
	return edges
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
	//fmt.Printf("---addPolygon.shell=%#v\n", shell)
	e.addPolygonRing(shell, false, geomIndex)
	for _, hole := range poly.Holes() {
		//fmt.Printf("---addPolygon.hole=%#v\n", shell)
		e.addPolygonRing(hole, true, geomIndex)
	}

}

// addPolygonRing Adds a polygon ring to the graph. Empty rings are ignored.
func (e *EdgeNodingBuilder) addPolygonRing(ring space.Ring, isHole bool, index int) {
	if ring.IsEmpty() {
		return
	}
	if e.isClippedCompletely(ring.ComputeEnvelopeInternal()) {
		return
	}
	pts := e.clip(ring)
	//fmt.Printf("---addPolygonRing.pts=%#v\n", pts)
	// Don't add edges that collapse to a point
	if len(pts) < 2 {
		return
	}
	depthDelta := e.computeDepthDelta(ring, isHole)
	info := NewEdgeSourceInfo(index, depthDelta, isHole)
	//fmt.Printf("---depthDelta=%v\n", depthDelta)
	//fmt.Printf("---info=%v\n", info)
	e.addEdge(pts, info)
	//for _, edge := range e.inputEdges {
	//	fmt.Printf("e.edge=%#v\n", edge.GetCoordinates())
	//}
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
	pts := make([]matrix.Matrix, 0)
	for _, f := range ring {
		pts = append(pts, f)
	}
	env := ring.ComputeEnvelopeInternal()
	//fmt.Printf("----clip.pts=%#v\n", pts)
	//fmt.Printf("----clip.env=%#v\n", env)

	/**
	 * If no clipper or ring is completely contained then no need to clip.
	 * But repeated points must be removed to ensure correct noding.
	 */
	if e.clipper == nil || e.clipEnv.Covers(env) {
		coordList := matrix.CoordinateList{}
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
	isCCW := o.IsCCW(matrix.LineMatrix(ring))
	/**
	 * Compute whether ring is in canonical orientation or not.
	 * Canonical orientation for the overlay process is
	 * Shells : CW, Holes: CCW
	 */
	isOriented := true
	if !isHole {
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

// addEdge...
func (e *EdgeNodingBuilder) addEdge(pts []matrix.Matrix, info *EdgeSourceInfo) {
	ss := noding.NewNodedSegmentString(pts, info)
	e.inputEdges = append(e.inputEdges, ss)
	//for _, edge := range e.inputEdges {
	//	fmt.Println("--addEdge=", edge)
	//}
}

// hasEdgesFor Reports whether there are noded edges for the given input geometry.
// If there are none, this indicates that either the geometry was empty, or has completely
// collapsed (because it is smaller than the noding precision).
func (e *EdgeNodingBuilder) hasEdgesFor(geomIndex int) bool {
	return e.hasEdges[geomIndex]
}
