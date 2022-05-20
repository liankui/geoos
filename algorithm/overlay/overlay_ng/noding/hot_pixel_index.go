package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/matrix/envelope"
	"github.com/spatial-go/geoos/index"
	"github.com/spatial-go/geoos/index/kdtree"
	"math/rand"
	"time"
)

// HotPixelIndex An index which creates unique HotPixels for provided points, and performs range
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

// addPts Adds a list of points as non-node pixels.
func (h *HotPixelIndex) addPts(pts []matrix.Matrix) {
	/**
	 * Shuffle the points before adding.
	 * This avoids having long monontic runs of points
	 * causing an unbalanced KD-tree, which would create
	 * performance and robustness issues.
	 */
	it := CoordinateShuffler(pts)
	for _, i := range it {
		h.add(i)
	}
}

// CoordinateShuffler Utility class to shuffle an array of Coordinates using the Fisher-Yates shuffle algorithm
func CoordinateShuffler(pts []matrix.Matrix) []matrix.Matrix {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pts), func(i, j int) {
		pts[i], pts[j] = pts[j], pts[i]
	})
	return pts
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

// find...
func (h *HotPixelIndex) find(pixelPt matrix.Matrix) *HotPixel {
	kdNode := h.index.QueryMatrix(pixelPt)
	if kdNode == nil {
		return nil
	}
	return kdNode.Data.(*HotPixel)
}

// query Visits all the hot pixels which may intersect a segment (p0-p1). The visitor must determine whether each hot pixel actually intersects the segment.
func (h *HotPixelIndex) query(p0, p1 matrix.Matrix, visitor index.ItemVisitor) { // todo visitor
	queryEnv := envelope.TwoMatrix(p0, p1)
	// expand query range to account for HotPixel extent
	// expand by full width of one pixel to be safe
	queryEnv.ExpandBy(1.0 / h.scaleFactor)
	h.index.QueryVisitor(queryEnv, visitor)
}
