package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"github.com/spatial-go/geoos/algorithm/measure"
	"log"
	"math"
)

// HotPixel Implements a "hot pixel" as used in the Snap Rounding algorithm. A hot pixel
// is a square region centred on the rounded valud of the coordinate given, and
// of width equal to the size of the scale factor. It is a partially open region,
// which contains the interior of the tolerance square and the boundary minus the
// top and right segments. This ensures that every point of the space lies in a
// unique hot pixel. It also matches the rounding semantics for numbers.
// The hot pixel operations are all computed in the integer domain to avoid rounding problems.
// Hot Pixels support being marked as nodes. This is used to prevent introducing
// nodes at line vertices which do not have other lines snapped to them.
type HotPixel struct {
	TOLERANCE   float64 // default = 0.5
	originalPt  matrix.Matrix
	scaleFactor float64
	hpx, hpy    float64
	isNode      bool
}

// NewHotPixel Creates a new hot pixel centered on a rounded point, using a given
// scale factor. The scale factor must be strictly positive (non-zero).
func NewHotPixel(pt matrix.Matrix, scaleFactor float64) *HotPixel {
	h := HotPixel{
		originalPt:  pt,
		scaleFactor: scaleFactor,
	}
	if scaleFactor <= 0 {
		log.Printf("Scale factor must be non-zero")
		return nil
	}
	if scaleFactor != 1.0 {
		h.hpx = math.Round(pt[0] * scaleFactor)
		h.hpy = math.Round(pt[1] * scaleFactor)
	} else {
		h.hpx = pt[0]
		h.hpy = pt[1]
	}
	return &h
}

func (h *HotPixel) setToNode() {
	h.isNode = true
}

// scale Scale without rounding. This ensures intersections are checked against original linework.
// This is required to ensure that intersections are not missed because the segment is moved by snapping.
func (h *HotPixel) scale(val float64) float64 {
	return val * h.scaleFactor
}

// intersects Tests whether a coordinate lies in (intersects) this hot pixel.
func (h *HotPixel) intersects(p matrix.Matrix) bool {
	x := h.scale(p[0])
	y := h.scale(p[1])
	if x >= h.hpx+h.TOLERANCE {
		return false
	}
	// check Left side
	if x < h.hpx-h.TOLERANCE {
		return false
	}
	// check Top side
	if y >= h.hpy+h.TOLERANCE {
		return false
	}
	// check Bottom side
	if y < h.hpy-h.TOLERANCE {
		return false
	}
	return true
}

func (h *HotPixel) intersects2(p0, p1 matrix.Matrix) bool {
	if h.scaleFactor == 1.0 {
		return h.intersectsScaled(p0[0], p0[1], p1[0], p1[1])
	}

	sp0x := h.scale(p0[0])
	sp0y := h.scale(p0[1])
	sp1x := h.scale(p1[0])
	sp1y := h.scale(p1[1])
	return h.intersectsScaled(sp0x, sp0y, sp1x, sp1y)
}

func (h *HotPixel) intersectsScaled(p0x, p0y, p1x, p1y float64) bool {
	// determine oriented segment pointing in positive X direction
	px := p0x
	py := p0y
	qx := p1x
	qy := p1y
	if px > qx {
		px = p1x
		py = p1y
		qx = p0x
		qy = p0y
	}

	/**
	 * Report false if segment env does not intersect pixel env.
	 * This check reflects the fact that the pixel Top and Right sides
	 * are open (not part of the pixel).
	 */
	// check Right side
	maxx := h.hpx + h.TOLERANCE
	segMinx := math.Min(px, qx)
	if segMinx >= maxx {
		return false
	}
	// check Left side
	minx := h.hpx - h.TOLERANCE
	segMaxx := math.Max(px, qx)
	if segMaxx < minx {
		return false
	}
	// check Top side
	maxy := h.hpy + h.TOLERANCE
	segMiny := math.Min(py, qy)
	if segMiny >= maxy {
		return false
	}
	// check Bottom side
	miny := h.hpy - h.TOLERANCE
	segMaxy := math.Max(py, qy)
	if segMaxy < miny {
		return false
	}

	/**
	 * Vertical or horizontal segments must now intersect
	 * the segment interior or Left or Bottom sides.
	 */
	//---- check vertical segment
	if px == qx {
		return true
	}
	//---- check horizontal segment
	if py == qy {
		return true
	}

	/**
	 * Now know segment is not horizontal or vertical.
	 *
	 * Compute orientation WRT each pixel corner.
	 * If corner orientation == 0,
	 * segment intersects the corner.
	 * From the corner and whether segment is heading up or down,
	 * can determine intersection or not.
	 *
	 * Otherwise, check whether segment crosses interior of pixel side
	 * This is the case if the orientations for each corner of the side are different.
	 */
	var cgDD measure.CGAlgorithmsDD
	orientUL := cgDD.OrientationIndex(px, py, qx, qy, minx, maxy)
	if orientUL == 0 {
		// upward segment does not intersect pixel interior
		if py < qy {
			return false
		}
		// downward segment must intersect pixel interior
		return true
	}

	orientUR := cgDD.OrientationIndex(px, py, qx, qy, maxx, maxy)
	if orientUR == 0 {
		// downward segment does not intersect pixel interior
		if py > qy {
			return false
		}
		// upward segment must intersect pixel interior
		return true
	}
	//--- check crossing Top side
	if orientUL != orientUR {
		return true
	}

	orientLL := cgDD.OrientationIndex(px, py, qx, qy, minx, miny)
	if orientLL == 0 {
		// segment crossed LL corner, which is the only one in pixel interior
		return true
	}
	//--- check crossing Left side
	if orientLL != orientUL {
		return true
	}

	orientLR := cgDD.OrientationIndex(px, py, qx, qy, maxx, miny)
	if orientLR == 0 {
		// upward segment does not intersect pixel interior
		if py < qy {
			return false
		}
		// downward segment must intersect pixel interior
		return true
	}

	//--- check crossing Bottom side
	if orientLL != orientLR {
		return true
	}
	//--- check crossing Right side
	if orientLR != orientUR {
		return true
	}

	// segment does not intersect pixel
	return false
}
