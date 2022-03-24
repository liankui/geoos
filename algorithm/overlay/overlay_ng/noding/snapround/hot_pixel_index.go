package snapround

import (
	"github.com/spatial-go/geoos/algorithm/overlay/overlay_ng/noding"
	"github.com/spatial-go/geoos/index/kdtree"
)

// An index which creates unique HotPixels for provided points, and performs range
// queries on them. The points passed to the index do not needed to be rounded to the
// specified Scale factor; this is done internally when creating the HotPixels for them.
type HotPixelIndex struct {
	precisionModel *noding.PrecisionModel
	scaleFactor    float64
	/**
	 * Use a kd-tree to index the pixel centers for optimum performance.
	 * Since HotPixels have an extent, range queries to the
	 * index must enlarge the query range by a suitable value
	 * (using the pixel width is safest).
	 */
	index *kdtree.KdTree
}

func NewHotPixelIndex(pm *noding.PrecisionModel) *HotPixelIndex {
	return &HotPixelIndex{
		precisionModel: pm,
		scaleFactor:    pm.Scale,
	}
}
