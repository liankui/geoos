package graph

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/measure"
	"github.com/spatial-go/geoos/space"
	"log"
)

type OverlayEdgeRing struct {
	startEdge *OverlayEdge
	ring      space.Ring
	isHole    bool
	ringPts   []matrix.Matrix
	//locator   IndexedPointInAreaLocator
	shell     *OverlayEdgeRing
	holes     []*OverlayEdgeRing // a list of EdgeRings which are holes in this EdgeRing
}

// NewOverlayEdgeRing...
func NewOverlayEdgeRing(start *OverlayEdge) *OverlayEdgeRing {
	overlayEdgeRing := new(OverlayEdgeRing)
	overlayEdgeRing.startEdge = start
	overlayEdgeRing.ringPts = overlayEdgeRing.computeRingPts(start)
	overlayEdgeRing.computeRing(overlayEdgeRing.ringPts)
	return overlayEdgeRing
}

// setShell Sets the containing shell ring of a ring that has been determined to be a hole.
func (o *OverlayEdgeRing) setShell(shell *OverlayEdgeRing) {
	o.shell = shell
	if shell != nil {
		shell.holes = append(shell.holes, o)
	}
}

// computeRingPts...
func (o *OverlayEdgeRing) computeRingPts(start *OverlayEdge) []matrix.Matrix {
	edge := start
	pts := make([]matrix.Matrix, 0)
	_tk := true
	for _tk || edge != start {
		_tk = false
		if edge.edgeRing == o {
			log.Printf("Edge visited twice during ring-building at %v %v\n", edge.pts, edge.pts)
			return nil
		}
		edge.addCoordinates(pts)
		edge.edgeRing = o
		if edge.nextResultEdge == nil {
			log.Printf("Found null edge in ring %v\n", edge.DirectionPt())
			return nil
		}
		edge = edge.nextResultEdge
	}
	var coordList matrix.CoordinateList = pts
	coordList.CloseRing()
	return coordList
}

// computeRing...
func (o *OverlayEdgeRing) computeRing(ringPts []matrix.Matrix) {
	if o.ring != nil {
		return
	}

	o.ring = space.Ring(matrix.CoordinateList(ringPts).ToLineString()) // geometryFactory.createLinearRing(ringPts)

	var orient measure.Orientation
	o.isHole = orient.IsCCW(matrix.LineMatrix(o.ring))
}

// findEdgeRingContaining Finds the innermost enclosing shell OverlayEdgeRing containing
// this OverlayEdgeRing, if any. The innermost enclosing ring is the smallest enclosing ring.
// The algorithm used depends on the fact that: ring A contains ring B if envelope(ring A)
// contains envelope(ring B) This routine is only safe to use if the chosen point of the hole
// is known to be properly contained in a shell (which is guaranteed to be the case if the hole
// does not touch its shell)
// To improve performance of this function the caller should make the passed shellList as small as
// possible (e.g. by using a spatial index filter beforehand).
// Returns:
//		containing EdgeRing, if there is one or null if no containing EdgeRing is found
func (o *OverlayEdgeRing) findEdgeRingContaining(erList []*OverlayEdgeRing) *OverlayEdgeRing {
	testRing := o.ring
	testEnv := testRing.ComputeEnvelopeInternal()

	minRing := new(OverlayEdgeRing)
	//minRingEnv := new(envelope.Envelope)
	for _, tryEdgeRing := range erList {
		tryRing := tryEdgeRing.ring
		tryShellEnv := tryRing.ComputeEnvelopeInternal()
		// the hole envelope cannot equal the shell envelope
		// (also guards against testing rings against themselves)
		if tryShellEnv.Equals(testEnv) {
			continue
		}
		// hole must be contained in shell
		if !tryShellEnv.Contains(testEnv) {
			continue
		}

		// todo ptNotInList
		//testPt := ptNotInList(testRing.ToMatrix().Bound(), tryEdgeRing.ringPts)
		//isContained := tryEdgeRing.isInRing(testPt)

		// check if the new containing ring is smaller than the current minimum ring
		//if isContained {
		//	if minRing == nil  || minRingEnv.Contains(tryShellEnv) {
		//		minRing = tryEdgeRing
		//		minRingEnv = minRing.ring.GetEnvelopeInternal()
		//	}
		//}
	}
	return minRing
}

// toPolygon Computes the Polygon formed by this ring and any contained holes.
// Returns:
//		the Polygon formed by this ring and its holes.
func (o *OverlayEdgeRing) toPolygon() space.Polygon {
	holeLR := make([]space.LineString, 0)
	if o.holes != nil {
		holeLR = make([]space.LineString, len(holeLR))
		for i := 0; i < len(o.holes); i++ {
			holeLR[i] = space.LineString(o.holes[i].ring)
		}
	}
	var poly space.Polygon
	poly = append(poly, o.ring)
	for _, hole := range holeLR {
		poly = append(poly, hole)
	}
	return poly
}