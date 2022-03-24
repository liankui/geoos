package snapround

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/overlay/overlay_ng/noding"
)

// Uses Snap Rounding to compute a rounded, fully noded arrangement from a set of
// SegmentStrings, in a performant way, and avoiding unnecessary noding.
// Implements the Snap Rounding technique described in the papers by Hobby,
// Guibas & Marimont, and Goodrich et al. Snap Rounding enforces that all output
// vertices lie on a uniform grid, which is determined by the provided Pm.
// Input vertices do not have to be rounded to the grid beforehand; this is done
// during the snap-rounding process. In fact, rounding cannot be done a priori,
// since rounding vertices by themselves can distort the rounded topology of the arrangement
// (i.e. by moving segments away from hot pixels that would otherwise intersect them,
// or by moving vertices across segments).
// To minimize the number of introduced nodes, the Snap-Rounding Noder avoids
// creating nodes at edge vertices if there is no intersection or snap at that location.
// However, if two different input edges contain identical segments, each of the segment
// vertices will be noded. This still provides fully-noded output. This is the same behaviour
// provided by other noders, such as MCIndexNoder and org.locationtech.jts.noding.snap.SnappingNoder.
type SnapRoundingNoder struct {
	precisionModel *noding.PrecisionModel
	pixelIndex     *HotPixelIndex
	snappedResult  []*noding.NodedSegmentString
}

// NewSnapRoundingNoder...
func NewSnapRoundingNoder(pm *noding.PrecisionModel) *SnapRoundingNoder {
	return &SnapRoundingNoder{
		precisionModel: pm,
		pixelIndex:     NewHotPixelIndex(pm),
	}
}

// computeNodes Computes the nodes in the snap-rounding line arrangement.
// The nodes are added to the NodedSegmentStrings provided as the input.
func (s *SnapRoundingNoder) ComputeNodes(segStrings interface{}) {
	fmt.Println("====computeNodes8")
	s.snappedResult = s.snapRound(segStrings.([]*noding.NodedSegmentString))
}

// getNodedSubstrings...
func (s *SnapRoundingNoder) GetNodedSubstrings() interface{} {

	return nil
}

// snapRound...
func (s *SnapRoundingNoder) snapRound(segStrings []*noding.NodedSegmentString) []*noding.NodedSegmentString {
	/**
	 * Determine hot pixels for intersections and vertices.
	 * This is done BEFORE the input lines are rounded,
	 * to avoid distorting the line arrangement
	 * (rounding can cause vertices to move across edges).
	 */
	addIntersectionPixels(segStrings)
	addVertexPixels(segStrings)

	snapped := s.computeSnaps(segStrings)
	return snapped
}

// computeSnaps Computes new segment strings which are rounded and contain intersections
// added as a result of snapping segments to snap points (hot pixels).
func (s *SnapRoundingNoder) computeSnaps(segStrings []*noding.NodedSegmentString) []*noding.NodedSegmentString {
	snapped := make([]*noding.NodedSegmentString, 0)
	for _, ss := range segStrings {
		snappedSS := s.computeSegmentSnaps(ss)
		if snappedSS != nil {
			snapped = append(snapped, snappedSS)
		}
	}
	/**
	 * Some intersection hot pixels may have been marked as nodes in the previous
	 * loop, so add nodes for them.
	 */
	for _, ss := range snapped {
		s.addVertexNodeSnaps(ss)
	}

	return snapped
}

// Add snapped vertices to a segment string. If the segment string collapses completely due to rounding, null is returned.
func (s *SnapRoundingNoder) computeSegmentSnaps(ss *noding.NodedSegmentString) *noding.NodedSegmentString {
	/**
	 * Get edge coordinates, including added intersection nodes.
	 * The coordinates are now rounded to the grid,
	 * in preparation for snapping to the Hot Pixels
	 */
	pts := ss.GetCoordinates()
	ptsRound := s.round(pts)

	// if complete collapse this edge can be eliminated
	if len(ptsRound) <= 1 {
		return nil
	}

	// Create new nodedSS to allow adding any hot pixel nodes
	snapSS := noding.NewNodedSegmentString(ptsRound, ss.GetData())
	snapSSindex := 0
	for i := 0; i < len(pts)-1; i++ {
		currSnap := snapSS.GetCoordinate(i)
		// If the segment has collapsed completely, skip it
		p1 := pts[i+1]
		p1Round := s.roundPt(p1)
		if p1Round.Equals(currSnap) {
			continue
		}
		p0 := pts[i]

		/**
		 * Add any Hot Pixel intersections with *original* segment to rounded segment.
		 * (It is important to check original segment because rounding can
		 * move it enough to intersect other hot pixels not intersecting original segment)
		 */
		s.snapSegment(p0, p1, snapSS, snapSSindex)
		snapSSindex++
	}
	return snapSS
}

// snapSegment Snaps a segment in a segmentString to HotPixels that it intersects.
func (s *SnapRoundingNoder) snapSegment(p0, p1 matrix.Matrix, ss *noding.NodedSegmentString, segIndex int) {
	KdNodeVisitor()
}

// round Gets a list of the rounded coordinates. Duplicate (collapsed) coordinates are removed.
func (s *SnapRoundingNoder) round(pts []matrix.Matrix) []matrix.Matrix {
	roundPts := matrix.CoordinateList{}
	for i := 0; i < len(pts); i++ {
		roundPts.AddToEndList(s.roundPt(pts[i]), false)
	}
	return roundPts.ToCoordinateArray(true)
}

// roundPt...
func (s *SnapRoundingNoder) roundPt(pt matrix.Matrix) matrix.Matrix {
	var p2 matrix.Matrix
	copy(p2, pt)
	s.precisionModel.MakePrecise(p2)
	return p2
}
