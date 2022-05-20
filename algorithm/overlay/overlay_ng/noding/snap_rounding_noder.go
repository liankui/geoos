package noding

import (
	"fmt"
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/index/kdtree"
)

const NEARNESS_FACTOR = 100

// SnapRoundingNoder Uses Snap Rounding to compute a rounded, fully noded arrangement from a set of
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
	precisionModel *PrecisionModel
	pixelIndex     *HotPixelIndex
	snappedResult  []*NodedSegmentString
}

// NewSnapRoundingNoder ...
func NewSnapRoundingNoder(pm *PrecisionModel) *SnapRoundingNoder {
	return &SnapRoundingNoder{
		precisionModel: pm,
		pixelIndex:     NewHotPixelIndex(pm),
	}
}

// ComputeNodes Computes the nodes in the snap-rounding line arrangement.
// The nodes are added to the NodedSegmentStrings provided as the input.
func (s *SnapRoundingNoder) ComputeNodes(segStrings interface{}) {
	fmt.Println("====computeNodes8")
	s.snappedResult = s.snapRound(segStrings.([]*NodedSegmentString))
}

// GetNodedSubstrings ...
func (s *SnapRoundingNoder) GetNodedSubstrings() interface{} {
	return nil
}

// snapRound...
func (s *SnapRoundingNoder) snapRound(segStrings []*NodedSegmentString) []*NodedSegmentString {
	/**
	 * Determine hot pixels for intersections and vertices.
	 * This is done BEFORE the input lines are rounded,
	 * to avoid distorting the line arrangement
	 * (rounding can cause vertices to move across edges).
	 */
	s.addIntersectionPixels(segStrings)
	s.addVertexPixels(segStrings)

	snapped := s.computeSnaps(segStrings)
	return snapped
}

// addIntersectionPixels Detects interior intersections in the collection of SegmentStrings,
// and adds nodes for them to the segment strings. Also creates HotPixel nodes for the intersection points.
func (s *SnapRoundingNoder) addIntersectionPixels(segStrings []*NodedSegmentString) {
	// nearness tolerance is a small fraction of the grid size.
	snapGridSize := 1.0 / s.precisionModel.Scale
	nearnessTol := snapGridSize / NEARNESS_FACTOR

	intAdder := NewSnapRoundingIntersectionAdder(nearnessTol)
	noder := NewMCIndexNoderByTolerance(intAdder, nearnessTol)
	noder.ComputeNodes(segStrings)
	s.pixelIndex.addNodes(intAdder.intersections)
}

// addVertexPixels Creates HotPixels for each vertex in the input segStrings. The HotPixels
// are not marked as nodes, since they will only be nodes in the final line arrangement if
// they interact with other segments (or they are already created as intersection nodes).
func (s *SnapRoundingNoder) addVertexPixels(segStrings []*NodedSegmentString) {
	for _, ss := range segStrings {
		s.pixelIndex.addPts(ss.GetCoordinates())
	}
}

// computeSnaps Computes new segment strings which are rounded and contain intersections
// added as a result of snapping segments to snap points (hot pixels).
func (s *SnapRoundingNoder) computeSnaps(segStrings []*NodedSegmentString) []*NodedSegmentString {
	snapped := make([]*NodedSegmentString, 0)
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
func (s *SnapRoundingNoder) computeSegmentSnaps(ss *NodedSegmentString) *NodedSegmentString {
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
	snapSS := NewNodedSegmentString(ptsRound, ss.GetData())
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
func (s *SnapRoundingNoder) snapSegment(p0, p1 matrix.Matrix, ss *NodedSegmentString, segIndex int) {
	s.pixelIndex.query(p0, p1, new(SnapRoundingNoder))
	s.visit(p0, p1, ss, segIndex, new(kdtree.KdNode)) // todo 验证KdNode override处理
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

// addVertexNodeSnaps Add nodes for any vertices in hot pixels that were added as nodes during segment noding.
func (s *SnapRoundingNoder) addVertexNodeSnaps(ss *NodedSegmentString) {
	pts := ss.GetCoordinates()
	for i := 1; i < len(pts)-1; i++ {
		p0 := pts[i]
		s.snapVertexNode(p0, ss, i)
	}
}

// snapVertexNode...
func (s *SnapRoundingNoder) snapVertexNode(p0 matrix.Matrix, ss *NodedSegmentString, segIndex int) {
	s.pixelIndex.query(p0, p0, s)
}

// VisitItem ...
func (s *SnapRoundingNoder) VisitItem(item interface{}) {
	// todo
	//node := item.(*kdtree.KdNode)
	//hp := node.Data.(*HotPixel)
	///**
	// * If vertex pixel is a node, add it.
	// */
	//if hp.isNode && hp.originalPt.Equals(p0) {
	//	ss.addIntersectionNode(p0, segIndex)
	//}
}

func (s *SnapRoundingNoder) Items() interface{} {
	return nil
}

func (s *SnapRoundingNoder) visit(p0, p1 matrix.Matrix, ss *NodedSegmentString, segIndex int, node *kdtree.KdNode) {
	hp := node.Data.(HotPixel)
	/**
	 * If the hot pixel is not a node, and it contains one of the segment vertices,
	 * then that vertex is the source for the hot pixel.
	 * To avoid over-noding a node is not added at this point.
	 * The hot pixel may be subsequently marked as a node,
	 * in which case the intersection will be added during the final vertex noding phase.
	 */
	if !hp.isNode {
		if hp.intersects(p0) || hp.intersects(p1) {
			return
		}
	}
	/**
	 * Add a node if the segment intersects the pixel.
	 * Mark the HotPixel as a node (since it may not have been one before).
	 * This ensures the vertex for it is added as a node during the final vertex noding phase.
	 */
	if hp.intersects2(p0, p1) {
		//System.out.println("Added intersection: " + hp.getCoordinate());
		ss.addIntersectionNode(hp.originalPt, segIndex)
		hp.setToNode()
	}
}
