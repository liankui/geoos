package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/index/kdtree"
)

// An index which creates unique HotPixels for provided points, and performs range
// queries on them. The points passed to the index do not needed to be rounded to the
// specified Scale factor; this is done internally when creating the HotPixels for them.
type HotPixelIndex struct {
	precisionModel *PrecisionModel
	scaleFactor    float64
	/**
	 * Use a kd-tree to index the pixel centers for optimum performance.
	 * Since HotPixels have an extent, range queries to the
	 * index must enlarge the query range by a suitable value
	 * (using the pixel width is safest).
	 */
	index *kdtree.KdTree
}

func NewHotPixelIndex(pm *PrecisionModel) *HotPixelIndex {
	return &HotPixelIndex{
		precisionModel: pm,
		scaleFactor:    pm.Scale,
	}
}

// addNodes Adds a list of points as node pixels.
func (h *HotPixelIndex) addNodes(pts []matrix.Matrix) {
	/**
	 * Node points are not shuffled, since they are
	 * added after the vertex points, and hence the KD-tree should
	 * be reasonably balanced already.
	 */
	for _, pt := range pts {
		hp := h.add(pt)
		hp.setToNode()
	}
}

// add Adds a point as a Hot Pixel. If the point has been added already, it is marked as a node.
func (h *HotPixelIndex) add(p matrix.Matrix) *HotPixel {
	pRound := h.round(p)

	hp := h.find(pRound)
	/**
	 * Hot Pixels which are added more than once
	 * must have more than one vertex in them
	 * and thus must be nodes.
	 */
	if hp != nil {
		hp.setToNode()
		return hp
	}
	/**
	 * A pixel containing the point was not found, so create a new one.
	 * It is initially set to NOT be a node
	 * (but may become one later on).
	 */
	hp = NewHotPixel(pRound, h.scaleFactor)
	h.index.InsertMatrix(hp.originalPt, hp)
	return hp
}

// round...
func (h *HotPixelIndex) round(p matrix.Matrix) matrix.Matrix {
	p2 := make(matrix.Matrix, 0)
	copy(p2, p)
	h.precisionModel.MakePrecise(p2)
	return p2
}

func (h *HotPixelIndex) find(pixelPt matrix.Matrix) *HotPixel {
	kdNode := h.index.QueryMatrix(pixelPt)
	if kdNode == nil {
		return nil
	}
	return kdNode.Data.(*HotPixel)
}
