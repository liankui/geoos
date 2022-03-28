package noding

import (
	"github.com/spatial-go/geoos/algorithm/matrix"
	"log"
	"math"
)

// Implements a "hot pixel" as used in the Snap Rounding algorithm. A hot pixel
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
