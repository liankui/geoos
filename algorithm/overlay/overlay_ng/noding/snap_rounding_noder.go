package noding

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
	precisionModel *PrecisionModel
	pixelIndex     *HotPixelIndex
	snappedResult  []NodedSegmentString
}

// NewSnapRoundingNoder...
func NewSnapRoundingNoder(pm *PrecisionModel) *SnapRoundingNoder {
	return &SnapRoundingNoder{
		precisionModel: pm,
		pixelIndex:     NewHotPixelIndex(pm),
	}
}

// computeNodes...
func (s *SnapRoundingNoder) ComputeNodes(segStrings interface{}) {

}

// getNodedSubstrings...
func (s *SnapRoundingNoder) GetNodedSubstrings() interface{} {

	return nil
}
